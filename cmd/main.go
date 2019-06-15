package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/zhashkevych/goutalk/application"
	"os"
)

var (
	addr string
)

func main() {
	setupLogging()

	flag.StringVar(&addr,
		"addr",
		"8000",
		"port used to run application")
	flag.Parse()

	log.Infof("Starting GouTalk server")

	app := application.NewApp()
	if err := app.Run(addr); err != nil {
		log.Fatal(err)
	}
	log.Warnf("Gracefully stopping GouTalk")
	app.Stop()
}

func setupLogging() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint:      true,
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
}