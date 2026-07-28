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

单机上限 **30k QPS**（go-zero + net/http 栈的物理极限）。多机器水平扩展可达 N × 30k。

| 场景 | 配置 | 总请求 | QPS | 卖出 | 超卖 |
|------|------|--------|-----|------|------|
| 500 人抢 100 张 (raw) | 单实例直连 | 917,090 | 30,566 | 100 | 0 |
| 3×300 人抢 50 张 (raw) | Nginx + 3 实例 | 928,124 | 30,909 | 150 | 0 |

瓶颈在 HTTP 处理本身，不在 Redis / gRPC / MQ。架构支持加机器线性扩展。

## 快速启动

### 1. 前置依赖

```bash
# MySQL / Redis / etcd
sudo apt install mysql-server redis-server etcd -y

# RabbitMQ
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=123456 \
  rabbitmq:3-management
```

### 2. 启动服务

```bash
make start          # 一键启动所有服务
make build          # 编译检查
make stop           # 停止全部
```

### 3. Nginx 负载均衡（可选，多实例用）

```bash
# 安装
sudo apt install nginx

# 把项目配置软链过去
sudo ln -sf /home/xfdhm/Ticket/nginx/nginx.conf /etc/nginx/nginx.conf

# 启动
sudo nginx

# 重载（修改配置后）
sudo nginx -s reload

# 停止
sudo nginx -s stop
```

默认监听 8090。配置在 `nginx/nginx.conf`，3 台 order-api 轮询分发。非 80 端口无需 root，直接：

```bash
nginx -c /home/xfdhm/Ticket/nginx/nginx.conf
```

## 测试

```bash
go run test/idem/main.go                        # 幂等 + 分布式锁
go run test/rate/main.go                        # 三层限流
go run test/e2e/main.go                         # 端到端全链路
go run test/bench/main.go -c 200 -n 50          # 压测: N人抢M张
go run test/bench/main.go -c 200 -n 50 -raw -addr 127.0.0.1:8090  # 通过 Nginx
```

全接口 curl 手册见 [test/test.md](test/test.md)。
