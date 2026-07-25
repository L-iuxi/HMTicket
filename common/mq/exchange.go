package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TicketExchange = "ticket_exchange"
	DeadExchange   = "dead_exchange"
	DeadQueue      = "dead_queue"
)

// 路由 key
const (
	RoutingOrderCreate = "order.create"
	RoutingDeadMessage = "dead.message" // 与 InitQueue 中的 x-dead-letter-routing-key 一致
)

func (mq *RabbitMQ) InitExchange() error {
	return mq.Ch.ExchangeDeclare(
		TicketExchange, "direct", true, false, false, false, nil,
	)
}

func (mq *RabbitMQ) InitDeadExchange() error {
	return mq.Ch.ExchangeDeclare(
		DeadExchange, "direct", true, false, false, false, nil,
	)
}

// InitQueue 声明带死信配置的业务队列
func (mq *RabbitMQ) InitQueue(queueName string) error {
	args := amqp.Table{
		"x-dead-letter-exchange":    DeadExchange,
		"x-dead-letter-routing-key": RoutingDeadMessage,
	}
	_, err := mq.Ch.QueueDeclare(
		queueName, true, false, false, false, args,
	)
	return err
}

func (mq *RabbitMQ) InitDeadQueue() error {
	_, err := mq.Ch.QueueDeclare(
		DeadQueue, true, false, false, false, nil,
	)
	return err
}

func (mq *RabbitMQ) BindQueue(queueName, routingKey string) error {
	return mq.Ch.QueueBind(
		queueName, routingKey, TicketExchange, false, nil,
	)
}

func (mq *RabbitMQ) BindDeadQueue() error {
	return mq.Ch.QueueBind(
		DeadQueue, RoutingDeadMessage, DeadExchange, false, nil,
	)
}
