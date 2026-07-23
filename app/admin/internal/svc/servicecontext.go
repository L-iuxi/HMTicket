// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Ticket/app/admin/internal/config"
	"Ticket/app/admin/internal/middleware"
	"Ticket/app/event/event-rpc/eventclient"
	"Ticket/internal/pkg/db"
	"Ticket/internal/pkg/jwt"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.RedisClient
	DB             *gorm.DB
	Jwt            *jwt.JWT
	AuthMiddleware *middleware.AuthMiddleware
	Admin          rest.Middleware
	EventRpc       eventclient.Event
}

func NewServiceContext(c config.Config) *ServiceContext {
	//连接sql
	conn, err := db.NewConnection(
		c.Mysql.DSN,
		c.Mysql.MaxOpenConns,
		c.Mysql.MaxIdleConns,
		c.Mysql.ConnMaxLifetime,
	)

	if err != nil {
		panic(err)
	}

	//连接redis
	redisClient, err := redis.InitRedis(
		c.Redis.Addr,
		c.Redis.Password,
		c.Redis.DB,
	)
	if err != nil {
		panic(err)
	}
	//初始化jwt
	jwtTool := jwt.NewJWT(c.Auth.AccessSecret)
	//依赖注入
	return &ServiceContext{
		Config:         c,
		DB:             conn,
		Admin:          middleware.NewAdminMiddleware().Handle,
		Jwt:            jwtTool,
		Redis:          redisClient,
		AuthMiddleware: middleware.NewAuthMiddleware(c.Auth.AccessSecret),
		EventRpc:       eventclient.NewEvent(zrpc.MustNewClient(c.EventRpc)),
	}
}
