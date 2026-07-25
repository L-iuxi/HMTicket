// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Ticket/app/event/event-api/internal/config"
	"Ticket/app/event/event-rpc/eventclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	EventRpc eventclient.Event
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		EventRpc: eventclient.NewEvent(zrpc.MustNewClient(c.EventRpc)),
	}
}
