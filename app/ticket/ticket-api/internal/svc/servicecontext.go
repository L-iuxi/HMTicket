// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Ticket/app/ticket/ticket-api/internal/config"
	"Ticket/app/ticket/ticket-rpc/ticketclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	TicketRpc ticketclient.Ticket
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		TicketRpc: ticketclient.NewTicket(zrpc.MustNewClient(c.TicketRpc)),
	}
}
