package pubsub

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	kafka.Writer
}

func NewKafkaWriter(topic string) *Producer {
	//TODO: Configure TLS Support
	//TODO: Configure SASL Support
	return &Producer{
		kafka.Writer{
			Addr:         kafka.TCP("localhost:9092"),
			Topic:        topic,
			WriteTimeout: 10 * time.Second,
			Balancer:     &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) NewMessage(ctx context.Context, us []byte) error {
	err := p.WriteMessages(ctx, kafka.Message{Value: us})
	if err != nil {
		return err
	}

	return nil
}
