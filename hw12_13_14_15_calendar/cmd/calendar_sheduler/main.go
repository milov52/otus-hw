package main

import (
	"flag"
	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"github.com/milov52/hw12_13_14_15_calendar/internal/logger"
	"github.com/milov52/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	memorystorage "github.com/milov52/hw12_13_14_15_calendar/internal/repository/event/memory"
	sqlstorage "github.com/milov52/hw12_13_14_15_calendar/internal/repository/event/sql"
	"github.com/milov52/hw12_13_14_15_calendar/internal/service/scheduler"
	"golang.org/x/net/context"
	"os"
)

const (
	inMemory = "in-memory"
	sql      = "sql"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yaml", "Path to configuration file")
	flag.Parse()

	code := run()
	os.Exit(code)
}

func run() int {
	cfg := config.MustLoad(configFile)
	logg := logger.SetupLogger(cfg.Env)

	var storage scheduler.Storage

	switch cfg.DefaultStorage {
	case inMemory:
		storage = memorystorage.New()
	case sql:
		sqlStorage := sqlstorage.New()
		ctx := context.Background()

		if err := sqlStorage.Connect(ctx, *cfg); err != nil {
			logg.Error("failed to connect to database: " + err.Error())
			return 1 // Возвращаем код ошибки, чтобы завершить программу
		}

		storage = sqlStorage
		defer sqlStorage.Close(ctx) // Убеждаемся, что соединение закроется
	}

	eventQueue, err := queue.NewQueue(cfg)
	if err != nil {
		logg.Error("failed to create queue: " + err.Error())
		return 1 // Возвращаем код ошибки, чтобы завершить программу
	}

	eventScheduler := scheduler.NewScheduler(*logg, storage, eventQueue)
	eventScheduler.Start(context.Background(), cfg.Scheduler.LaunchFrequency)

	return 0
}
