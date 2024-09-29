package calendar

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	"golang.org/x/net/context"
)

const (
	DAY   = 0
	WEEK  = 7
	MONTH = 31
)

type Storage interface {
	CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvents(ctx context.Context, date time.Time, offset int) ([]model.Event, error)
}

type Service struct {
	logger     slog.Logger
	repository Storage
}

func NewEventService(logger slog.Logger, eventProvider Storage) *Service {
	return &Service{
		logger:     logger,
		repository: eventProvider,
	}
}

func (s *Service) CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error) {
	id, err := s.repository.CreateEvent(ctx, event)
	if err != nil {
		s.logger.Error("failed create new event", "err", err)
		return uuid.Nil, err
	}
	s.logger.Info("created new event with id: %s", "id", id)
	return id, nil
}

func (s *Service) UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error {
	err := s.repository.UpdateEvent(ctx, id, event)
	if err != nil {
		s.logger.Error("failed update event", "err", err)
	}
	s.logger.Info("updated event")
	return nil
}

func (s *Service) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	err := s.repository.DeleteEvent(ctx, id)
	if err != nil {
		s.logger.Error("failed delete event", "err", err)
	}
	s.logger.Info("deleted event")
	return nil
}

func (s *Service) DayEventList(ctx context.Context, date time.Time) ([]model.Event, error) {
	eventList, err := s.repository.GetEvents(ctx, date, DAY)
	if err != nil {
		s.logger.Error("failed list event", "err", err)
	}
	s.logger.Info("list event")
	return eventList, nil
}

func (s *Service) WeekEventList(ctx context.Context, startDate time.Time) ([]model.Event, error) {
	eventList, err := s.repository.GetEvents(ctx, startDate, WEEK)
	if err != nil {
		s.logger.Error("failed list event", "err", err)
	}
	s.logger.Info("list event")
	return eventList, nil
}

func (s *Service) MonthEventList(ctx context.Context, startDate time.Time) ([]model.Event, error) {
	eventList, err := s.repository.GetEvents(ctx, startDate, MONTH)
	if err != nil {
		s.logger.Error("failed list event", "err", err)
	}
	s.logger.Info("list event")
	return eventList, nil
}
