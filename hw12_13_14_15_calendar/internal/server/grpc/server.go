package internalgrpc

import (
	"log"
	"log/slog"
	"net"

	"github.com/milov52/hw12_13_14_15_calendar/internal/api/event"
	desc "github.com/milov52/hw12_13_14_15_calendar/pkg/api/event/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	logger     slog.Logger
	controller *event.Controller
}

func NewServer(logger slog.Logger, controller event.Controller) *Server {
	return &Server{
		logger:     logger,
		grpcServer: grpc.NewServer(),
		controller: &controller,
	}
}

func (s *Server) Start(lis net.Listener) error {
	reflection.Register(s.grpcServer)
	desc.RegisterCalendarServer(s.grpcServer, s.controller)

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil { // запускаем grpc сервер
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return nil
}
