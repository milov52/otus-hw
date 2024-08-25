package storage

import (
	"github.com/google/uuid"
	"time"
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
