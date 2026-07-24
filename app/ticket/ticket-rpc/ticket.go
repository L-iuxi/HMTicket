package main

import (
	"flag"
	"fmt"

	"Ticket/app/ticket/ticket-rpc/internal/config"
	"Ticket/app/ticket/ticket-rpc/internal/server"
	"Ticket/app/ticket/ticket-rpc/internal/svc"
	"Ticket/app/ticket/ticket-rpc/ticket"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/ticket.yaml", "the config file")
var port = flag.Int("port", 0, "override listen port (0=use config)")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	if *port > 0 {
		c.ListenOn = fmt.Sprintf("0.0.0.0:%d", *port)
	}
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		ticket.RegisterTicketServer(grpcServer, server.NewTicketServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(xerr.ErrorInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
