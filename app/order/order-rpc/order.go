package main

import (
	"flag"
	"fmt"

	"Ticket/app/order/order-rpc/internal/config"
	"Ticket/app/order/order-rpc/internal/logic"
	"Ticket/app/order/order-rpc/internal/server"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/app/order/order-rpc/order"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")
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
		order.RegisterOrderServer(grpcServer, server.NewOrderServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(xerr.ErrorInterceptor)
	defer s.Stop()

	// 启动订单 MQ 消费者（异步建订单）
	if ctx.MQ != nil {
		if err := logic.StartOrderConsumer(ctx); err != nil {
			fmt.Printf("[order-rpc] MQ order consumer 启动失败: %v\n", err)
		}
	}
	// 启动死信消费者（建订单失败补偿 + 超时取消）
	if ctx.MQ != nil {
		if err := logic.StartDeadConsumer(ctx); err != nil {
			fmt.Printf("[order-rpc] MQ dead consumer 启动失败: %v\n", err)
		}
	}
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
