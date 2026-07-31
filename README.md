# TicketX — 高并发微服务分布式购票平台

基于 go-zero 从零搭建的微服务票务系统。RPC 集群双实例 + etcd 客户端负载均衡，Redis Sentinel 哨兵故障自动转移，RabbitMQ 异步削峰，OpenTelemetry + Jaeger 全链路追踪。

## 技术栈

| 层 | 技术 |
|---|------|
| 框架 | go-zero v1.10 + gRPC + GORM |
| 服务发现 | etcd |
| 存储 | MySQL, Redis Sentinel（1主1从3哨兵） |
| 消息队列 | RabbitMQ（DLX 死信重试） |
| 网关 | Nginx |
| 安全 | bcrypt + JWT + HMAC 内存缓存 |
| 可观测 | OpenTelemetry + Jaeger 分布式链路追踪 |
| 高可用 | go-zero 内置熔断/降载/超时 + RPC 双实例 |
| 容器化 | Docker + Docker Compose（25 容器） |

## 快速启动

### 方式一：Docker Compose（推荐）

```bash
git clone <repo> && cd Ticket

# 限制并行构建（11 个 Go 服务并行编译可能 OOM）
COMPOSE_PARALLEL_LIMIT=2 docker compose build

# 启动全部 25 个容器
docker compose up -d
```

首次构建 3-5 分钟（拉镜像 + 编译），后续 `docker compose up -d` 秒起。

Docker 模式下所有依赖（MySQL、etcd、Redis Sentinel、RabbitMQ、Nginx）自动启动，无需手动安装。

### 容器化部署

```bash
# 构建（限制并行数避免内存爆）
COMPOSE_PARALLEL_LIMIT=2 docker compose build

# 启动
docker compose up -d

# 查看所有容器状态
docker compose ps

# 单独重建某服务（代码修改后）
docker compose up -d --build order-api-1

# 扩容
docker compose up -d order-api-2

# 停止
docker compose stop

# 停止并清理（数据丢失！）
docker compose down -v
```

**端口映射：**

| 宿主机端口 | 容器端口 | 服务 |
|-----------|---------|------|
| 8090 | 80 | Nginx 统一入口 |
| 3307 | 3306 | MySQL |
| 2379 | 2379 | etcd |
| 6381 | 6379 | Redis Master |
| 6380 | 6380 | Redis Slave |
| 26379-26381 | 26379-26381 | Redis Sentinel ×3 |
| 5672 | 5672 | RabbitMQ |
| 15672 | 15672 | RabbitMQ 管理界面 |
| 4317 | 4317 | OpenTelemetry Collector |

> 端口映射避开宿主机常用端口。如宿主机有 etcd/MySQL/Redis 占用，冲突容器会卡在 `Created` 状态。`sudo fuser -k 2379/tcp` 释放端口再 `docker compose up -d`。

### 查看日志

```bash
# 所有容器实时日志
docker compose logs -f

# 特定服务
docker logs -f ticket-order-api-1

# 最近 100 行
docker logs --tail 100 ticket-admin-api

# 按时间过滤（Docker 日志带时间戳）
docker logs --since 5m ticket-event-rpc

# 搜索错误
docker logs ticket-order-rpc 2>&1 | grep -i error
```

所有 Go 服务日志打到 stdout，Docker 自动收集。不落盘文件。

### 压测

```bash
# 走 Nginx（推荐，模拟生产环境）
go run test/bench/main.go -c 300 -n 20 -addr 127.0.0.1:8090

# 直连服务端口（跳过 Nginx）
go run test/bench/main.go -c 300 -n 20

# raw 模式：无客户端限速，打满 QPS
go run test/bench/main.go -c 300 -n 20 -raw -addr 127.0.0.1:8090
```

| 参数 | 说明 |
|------|------|
| `-c` | 并发用户数 |
| `-n` | 库存票数 |
| `-addr` | Nginx 地址，如 `127.0.0.1:8090` |
| `-raw` | 裸测模式，去掉客户端 100ms sleep，打满 QPS |
| `-t` | 抢购超时秒数，默认 30s |

每轮压测自动创建新活动+场次+票种，注册新用户，互不干扰。

**预期结果：** 0 超卖，全部票售罄。

```bash
# raw 模式（47K QPS，30s 卖出 20 张）
库存: 20  用户: 300  耗时: 30.0s  总请求: 1418064  QPS: 47267
卖出: 20 张  ✅ 精准卖完，0 超卖！

# 普通模式（2.7K QPS，30s 卖出 20 张）
库存: 20  用户: 300  耗时: 30.1s  总请求: 82900   QPS: 2755
卖出: 20 张  ✅ 精准卖完，0 超卖！
```

### 常见问题

**etcd 端口冲突：**
```
Error: failed to bind host port 0.0.0.0:2379/tcp: address already in use
```
```bash
sudo fuser -k 2379/tcp && docker compose up -d
```

**Sentinel DNS 解析失败（Alpine 容器）：**
```
*** FATAL CONFIG FILE ERROR ***
Can't resolve instance hostname.
```
→ sentinel 配置使用硬编码 IP `172.18.0.2`。如 docker network 重建后 IP 变化，查 IP：`docker inspect ticket-redis-master | grep IPAddress`，更新 `internal/redis/conf/docker-sentinel-*.conf`。

**后端 502 Bad Gateway：**
```bash
docker compose restart nginx   # nginx DNS 缓存旧 IP
```

**订单未创建 / MQ 连接失败：**
→ 确认 `order-api-1/2/3` 的 `docker-compose.yml` 有 `RABBITMQ_URL: "amqp://admin:123456@rabbitmq:5672/"` 环境变量。

**MySQL 无 admin 账号：**
```bash
# 先注册用户
curl -X POST http://127.0.0.1:8090/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456","email":"admin@ticket.com"}'

# 改 role 为 admin
docker exec -it ticket-mysql mysql -uroot -p123456 ticket \
  -e "UPDATE users SET role='admin' WHERE username='admin';"
```

### 方式二：裸机启动

**前置依赖：** Go 1.25+、MySQL、etcd、Redis、RabbitMQ

```bash
# 确保 MySQL ticket 库已创建
make start          # 一键启动全部（25 进程，集群模式）
make build          # 编译检查
make stop           # 停止全部
```

```bash
# Nginx（如已安装）
nginx -p /home/xfdhm/nginx -c /home/xfdhm/Ticket/nginx/nginx.conf
```

### Jaeger 链路追踪

```bash
docker run -d --name jaeger -p 4317:4317 -p 16686:16686 jaegertracing/all-in-one:latest
```

浏览器打开 `http://127.0.0.1:16686` 查看调用链。

## 项目结构

```
Ticket/
├── common/
│   ├── xerr/             # 统一错误码 + gRPC 拦截器
│   ├── response/         # 统一 HTTP 响应
│   └── mq/               # RabbitMQ 封装（Docker 自动注入 RABBITMQ_URL）
├── internal/
│   ├── pkg/db/           # MySQL 连接
│   ├── pkg/jwt/          # JWT 签发（sync.Map 内存缓存）
│   └── redis/
│       ├── conf/         # Redis Sentinel 配置（1M1S3S + Docker 版）
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
│   ├── bench/            # 压测（单机 30k QPS，0 超卖）
│   ├── e2e/              # 端到端
│   ├── consistency/      # 一致性
│   └── idem/             # 幂等
├── Dockerfile            # 通用多阶段构建（一个文件打 11 个服务镜像）
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

## 接口手册

全接口 curl 手册见 [test/test.md](test/test.md)。
