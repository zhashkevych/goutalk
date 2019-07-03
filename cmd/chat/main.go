package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/zhashkevych/goutalk/application"
	"os"
)

var (
	httpAddr string
	wsAddr   string
	dbURI    string
)

func main() {
	setupLogging()

	flag.StringVar(&httpAddr,
		"httpAddr",
		"8000",
		"port used to run application's http server")
	flag.StringVar(&wsAddr,
		"wsAddr",
		"1030",
		"port used to run application's websocket server")
	flag.StringVar(&dbURI,
		"dbURI",
		"mongodb://localhost:27017",
		"mongodb host")

	flag.Parse()

	log.Infof("Starting GouTalk server")

	app := application.NewApp(dbURI)
	if err := app.Run(httpAddr, wsAddr); err != nil {
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
