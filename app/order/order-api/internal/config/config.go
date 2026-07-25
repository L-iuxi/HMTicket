// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	BizRedis struct {
		Addr     string
		Password string
		DB       int
	}
	OrderRpc     zrpc.RpcClientConf
	InventoryRpc zrpc.RpcClientConf
	PaymentRpc   zrpc.RpcClientConf
	TicketRpc    zrpc.RpcClientConf
	EventRpc     zrpc.RpcClientConf
	AdminAuth    struct {
		AccessSecret string
		AccessExpire int64
	}

	UserAuth struct {
		AccessSecret string
		AccessExpire int64
	}
}
