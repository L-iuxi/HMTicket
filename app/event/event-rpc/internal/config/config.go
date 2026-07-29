package config

import (
	"Ticket/internal/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DSN             string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime int
	}
	BizRedis     redis.RedisConf
	InventoryRpc zrpc.RpcClientConf
}
