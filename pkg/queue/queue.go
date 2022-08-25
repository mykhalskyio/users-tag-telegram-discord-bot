package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Queue struct {
	Kafka *kafka.Conn
}

func GetQueue(address string, topic string) (*Queue, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, 0)
	if err != nil {
		return nil, err
	}

	return &Queue{
		Kafka: conn,
	}, nil
}

func (q *Queue) SendToQueue(msg []byte) {
	q.Kafka.WriteMessages(
		kafka.Message{Value: msg},
	)
}
