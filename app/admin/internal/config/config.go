// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Mysql struct {
		DSN             string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime int
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	EventRpc zrpc.RpcClientConf
}
