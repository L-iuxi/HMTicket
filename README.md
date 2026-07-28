# Ticket — 分布式高并发抢票系统

go-zero 微服务票务平台。7 个服务，gRPC + RabbitMQ + Redis Lua 原子库存。

## 技术栈

| 层 | 技术 |
|---|------|
| 框架 | go-zero v1.10 + gRPC + GORM |
| 存储 | MySQL + Redis（Lua 原子脚本） |
| 消息队列 | RabbitMQ（异步建单 + 死信延迟取消） |
| 服务发现 | etcd |
| 网关 | Nginx（多实例负载均衡 + keepalive） |
| 安全 | bcrypt + JWT（缓存 1 分钟） |

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
│   └── redis/            # Lua 脚本 + 限流 + 幂等 + 分布式锁 + 补偿
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
       ┌────────────────────┼────────────────────┐
       ▼                    ▼                    ▼
  order-api:8894      order-api:8895      order-api:8896
       │                    │                    │
       └────────────────────┼────────────────────┘
                            │ gRPC
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
         event-rpc     inventory-rpc   order-rpc ←→ RabbitMQ
         ticket-rpc    payment-rpc
              │             │             │
              └─────────────┼─────────────┘
                            │
                  MySQL + Redis + etcd
```

## 抢票链路

```
POST /api/v1/order/buy
  → Redis Pipeline（限流 + 令牌桶 + 分布式锁，1 次往返）
  → Redis 幂等检查（SETNX + GET）
  → gRPC 查票种 / 活动 / 限购
  → gRPC DeductStock（Redis Lua 原子扣，0 超卖）
  → RabbitMQ 异步建订单
  → 立即返回 orderNo
```

## Redis Lua 脚本

| 脚本 | 作用 |
|------|------|
| `DeductStockLua` | GET + 校验 + DECRBY 原子执行，保证 0 超卖 |
| `RateLimitLua` | 用户固定窗口限流 |
| `TokenBucketLua` | 票种令牌桶平滑限流 |
| `UnlockLua` | 比对 value 后释放，防误删 |

## 中间件要点

- **Redis Pipeline** — 限流 + 令牌桶 + 分布式锁合并为 1 次 Redis 往返
- **RabbitMQ** — 扣库存后异步建订单，API 层只等 Redis（<1ms）
- **死信队列** — 15 分钟 TTL 自动取消未支付订单，回滚库存
- **补偿机制** — 回滚失败写补偿记录，后台定时自动执行（5 分钟）
- **Nginx 多实例** — 3 台 order-api + keepalive 长连接池
- **JWT 缓存** — 已验证 token 缓存 1 分钟，跳过重复 HMAC 计算

## 数据模型

### Redis Key

```
stock:ticket:{ticketTypeId}                           # 库存
lock:order:{userId}:{eventId}:{showId}:{ticketTypeId} # 分布式锁
idem:{userId}:{eventId}:{showId}:{ticketTypeId}:{requestId} # 幂等
limit:buy:user:{userId}                               # 用户限流
bucket:ticket:{ticketTypeId}                          # 令牌桶
compensate:list                                       # 补偿记录列表
```

### MySQL 核心表

- **User** — bcrypt 密码，user/admin 双角色
- **Event** — draft → ready → selling → closed
- **Show** — 场次
- **TicketType** — 票种（价格 / 库存 / 限购）
- **Order** — unpaid → paid / cancelled
- **Ticket** — 出票记录

## 性能

单机 **30k QPS**，0 超卖，加机器 N × 30k。详见 [BENCH.md](BENCH.md)。

## 快速启动


```bash
make start          # 一键启动所有服务
make build          # 编译检查
make stop           # 停止全部
```

全接口 curl 手册见 [test/test.md](test/test.md)。
