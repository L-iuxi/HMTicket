package config

import (
	"Ticket/internal/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	BizRedis redis.RedisConf
}
