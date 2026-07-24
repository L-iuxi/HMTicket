package svc

import (
	"Ticket/app/ticket/model"
	"Ticket/app/ticket/ticket-rpc/internal/config"
	"Ticket/internal/pkg/db"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	DB *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn, err := db.NewConnection(
		c.Mysql.DSN,
		c.Mysql.MaxOpenConns,
		c.Mysql.MaxIdleConns,
		c.Mysql.ConnMaxLifetime,
	)

	db.Migrate(conn, &model.Ticket{})

	if err != nil {
		panic(err)
	}

	//依赖注入
	return &ServiceContext{
		Config: c,
		DB:     conn,
	}

}
