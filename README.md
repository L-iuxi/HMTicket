# TicketX — 高并发微服务分布式购票平台

基于 go-zero 从零搭建的微服务票务系统。RPC 集群双实例 + etcd 客户端负载均衡，Redis Sentinel 哨兵故障自动转移，RabbitMQ 异步削峰，OpenTelemetry + Jaeger 全链路追踪。

## 技术栈

| 层 | 技术 |
|---|------|
| 框架 | go-zero v1.10 + gRPC + GORM |
| 服务发现 | etcd |
| 存储 | MySQL, Redis Sentinel |
| 消息队列 | RabbitMQ |
| 网关 | Nginx |
| 安全 | bcrypt + JWT + HMAC 内存缓存 |
| 可观测 | OpenTelemetry + Jaeger 分布式链路追踪 |
| 高可用 | go-zero 内置熔断/降载/超时 + RPC 双实例 |
| 容器化 | Docker + Docker Compose |

## 项目结构

```
Ticket/
├── common/
│   ├── xerr/             # 统一错误码 + gRPC 拦截器
│   ├── response/         # 统一 HTTP 响应
│   └── mq/               # RabbitMQ 封装
├── internal/
│   ├── pkg/db/           # MySQL 连接
│   ├── pkg/jwt/          # JWT 签发
│   └── redis/
│       ├── conf/         # Redis Sentinel 配置
│       └── *.go          # Sentinel FailoverClient + Lua 脚本 + 分布式锁/幂等/限流/令牌桶
├── app/
│   ├── user/             # 用户服务
│   ├── admin/            # 管理端 + 补偿任务
│   ├── event/            # 活动服务 (API + RPC)
│   ├── inventory/        # 库存服务 (API + RPC)
│   ├── order/            # 订单服务 (API ×3 + RPC)
│   ├── ticket/           # 票务服务 (API + RPC)
│   └── payment/          # 支付服务 (RPC)
├── nginx/
│   ├── nginx.conf        # 裸机 Nginx 配置
│   └── nginx-docker.conf # Docker Nginx 配置
├── test/
│   ├── bench/            # 压测
│   ├── consistency/      # 一致性
│   └── idem/             # 幂等
├── Dockerfile            # 通用多阶段构建
├── docker-compose.yml    # 25 容器编排
├── Makefile
└── BENCH.md              # 压测报告
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


```bash
#拉取仓库
git clone  && cd HMTicket

# 限制并行构建每次构建两个服务，防止cpu过载
COMPOSE_PARALLEL_LIMIT=2 docker compose build

# 启动全部 25 个容器
docker compose up -d

# 所有容器实时日志
docker compose logs -f
```

首次构建需要 3-5 分钟

Docker 模式下所有依赖（MySQL、etcd、Redis Sentinel、RabbitMQ、Nginx）自动启动，无需手动安装。



## 接口手册

全接口 curl 手册见 [test/test.md](test/test.md)。
