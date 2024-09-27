package scheduler

import (
	"fmt"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	"golang.org/x/net/context"
	"log/slog"
	"time"
)

type Storage interface {
	GetNotifications(ctx context.Context, date time.Time) ([]model.Notification, error)
}

type QueueMessage interface {
	Send(msg string) error
	Receive() (<-chan string, error)
}

type Scheduler struct {
	logger  slog.Logger
	storage Storage
	queue   QueueMessage
}

func NewScheduler(logger slog.Logger, storage Storage, queue QueueMessage) *Scheduler {
	return &Scheduler{
		logger:  logger,
		storage: storage,
		queue:   queue,
	}
}

func (s *Scheduler) Start(ctx context.Context, freq time.Duration) {
	s.logger.Info("Starting Scheduler...")

	ticker := time.NewTicker(1 * freq)
	for range ticker.C {
		s.processReminders(ctx)
	}
}

func (s *Scheduler) processReminders(ctx context.Context) {
	s.logger.Info("Processing reminders...")
	currentTime := time.Now()

	notifications, err := s.storage.GetNotifications(ctx, currentTime)
	if err != nil {
		s.logger.Error(err.Error())
	}

	for _, n := range notifications {
		msg := fmt.Sprintf("Notification to User: %s, Event ID: %d, Title: %s, Notify At: %s",
			n.UserID, n.EventID, n.Title, n.Date)

		err := s.queue.Send(msg)
		if err != nil {
			s.logger.Error("Error sending message to queue: %v", err)
		}
	}
}
