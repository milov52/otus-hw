package event

import (
	"github.com/google/uuid"
	"github.com/milov52/hw12_13_14_15_calendar/internal/converter/server"
	"github.com/milov52/hw12_13_14_15_calendar/internal/service/event"
	servicepb "github.com/milov52/hw12_13_14_15_calendar/pkg/api/event/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ servicepb.CalendarServer = (*Service)(nil)

type Service struct {
	servicepb.UnimplementedCalendarServer
	impl *event.Service
}

func NewEventServer(impl *event.Service) *Service {
	return &Service{impl: impl}
}

func (s *Service) CreateEvent(ctx context.Context, req *servicepb.CreateRequest) (*servicepb.CreateResponse, error) {
	eventDTO, err := server.EventFromReq(req.GetEvent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	id, err := s.impl.CreateEvent(ctx, *eventDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	return &servicepb.CreateResponse{UUID: id.String()}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *servicepb.UpdateRequest) (*emptypb.Empty, error) {
	eventDTO, err := server.EventFromReq(req.GetEvent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}
	eventID, err := uuid.Parse(req.UUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	err = s.impl.UpdateEvent(ctx, eventID, *eventDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	return nil, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *servicepb.DeleteRequest) (*emptypb.Empty, error) {
	eventID, err := uuid.Parse(req.UUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	err = s.impl.DeleteEvent(ctx, eventID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}
	return nil, nil
}

func (s *Service) GetDayEventList(ctx context.Context, req *servicepb.GetRequest) (*servicepb.GetResponse, error) {
	day := req.GetDate().AsTime()
	events, err := s.impl.DayEventList(ctx, day)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return server.EventsToResp(events), nil
}

func (s *Service) GetWeekEventList(ctx context.Context, req *servicepb.GetRequest) (*servicepb.GetResponse, error) {
	day := req.GetDate().AsTime()
	events, err := s.impl.WeekEventList(ctx, day)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return server.EventsToResp(events), nil
}

func (s *Service) GetMonthEventList(ctx context.Context, req *servicepb.GetRequest) (*servicepb.GetResponse, error) {
	day := req.GetDate().AsTime()
	events, err := s.impl.MonthEventList(ctx, day)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return server.EventsToResp(events), nil
}
