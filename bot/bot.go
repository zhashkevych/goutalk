package bot

import (
	"github.com/gorilla/websocket"
	"github.com/zhashkevych/goutalk/bot/nlu"
	"github.com/zhashkevych/goutalk/bot/nlu/dialogflow"
	"github.com/zhashkevych/goutalk/bot/queue"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

const urlScheme = "ws"

type ChatBot struct {
	wsURL     url.URL
	serverURL url.URL

	processor nlu.Processor

	taskQueue *queue.Queue
	response  chan *queue.Result
}

func NewChatBot(host, path string) (*ChatBot, error) {
	processor, err := dialogflow.NewDialogflowProcessor("goutalkbot-rtwdcq", "en-US", "creds.json")
	if err != nil {
		return nil, err
	}

	return &ChatBot{
		wsURL: url.URL{
			Scheme: urlScheme,
			Host:   host,
			Path:   path,
		},
		taskQueue: queue.NewQueue(processor, 10),
		response:  make(chan *queue.Result),
	}, nil
}

func (c *ChatBot) Run() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	c.taskQueue.Start(c.response)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			m, err := Parse(message)
			if err != nil {
				log.Println("error parsing message:", err)
				continue
			}

			// message contains no @bot mention
			if m == nil {
				continue
			}

			c.taskQueue.Enqueue(m.Text, m.RoomID, m.UserID)
		}
	}()

	for {
		select {
		case <-done:
			return nil
			// case recieved message from queue
		case r := <-c.response:
			if r.Err != nil {
				log.Println("error: %s", r.Err.Error())
				continue
			}

			log.Println("response: %s", r.ResponseMsg)

			err := conn.WriteMessage(websocket.TextMessage, []byte(r.ResponseMsg))
			if err != nil {
				log.Println("write:", err)
				return err
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
