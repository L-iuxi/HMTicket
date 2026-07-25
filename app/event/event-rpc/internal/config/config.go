package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DSN             string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime int
	}
	BizRedis struct {
		Addr     string
		Password string
		DB       int
	}
	InventoryRpc zrpc.RpcClientConf
}
