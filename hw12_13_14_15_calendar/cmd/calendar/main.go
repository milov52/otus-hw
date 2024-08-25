//nolint:depguard
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/milov52/hw12_13_14_15_calendar/internal/app"
	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"github.com/milov52/hw12_13_14_15_calendar/internal/server/http"
	"github.com/milov52/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/milov52/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	inMemory = "in-memory"
	sql      = "sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg := config.MustLoad(configFile)
	logg := setupLogger(cfg.Env)

	var storage app.Storage
	switch cfg.DefaultStorage {
	case inMemory:
		storage = memorystorage.New()
	case sql:
		storage = sqlstorage.New()
	}
	calendar := app.New(*logg, storage)

	server := internalhttp.NewServer(*logg, *cfg, *calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
