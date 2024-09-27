package model

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
	NotifyBefore time.Duration
	Sent         bool
}

type Notification struct {
	EventID string    `json:"event_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	UserID  string    `json:"user_id"`
}
