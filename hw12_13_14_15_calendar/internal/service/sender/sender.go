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
	//var forever chan struct{}

	messages, err := s.queue.Receive()
	if err != nil {
		s.logger.Error("Failed to receive messages: %v", err)
	}

	for msg := range messages {
		s.logger.Info("Received message: %s", msg)
	}
	//<-forever
}
