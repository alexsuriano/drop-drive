package queue

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL       string
	TopicName string
	Timeout   time.Time
}

type RabbitMQConnection struct {
	cfg  RabbitMQConfig
	conn *amqp.Connection
}

func newRabbitMQConn(cfg RabbitMQConfig) (rc *RabbitMQConnection, err error) {
	rc.cfg = cfg
	rc.conn, err = amqp.Dial(cfg.URL)

	return
}

func (rc *RabbitMQConnection) Publish(msg []byte) error {
	ch, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx,
		"",
		rc.cfg.TopicName,
		false,
		false,
		message)
}

func (rc *RabbitMQConnection) Consume(chDTO chan<- QueueDTO) error {
	ch, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(rc.cfg.TopicName,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for delivery := range messages {
		dto := QueueDTO{}
		dto.Unmarshal(delivery.Body)

		chDTO <- dto
	}

	return nil
}
