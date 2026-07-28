package mq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// SendMsg 发送消息到 ticket_exchange
func (mq *RabbitMQ) SendMsg(ctx context.Context, body []byte, routingKey string) error {
	return mq.Ch.PublishWithContext(
		ctx,
		TicketExchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // 消息持久化
			Body:         body,
		},
	)
}
