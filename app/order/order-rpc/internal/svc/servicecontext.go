package svc

import (
	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	"Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/config"
	"Ticket/common/mq"
	"Ticket/internal/pkg/db"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	InventoryRpc inventoryclient.Inventory
	Lock         *redis.DistributedLock
	Idempotent   *redis.Idempotent
	MQ           *mq.RabbitMQ
	Redis        *redis.RedisClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn, err := db.NewConnection(
		c.Mysql.DSN,
		c.Mysql.MaxOpenConns,
		c.Mysql.MaxIdleConns,
		c.Mysql.ConnMaxLifetime,
	)
	if err != nil {
		panic(err)
	}

	db.Migrate(conn, &model.Order{})

	redisClient, err := redis.NewRedisClient(c.BizRedis)
	if err != nil {
		panic(err)
	}

	// 连接 RabbitMQ
	rabbitMQ, err := mq.NewRabbitMQ("")
	if err != nil {
		logx.Errorf("[order-rpc] RabbitMQ 连接失败: %v", err)
		// 不 panic，MQ 不可用时 order-rpc 仍可接受 RPC 调用
	}

	return &ServiceContext{
		Config:       c,
		DB:           conn,
		Redis:        redisClient,
		Lock:         redis.NewDistributedLock(redisClient.Client()),
		Idempotent:   redis.NewIdempotent(redisClient.Client()),
		InventoryRpc: inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		MQ:           rabbitMQ,
	}
}
