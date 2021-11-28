package pubsub

import (
	"context"
	"io"
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

func (p *Producer) NewMessage(ctx context.Context, res *io.ReadCloser) error {

	return nil
}
