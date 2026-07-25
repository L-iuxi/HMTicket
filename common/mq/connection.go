package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const DefaultBrokerURL = "amqp://admin:123456@localhost:5672/"

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbitMQ(brokerURL string) (*RabbitMQ, error) {
	if brokerURL == "" {
		brokerURL = DefaultBrokerURL
	}

	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		Conn: conn,
		Ch:   ch,
	}, nil
}

func (mq *RabbitMQ) Close() error {
	if mq.Ch != nil {
		if err := mq.Ch.Close(); err != nil {
			return err
		}
	}
	if mq.Conn != nil {
		return mq.Conn.Close()
	}
	return nil
}
