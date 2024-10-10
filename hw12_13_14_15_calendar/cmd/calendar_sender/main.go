package main

import (
	"flag"
	"os"

	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"github.com/milov52/hw12_13_14_15_calendar/internal/logger"
	"github.com/milov52/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	"github.com/milov52/hw12_13_14_15_calendar/internal/service/sender"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yaml", "Path to configuration file")
	flag.Parse()

	cfg := config.MustLoad(configFile)
	logg := logger.SetupLogger(cfg.Env)

	eventQueue, err := queue.NewQueue(cfg)
	if err != nil {
		logg.Error("failed to create queue: " + err.Error())
		os.Exit(1)
	}
	eventSender := sender.NewSender(*logg, eventQueue)
	eventSender.ReadMessages()
}
