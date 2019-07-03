package bot

import (
	"github.com/gorilla/websocket"
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

	processor LanguageProcessor

	taskQueue *queue.Queue
	response  chan *queue.Result
}

func NewChatBot(host, path string) *ChatBot {
	return &ChatBot{
		wsURL: url.URL{
			Scheme: urlScheme,
			Host:   host,
			Path:   path,
		},
		taskQueue: queue.NewQueue(),
		response:  make(chan *queue.Result),
	}
}

func (c *ChatBot) Run() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return nil
			// case recieved message from queue
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
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
