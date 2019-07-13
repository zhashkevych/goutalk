package bot

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"time"
)

func (c *ChatBot) listen(done chan struct{}) {
	defer close(done)

	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		msg, err := Parse(message)
		if err != nil {
			log.Println("error parsing message:", err)
			continue
		}

		// message contains no @bot mention
		if msg == nil {
			continue
		}

		log.Printf("Recieved message from chat room with ID %s, message: %s", msg.RoomID, msg.Text)

		c.taskQueue.Push(msg.Text, msg.RoomID, msg.UserID)
	}
}

func (c *ChatBot) write(done chan struct{}, interrupt chan os.Signal) error {
	for {
		select {
		case <-done:
			return nil
			// case recieved message from queue
		case r := <-c.response:
			if r.Err != nil {
				c.sendMessage(r.RoomID, "I'm struggling to answer you back at the moment")
				log.Printf("error: %s", r.Err.Error())
				continue
			}

			if err := c.sendMessage(r.RoomID, r.ResponseMsg); err != nil {
				log.Printf("error on message send: %s", err.Error())
				continue
			}

			log.Printf("Sent message to chat room with ID %s, message: %s", r.RoomID, r.ResponseMsg)
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
