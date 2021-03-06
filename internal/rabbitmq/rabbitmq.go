package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitConfig struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
}

func NewConnection() (*RabbitConfig, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"r6index", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	return &RabbitConfig{
		connection: conn,
		channel:    ch,
		queue:      &q,
	}, nil
}

func (p *RabbitConfig) Close() error {
	err := p.channel.Close()
	if err != nil {
		log.Println("error trying to close rabbit channel")
	}

	err = p.connection.Close()
	if err != nil {
		log.Println("error trying to close rabbit connection")
	}

	return nil
}

func (p *RabbitConfig) Produce(b *[]byte) error {
	err := p.channel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        *b,
		})

	return err
}
