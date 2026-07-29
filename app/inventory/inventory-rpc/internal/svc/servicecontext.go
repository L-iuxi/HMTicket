package svc

import (
	"Ticket/app/inventory/inventory-rpc/internal/config"
	"Ticket/internal/redis"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.RedisClient
}

func NewServiceContext(c config.Config) *ServiceContext {

	redisClient, err := redis.NewRedisClient(c.BizRedis)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		Redis:  redisClient,
	}
}
