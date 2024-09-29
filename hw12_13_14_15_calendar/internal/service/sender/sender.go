package sender

import (
	"log/slog"
)

type QueueMessage interface {
	Send(msg string) error
	Receive() (<-chan string, error) // Возвращаем канал для чтения сообщений
}

type Sender struct {
	logger slog.Logger
	queue  QueueMessage
}

func NewSender(logger slog.Logger, queue QueueMessage) *Sender {
	return &Sender{
		logger: logger,
		queue:  queue,
	}
}

func (s *Sender) ReadMessages() {
	messages, err := s.queue.Receive()
	if err != nil {
		s.logger.Error("Received message", "err", err)
	}

	for msg := range messages {
		s.logger.Info("Received message", "msg", msg)
	}
}
