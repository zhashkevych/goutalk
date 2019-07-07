package bot

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zhashkevych/goutalk/nlu"
	"github.com/zhashkevych/goutalk/nlu/dialogflow"
	"github.com/zhashkevych/goutalk/queue"
	"github.com/zhashkevych/scheduler"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

const (
	wsURLScheme = "ws"
	lang        = "en-US"
)

type ChatBot struct {
	wsURL      url.URL
	serverHost string

	connection *websocket.Conn

	username  string
	password  string
	authToken string

	processor nlu.Processor

	taskQueue *queue.Queue
	response  chan *queue.Result

	scheduler *scheduler.Scheduler
}

func NewChatBot(wsHost, serverHost, username, password, projectID, credsPath string) (*ChatBot, error) {
	processor, err := dialogflow.NewDialogflowProcessor(projectID, lang, credsPath)
	if err != nil {
		return nil, err
	}

	return &ChatBot{
		wsURL: url.URL{
			Scheme: wsURLScheme,
			Host:   wsHost,
			Path:   "/",
		},
		serverHost: serverHost,

		username: username,
		password: password,

		taskQueue: queue.NewQueue(processor, 10),
		response:  make(chan *queue.Result),
		scheduler: scheduler.NewScheduler(),
	}, nil
}

func (c *ChatBot) Run() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL.String(), nil)
	if err != nil {
		return err
	}
	c.connection = conn

	ctx := context.Background()

	// Initial login request
	c.setAuthToken(ctx)
	// Logging in to GouTalk's server every 24 hours (token lifetime duration)
	c.scheduler.Add(ctx, c.setAuthToken, time.Hour*24)

	c.taskQueue.Start(c.response)

	done := make(chan struct{})

	go c.listen(done)

	return c.write(done, interrupt)
}

func (c *ChatBot) Stop() {
	c.taskQueue.Stop()

	if err := c.connection.Close(); err != nil {
		log.Printf("error occured on connection close: %s", err.Error())
	}
}

func (c *ChatBot) setAuthToken(ctx context.Context) {
	token, err := c.login()
	if err != nil {
		log.Printf("error logging in to GouTalk Chat Server: %s", err.Error())
	}

	c.authToken = token
}
