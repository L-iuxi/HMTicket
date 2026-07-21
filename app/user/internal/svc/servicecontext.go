package svc

import (
	"Ticket/app/user/internal/config"
	"Ticket/app/user/internal/middleware"
	"Ticket/internal/pkg/db"
	"Ticket/internal/pkg/jwt"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	DB            *gorm.DB
	Jwt           *jwt.JWT
	JwtMiddleware *middleware.JwtMiddleware
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

	jwtTool := jwt.NewJWT(c.Auth.AccessSecret)

	return &ServiceContext{
		Config: c,
		DB:     conn,
		Jwt:    jwtTool,

		JwtMiddleware: middleware.NewJwtMiddleware(jwtTool),
	}

}
