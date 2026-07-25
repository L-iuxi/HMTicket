package svc

import (
	"Ticket/app/payment/pay-rpc/internal/config"
	"Ticket/internal/redis"
)

type ServiceContext struct {
	Config config.Config
	Lock   *redis.DistributedLock
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisClient, err := redis.InitRedis(
		c.BizRedis.Addr,
		c.BizRedis.Password,
		c.BizRedis.DB,
	)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config: c,
		Lock:   redis.NewDistributedLock(redisClient.Client()),
	}
}
