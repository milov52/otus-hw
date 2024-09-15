package event

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/milov52/hw12_13_14_15_calendar/internal/converter/server"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	servicepb "github.com/milov52/hw12_13_14_15_calendar/pkg/api/event/v1"
)

var _ servicepb.CalendarServer = (*Controller)(nil)

type EventService interface {
	CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	DayEventList(ctx context.Context, date time.Time) ([]model.Event, error)
	WeekEventList(ctx context.Context, date time.Time) ([]model.Event, error)
	MonthEventList(ctx context.Context, date time.Time) ([]model.Event, error)
}

type Controller struct {
	servicepb.UnimplementedCalendarServer
	eventService EventService
}

func NewEventController(eventService EventService) *Controller {
	return &Controller{eventService: eventService}
}

func (c *Controller) CreateEvent(ctx context.Context, req *servicepb.CreateRequest) (*servicepb.CreateResponse, error) {
	eventDTO, err := server.EventFromReq(req.GetEvent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	id, err := c.eventService.CreateEvent(ctx, *eventDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	return &servicepb.CreateResponse{UUID: id.String()}, nil
}

func (c *Controller) UpdateEvent(ctx context.Context, req *servicepb.UpdateRequest) (*emptypb.Empty, error) {
	eventDTO, err := server.EventFromReq(req.GetEvent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}
	eventID, err := uuid.Parse(req.UUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	err = c.eventService.UpdateEvent(ctx, eventID, *eventDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	return nil, nil
}

func (c *Controller) DeleteEvent(ctx context.Context, req *servicepb.DeleteRequest) (*emptypb.Empty, error) {
	eventID, err := uuid.Parse(req.UUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	err = c.eventService.DeleteEvent(ctx, eventID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}
	return nil, nil
}

func (c *Controller) GetDayEventList(ctx context.Context, req *servicepb.GetRequest) (*servicepb.GetResponse, error) {
	day := req.GetDate().AsTime()
	events, err := c.eventService.DayEventList(ctx, day)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return server.EventsToResp(events), nil
}

func (c *Controller) GetWeekEventList(ctx context.Context, req *servicepb.GetRequest) (*servicepb.GetResponse, error) {
	day := req.GetDate().AsTime()
	events, err := c.eventService.WeekEventList(ctx, day)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return server.EventsToResp(events), nil
}

func (c *Controller) GetMonthEventList(ctx context.Context, req *servicepb.GetRequest) (*servicepb.GetResponse, error) {
	day := req.GetDate().AsTime()
	events, err := c.eventService.MonthEventList(ctx, day)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return server.EventsToResp(events), nil
}
