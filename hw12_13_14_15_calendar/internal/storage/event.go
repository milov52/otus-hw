package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDateBusy      = errors.New("date is busy for this event")
	ErrEventNotFound = errors.New("event not found")
)

type Event struct {
	ID           uuid.UUID
	Title        string
	StartTime    time.Time
	Duration     time.Duration
	Description  string
	UserID       string
	Notification *time.Duration
}

type Notification struct {
	EventID string
	Title   string
	Date    time.Time
	UserID  string
}
