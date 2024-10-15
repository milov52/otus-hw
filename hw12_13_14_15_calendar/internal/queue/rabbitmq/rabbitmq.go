package queue

import (
	"fmt"
	"log"

	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}

func NewQueue(cfg *config.Config) (*Queue, error) {
	amqpConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.Username, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	conn, err := amqp.Dial(amqpConnectionString)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"notifications", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Queue{
		Connection: conn,
		Channel:    ch,
		Queue:      &q,
	}, nil
}

func (q *Queue) Send(msg string) error {
	err := q.Channel.Publish(
		"",           // Exchange
		q.Queue.Name, // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return err
	}

	log.Printf("Event sent to queue: %s", msg)
	return nil
}

func (q *Queue) Receive() (<-chan string, error) {
	msgs, err := q.Channel.Consume(
		q.Queue.Name, // Queue name
		"",           // Consumer
		true,         // Auto-Ack
		false,        // Exclusive
		false,        // No-local
		false,        // No-wait
		nil,          // Args
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
	}

	// Создаем новый канал для строк
	strChan := make(chan string)

	go func() {
		defer close(strChan)

		for msg := range msgs {
			strChan <- string(msg.Body)
		}
	}()

	return strChan, nil
}
