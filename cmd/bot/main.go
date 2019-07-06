package main

import (
	"github.com/kylelemons/go-gypsy/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/zhashkevych/goutalk/bot"
	"os"
)

func main() {
	config, err := yaml.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Error occured while reading config file: %s", err.Error())
		return
	}


	setupLogging()
	chatBot, err := setupChatBot(config)
	if err != nil {
		log.Printf("Error occured while reading config file: %s", err.Error())
		return
	}

	log.Infof("Starting GouTalk Chat Bot")

	if err := chatBot.Run(); err != nil {
		log.Fatalf("Error occured while running ChatBot: %s", err.Error())
	}

	chatBot.Stop()
}

func setupLogging() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint:      true,
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
}

func setupChatBot(cfg *yaml.File) (*bot.ChatBot, error) {
	wsHost, err := cfg.Get("wsHost")
	if err != nil {
		return nil, err
	}

	serverHost, err := cfg.Get("serverHost")
	if err != nil {
		return nil, err
	}

	projectID, err := cfg.Get("projectID")
	if err != nil {
		return nil, err
	}

	jsonPath, err := cfg.Get("jsonPath")
	if err != nil {
		return nil, err
	}

	username, err := cfg.Get("username")
	if err != nil {
		return nil, err
	}

	password, err := cfg.Get("password")
	if err != nil {
		return nil, err
	}

	chatBot, err := bot.NewChatBot(wsHost, serverHost, username, password, projectID, jsonPath)
	if err != nil {
		return nil, err
	}

	return chatBot, nil
}