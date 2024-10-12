package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	queue "github.com/milov52/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/suite"
)

type QueueMessage interface {
	Send(msg string) error
	Receive() (<-chan string, error)
}

type ShedulerSuite struct {
	suite.Suite
	q  QueueMessage
	ch *amqp.Channel
}

func TestServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ShedulerSuite))
}

func (s *ShedulerSuite) SetupSuite() {
	cfg := config.Config{
		RabbitMQ: config.RabbitMQ{
			Host:     "rabbitmq",
			Port:     "5672",
			Username: "guest",
			Password: "guest",
		},
	}
	eventQueue, err := queue.NewQueue(&cfg)
	s.Require().NoError(err)

	s.q = eventQueue
	s.ch = eventQueue.Channel
}

func (s *ShedulerSuite) TestSendMessage() {
	n := model.Notification{
		EventID: uuid.MustParse("53aa35c8-e659-44b2-882f-f6056e443c99"),
		Title:   "notification title",
		Date:    time.Now(),
		UserID:  "1000",
	}
	msg := fmt.Sprintf("Notification to User: %s, Event ID: %s, Title: %s, Notify At: %s",
		n.UserID, n.EventID, n.Title, n.Date)

	err := s.q.Send(msg)
	s.Require().NoError(err)
	// Подписываемся на получение сообщений из очереди
	messages, err := s.ch.Consume(
		"notifications", // Имя очереди
		"",              // Тег потребителя (оставляем пустым, чтобы RabbitMQ сгенерировал его автоматически)
		true,            // Автоматическое подтверждение сообщений
		false,           // Эксклюзивное использование
		false,           // Запрещаем сообщения с локального подключения
		false,           // Ожидание сообщений
		nil,             // Дополнительные аргументы
	)
	s.Require().NoError(err)
	for msg := range messages {
		//fmt.Printf("Received message: %s\n", string(msg.Body))
		s.Require().NotEmpty(msg)
		break // Читаем одно сообщение и выходим
	}
}
