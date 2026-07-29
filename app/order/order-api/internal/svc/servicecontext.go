// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Ticket/app/event/event-rpc/eventclient"
	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	"Ticket/app/order/order-api/internal/config"
	"Ticket/app/order/order-rpc/orderclient"
	"Ticket/app/payment/pay-rpc/paymentclient"
	"Ticket/app/ticket/ticket-rpc/ticketclient"
	"Ticket/common/mq"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	OrderRpc     orderclient.Order
	InventoryRpc inventoryclient.Inventory
	PaymentRpc   paymentclient.Payment
	TicketRpc    ticketclient.Ticket
	EventRpc     eventclient.Event
	Redis        *redis.RedisClient
	Idempotent   *redis.Idempotent
	RateLimiter  *redis.RateLimiter
	TokenBucket  *redis.TokenBucket
	Lock         *redis.DistributedLock
	MQ           *mq.RabbitMQ
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisClient, err := redis.NewRedisClient(c.BizRedis)
	if err != nil {
		panic(err)
	}

	rabbitMQ, err := mq.NewRabbitMQ("")
	if err != nil {
		logx.Errorf("[order-api] RabbitMQ 连接失败: %v", err)
	}

	return &ServiceContext{
		Config:       c,
		OrderRpc:     orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		InventoryRpc: inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		PaymentRpc:   paymentclient.NewPayment(zrpc.MustNewClient(c.PaymentRpc)),
		TicketRpc:    ticketclient.NewTicket(zrpc.MustNewClient(c.TicketRpc)),
		EventRpc:     eventclient.NewEvent(zrpc.MustNewClient(c.EventRpc)),
		Redis:        redisClient,
		Idempotent:   redis.NewIdempotent(redisClient.Client()),
		RateLimiter:  redis.NewRateLimiter(redisClient.Client()),
		TokenBucket:  redis.NewTokenBucket(redisClient.Client()),
		Lock:         redis.NewDistributedLock(redisClient.Client()),
		MQ:           rabbitMQ,
	}
}
