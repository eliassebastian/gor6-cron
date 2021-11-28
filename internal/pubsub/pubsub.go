package pubsub

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Producer *kafka.Conn
	Topic    string
	Ctx      context.Context
}

func NewKafkaConnection(ctx context.Context, topic string) (*Producer, error) {
	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, 0)
	if err != nil {
		return nil, err
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: conn,
		Topic:    topic,
		Ctx:      ctx,
	}, nil
}

func NewProducer(conn *kafka.Conn, topic string) *Producer {
	return &Producer{
		Producer: conn,
		Topic:    topic,
	}
}

func (p *Producer) GetProducer() *kafka.Conn {
	return p.Producer
}
