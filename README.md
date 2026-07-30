# TicketX — 高并发微服务分布式购票平台

基于 go-zero 从零搭建的微服务票务系统。RPC 集群双实例 + etcd 客户端负载均衡，Redis Sentinel 哨兵故障自动转移，RabbitMQ 异步削峰，OpenTelemetry + Jaeger 全链路追踪。

## 技术栈

| 层 | 技术 |
|---|------|
| 框架 | go-zero v1.10 + gRPC + GORM |
| 服务发现 | etcd |
| 存储 | MySQL，Redis  |
| 消息队列 | RabbitMQ|
| 网关 | Nginx |
| 安全 | bcrypt + JWT |
| 可观测 | OpenTelemetry + Jaeger 分布式链路追踪 |
| 高可用 | go-zero 内置熔断/降载/超时 + RPC 双实例 |

## 项目结构

```
Ticket/
├── common/
│   ├── xerr/             # 统一错误码 + gRPC 拦截器
│   ├── response/         # 统一 HTTP 响应
│   └── mq/               # RabbitMQ 封装
├── internal/
│   ├── pkg/db/           # MySQL 连接
│   ├── pkg/jwt/          # JWT 签发（带内存缓存）
│   └── redis/
│       └── conf/         # Redis Sentinel 配置文件（1M1S3S）
├── app/
│   ├── user/             # 用户服务
│   ├── admin/            # 管理端
│   ├── event/            # 活动服务 (API + RPC)
│   ├── inventory/        # 库存服务 (API + RPC)
│   ├── order/            # 订单服务 (API + RPC)
│   ├── ticket/           # 票务服务 (API + RPC)
│   └── payment/          # 支付服务 (RPC)
├── nginx/nginx.conf      # Nginx 负载均衡配置
├── test/bench/           # 压测工具
└── Makefile
```

## 服务架构

```
                        Nginx :8090
                            │
    ┌───────┬───────┬───────┼───────┬───────┐
    ▼       ▼       ▼       ▼       ▼       ▼
 user-api  admin  event-api ticket-api    order-api
 :8888    :8889   :8890     :8892         (×3)
                                     :8894 :8895 :8896
    │       │       │       │       │       │       │
    └───────┴───────┴───────┴───────┴───────┴───────┘
                            │ gRPC (go-zero + etcd 客户端 lb)
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
         event-rpc     inventory-rpc   order-rpc ←→ RabbitMQ
         ticket-rpc    payment-rpc
         (各 2 实例)    (各 2 实例)     (各 2 实例)
              │             │             │
              └─────────────┼─────────────┘
                            │
           ┌────────────────┼────────────────┐
           ▼                ▼                 ▼
        MySQL          Redis Sentinel        etcd
                     (1M 1S 3S, 自动故障转移)
```

## 快速启动

### 前置依赖

- Go 1.22+
- MySQL（root:123456@127.0.0.1:3306）
- etcd（127.0.0.1:2379）
- RabbitMQ（127.0.0.1:5672）

启动前确保数据库 `ticket` 已创建。

### 启动

```bash
make start          # 一键启动全部（18 个进程，集群模式）
make build          # 编译检查
make stop           # 停止全部
```

启动后建议运行 Nginx：

```bash
cd /usr/local/nginx/sbin
./nginx -c /home/xfdhm/Ticket/nginx/nginx.conf
```

### Jaeger 链路追踪

```bash
docker run -d --name jaeger -p 4317:4317 -p 16686:16686 jaegertracing/all-in-one:latest
```

浏览器打开 `http://127.0.0.1:16686` 查看调用链。


## 接口手册

全接口 curl 手册见 [test/test.md](test/test.md)。
