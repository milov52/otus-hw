//nolint:depguard
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/milov52/hw12_13_14_15_calendar/internal/api/event"
	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"github.com/milov52/hw12_13_14_15_calendar/internal/logger"
	"github.com/milov52/hw12_13_14_15_calendar/internal/repository/event/memory"
	"github.com/milov52/hw12_13_14_15_calendar/internal/repository/event/sql"
	internalgrpc "github.com/milov52/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/milov52/hw12_13_14_15_calendar/internal/server/http"
	sevent "github.com/milov52/hw12_13_14_15_calendar/internal/service/calendar"
	desc "github.com/milov52/hw12_13_14_15_calendar/pkg/api/event/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	inMemory = "in-memory"
	sql      = "sql"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yaml", "Path to configuration file")
	flag.Parse()

	cfg := config.MustLoad(configFile)
	logg := logger.SetupLogger(cfg.Env)

	var storage sevent.Storage

	switch cfg.DefaultStorage {
	case inMemory:
		storage = memorystorage.New()
	case sql:
		sqlStorage := sqlstorage.New()
		// Подключаемся к базе данных
		ctx := context.Background()
		if err := sqlStorage.Connect(ctx, *cfg); err != nil {
			logg.Error("failed to connect to database: " + err.Error())
			os.Exit(1)
		}
		storage = sqlStorage
		defer sqlStorage.Close(ctx) // Закрываем соединение при завершении программы
	}

	calendarService := sevent.NewEventService(*logg, storage)
	controller := event.NewEventController(calendarService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServer.Port)) // :82
	if err != nil {
		slog.Error("failed to listen: %v", err)
	}

	grpcServer := internalgrpc.NewServer(*logg, *controller)
	err = grpcServer.Start(lis)
	if err != nil {
		slog.Error("grpc server error", err)
	}

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to dial server", err)
	}

	mux := runtime.NewServeMux()
	err = desc.RegisterCalendarHandler(context.Background(), mux, conn)
	if err != nil {
		slog.Error("failed to register calendar handler", err)
	}

	server := internalhttp.NewServer(*logg, *cfg)
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

	if err := server.Start(mux); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
