package queue

import (
	"fmt"
	"log"
	"reflect"
)

const (
	RabbitMQ QueueType = iota
)

type QueueType int

type QueueConnection interface {
	Publish([]byte) error
	Consume(chan<- QueueDTO) error
}

type Queue struct {
	queueConnection QueueConnection
}

func New(queueType QueueType, cfg any) (queue *Queue, err error) {
	rt := reflect.TypeOf(cfg)

	switch queueType {
	case RabbitMQ:
		if rt.Name() != "RabbitMQConfig" {
			return nil, fmt.Errorf("config need's to be of type RabbitMQConfi")
		}
		conn, err := newRabbitMQConn(cfg.(RabbitMQConfig))
		if err != nil {
			return nil, err
		}

		queue.queueConnection = conn
	default:
		log.Fatal("queue type not implemented")
	}

	return
}

func (q *Queue) Publish(msg []byte) error {
	return q.queueConnection.Publish(msg)
}

func (q *Queue) Consume(chDTO chan<- QueueDTO) error {
	return q.queueConnection.Consume(chDTO)
}
