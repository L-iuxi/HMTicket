package svc

import (
	"Ticket/app/event/event-rpc/internal/config"
	"Ticket/app/event/model"
	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	"Ticket/internal/pkg/db"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	Redis        *redis.RedisClient
	DB           *gorm.DB
	InventoryRpc inventoryclient.Inventory
}

func NewServiceContext(c config.Config) *ServiceContext {

	conn, err := db.NewConnection(
		c.Mysql.DSN,
		c.Mysql.MaxOpenConns,
		c.Mysql.MaxIdleConns,
		c.Mysql.ConnMaxLifetime,
	)
	if err != nil {
		panic(err)
	}

	db.Migrate(conn, &model.Event{}, &model.Show{}, &model.TicketType{})

	//连接redis
	redisClient, err := redis.InitRedis(
		c.BizRedis.Addr,
		c.BizRedis.Password,
		c.BizRedis.DB,
	)
	if err != nil {
		panic(err)
	}
	//依赖注入
	return &ServiceContext{
		Config:       c,
		DB:           conn,
		Redis:        redisClient,
		InventoryRpc: inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
	}
}
