package app

import (
	"context"
	"github.com/milov52/hw12_13_14_15_calendar/internal/storage"
	"log/slog"
	"time"
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

type Logger interface { // TODO
}

type Storage interface {
	CreateEvent(event storage.Event) (string, error)
	UpdateEvent(id string, event storage.Event) error
	DeleteEvent(id string) error
	GetEvents(date time.Time, offset int) ([]storage.Event, error)
}

func New(logger slog.Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	id, err := a.storage.CreateEvent(event)

	if err != nil {
		a.logger.Error("failed create new event", err)
		return err
	}
	a.logger.Info("created new event with id: %s", id)
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	err := a.storage.UpdateEvent(id, event)
	if err != nil {
		a.logger.Error("failed update event", err)
	}
	a.logger.Info("updated event")
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id string, event storage.Event) error {
	err := a.storage.DeleteEvent(id)
	if err != nil {
		a.logger.Error("failed delete event", err)
	}
	a.logger.Info("deleted event")
	return nil
}

func (a *App) DayEventList(ctx context.Context, date time.Time) ([]storage.Event, error) {
	eventList, err := a.storage.GetEvents(date, DAY)
	if err != nil {
		a.logger.Error("failed list event", err)
	}
	a.logger.Info("list event")
	return eventList, nil
}

func (a *App) WeekEventList(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	eventList, err := a.storage.GetEvents(startDate, WEEK)
	if err != nil {
		a.logger.Error("failed list event", err)
	}
	a.logger.Info("list event")
	return eventList, nil
}

func (a *App) MonthEventList(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	eventList, err := a.storage.GetEvents(startDate, MONTH)
	if err != nil {
		a.logger.Error("failed list event", err)
	}
	a.logger.Info("list event")
	return eventList, nil
}
