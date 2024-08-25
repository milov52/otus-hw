package internalhttp

import (
	"context"
	"github.com/milov52/hw12_13_14_15_calendar/internal/app"
	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"log/slog"
	"net"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logger     slog.Logger
	app        app.App
}

func NewServer(logger slog.Logger, cfg config.Config, app app.App) *Server {
	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:         net.JoinHostPort(cfg.HttpServer.Host, cfg.HttpServer.Port),
			ReadTimeout:  cfg.HttpServer.Timeout,
			WriteTimeout: cfg.HttpServer.Timeout,
			IdleTimeout:  cfg.HttpServer.IdleTimeout,
		},
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	})
	s.httpServer.Handler = loggingMiddleware(mux)
	s.logger.Info("starting http server with address", "address", s.httpServer.Addr)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("could not listen on", "address", s.httpServer.Addr, ":", err)
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("shutting down http server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("failed to shutdown http server:", err)
	}
	s.logger.Info("shutting down http server gracefully")
	return nil
}
