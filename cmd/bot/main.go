package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhashkevych/goutalk/bot"
	"os"
)

func main() {
	setupLogging()

	log.Infof("Starting GouTalk Chat Bot")

	chatBot, err := bot.NewChatBot("localhost:1030", "/")
	if err != nil {
		log.Fatalf("Error initializing ChatBot: %s", err.Error())
	}

	if err := chatBot.Run(); err != nil {
		log.Fatalf("Error occured while running ChatBot: %s", err.Error())
	}
}

func setupLogging() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint:      true,
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
}
