// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Ticket/app/inventory/inventory-api/internal/config"
	"Ticket/app/inventory/inventory-rpc/inventoryclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	InventoryRpc inventoryclient.Inventory
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		InventoryRpc: inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
	}
}
