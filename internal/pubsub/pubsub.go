package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
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

/*func (p *Producer) Ping(ctx context.Context) {

}*/

func (p *Producer) NewMessage(ctx context.Context, us interface{}) error {
	b := new(bytes.Buffer)
	defer b.Reset()

	gob.NewEncoder(b).Encode(us)

	err := p.WriteMessages(ctx, kafka.Message{Value: b.Bytes()})
	if err != nil {
		return err
	}

	return nil
}
