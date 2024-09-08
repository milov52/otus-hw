package app

import (
	"context"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type App struct {
	logger  slog.Logger
	storage Storage
}

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

func New(logger slog.Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event model.Event) error {
	id, err := a.storage.CreateEvent(ctx, event)
	if err != nil {
		a.logger.Error("failed create new event", "err", err)
		return err
	}
	a.logger.Info("created new event with id: %s", "id", id)

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error {
	err := a.storage.UpdateEvent(ctx, id, event)
	if err != nil {
		a.logger.Error("failed update event", "err", err)
	}
	a.logger.Info("updated event")
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID, event model.Event) error {
	err := a.storage.DeleteEvent(ctx, id)
	if err != nil {
		a.logger.Error("failed delete event", "err", err)
	}
	a.logger.Info("deleted event")
	return nil
}

func (a *App) DayEventList(ctx context.Context, date time.Time) ([]model.Event, error) {
	eventList, err := a.storage.GetEvents(ctx, date, DAY)
	if err != nil {
		a.logger.Error("failed list event", "err", err)
	}
	a.logger.Info("list event")
	return eventList, nil
}

func (a *App) WeekEventList(ctx context.Context, startDate time.Time) ([]model.Event, error) {
	eventList, err := a.storage.GetEvents(ctx, startDate, WEEK)
	if err != nil {
		a.logger.Error("failed list event", "err", err)
	}
	a.logger.Info("list event")
	return eventList, nil
}

func (a *App) MonthEventList(ctx context.Context, startDate time.Time) ([]model.Event, error) {
	eventList, err := a.storage.GetEvents(ctx, startDate, MONTH)
	if err != nil {
		a.logger.Error("failed list event", "err", err)
	}
	a.logger.Info("list event")
	return eventList, nil
}
