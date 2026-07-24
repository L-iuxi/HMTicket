# Ticket — 分布式抢票系统

基于 go-zero 微服务框架的票务平台，支持活动管理、库存扣减、订单支付、出票验票全链路。

## 技术栈

- **语言**: Go 1.25
- **框架**: go-zero v1.10 + gRPC + GORM
- **存储**: MySQL + Redis（Lua 原子库存扣减）+ etcd（服务发现）
- **安全**: bcrypt + JWT

## 项目结构

```
Ticket/
├── common/
│   ├── xerr/             # 统一错误码 + gRPC 拦截器
│   │   ├── code.go       # 错误码常量（按模块分段: 1xxxxx 全局 / 2xxxxx 用户 / …）
│   │   ├── errors.go     # CodeError 结构 + 构造函数
│   │   ├── message.go    # 错误码 → 中文消息映射
│   │   └── interceptor.go # gRPC 拦截器: CodeError → status.Error
│   └── response/
│       └── response.go   # 统一 HTTP 响应 {code, msg, data}
├── internal/
│   ├── pkg/db/           # 共享 DB model（正在拆到各服务）
│   ├── pkg/jwt/          # JWT 签发
│   ├── pkg/util/         # 工具函数
│   └── redis/            # Redis 客户端 + 限流 + 幂等
├── app/
│   ├── user/             # 用户服务 (API)
│   ├── admin/            # 管理端 (API)
│   ├── event/            # 活动服务 (API + RPC)
│   ├── inventory/        # 库存服务 (API + RPC)
│   ├── order/            # 订单服务 (API + RPC)
│   ├── ticket/           # 票务服务 (API + RPC)
│   └── payment/          # 支付服务 (RPC)
└── Makefile
```

## 服务架构

```
用户端 API                    管理端
user-api      :8888          admin-api     :8889
event-api     :8890
inventory-api :8891
ticket-api    :8892
order-api     :8894

RPC 服务                  端口      职责
event-rpc     :8080       活动/场次/票种 CRUD
ticket-rpc    :8081       出票/验票/退票/转赠
payment-rpc   :8082       支付（模拟）
order-rpc     :8083       订单 CRUD + 超时取消
inventory-rpc :8084       Redis 库存扣减/释放
```

## 快速启动

```bash
# 前置：MySQL + Redis + etcd
make start        # 全量启动
make start-cluster # 集群模式（每 RPC 服务 2 实例）
make stop         # 停止所有
make build        # 编译检查
```

## 核心流程

```
POST /api/v1/order/buy
  → 限流检查（用户级 + 令牌桶）
  → 幂等检查（Redis idempotent key）
  → EventRpc.GetTicketType（查票价）
  → 用户限购校验
  → InventoryRpc.DeductStock（Redis Lua 原子扣库存）
  → OrderRpc.CreateOrder（分布式锁 + MySQL 建订单）
  → 支付成功 → TicketRpc.CreateTicket（出票）
  → 支付失败 → 回滚库存 + 取消订单
```

**安全机制**:
- 三层限流：用户级（滑动窗口）+ 票种级（令牌桶）+ 订单级
- 幂等保证：请求级幂等 key（Redis），5 分钟内重复请求直接返回
- 分布式锁：order-rpc / payment-rpc 加 Redis 锁防并发
- 超时取消：order-rpc 后台定时器，15 分钟未支付自动取消 + 回滚库存

## 错误处理体系

### 错误码分段

```
OK = 200
1xxxxx  全局（100001 服务错误 / 100002 参数错误 …）
2xxxxx  用户（200001 用户不存在 / 200002 密码错误 …）
3xxxxx  活动（300001 活动不存在 / 300002 名称不能为空 …）
4xxxxx  票务（400001 票不存在 / 400002 验票失败 …）
5xxxxx  库存（500001 库存不足 / 500002 初始化失败 …）
6xxxxx  订单（600001 订单不存在 / 600006 重复提交 …）
7xxxxx  支付（700001 支付失败 / 700003 重复操作 …）
```

### gRPC 拦截器

所有 RPC 服务自动注入 `xerr.ErrorInterceptor`，将 `*CodeError` 转换为 gRPC `status.Error`：

| xerr 错误 | gRPC Code |
|-----------|-----------|
| `*_NOT_FOUND` | `NotFound` |
| `REQUEST_PARAM_ERROR` | `InvalidArgument` |
| `TOKEN_EXPIRE_ERROR` | `Unauthenticated` |
| `EVENT_STATUS_INVALID` | `FailedPrecondition` |
| `STOCK_NOT_ENOUGH` | `ResourceExhausted` |
| `ORDER_DUPLICATE` | `AlreadyExists` |
| 其他 | `Internal` |

### HTTP 统一响应

所有 API handler 通过 `response.HttpResult` 返回统一格式：

```json
// 成功
{"code":200,"msg":"success","data":{...}}

// 业务错误
{"code":300001,"msg":"活动不存在"}

// 系统错误
{"code":100001,"msg":"服务器开小差啦，稍后再来试一下"}
```

## 数据模型

### User

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | varchar(6) | 用户 ID（6 位随机） |
| Username | varchar(64) | 用户名，唯一 |
| Password | varchar(255) | bcrypt 加密 |
| Email | varchar(128) | 邮箱，唯一 |
| Phone | varchar(20) | 手机号，唯一 |
| Role | varchar(32) | user / admin |
| Gender | uint8 | 0=男 1=女 |

### Event

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Title | varchar(255) | 活动名称 |
| Description | text | 描述 |
| Location | varchar(255) | 地点 |
| CoverImage | varchar(512) | 封面图 |
| StartTime | time | 开始时间 |
| EndTime | time | 结束时间 |
| Status | varchar(32) | draft → ready → selling → closed |
| TotalStock | int | 总票数 |

### Show

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| EventID | uint64 | 所属活动 |
| Name | varchar(128) | 场次名称 |
| ShowTime | time | 开始时间 |
| EndTime | time | 结束时间 |
| Status | varchar(32) | draft / selling / closed |
| Venue | varchar(128) | 场馆 |
| SoldCount | int | 已售数量 |
| SortOrder | int32 | 排序 |

### TicketType

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| EventID | uint64 | 所属活动 |
| ShowID | uint64 | 所属场次 |
| Name | varchar(128) | 票种名称（如 VIP / 普通票） |
| Price | float64 | 价格 |
| Stock | int32 | MySQL 库存 |
| MaxPerUser | int32 | 每人限购 |
| SortOrder | int32 | 排序 |

### Order

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | uint | 用户 |
| EventID | uint64 | 活动 |
| ShowID | uint | 场次 |
| TicketTypeID | uint | 票种 |
| OrderNo | varchar(64) | 订单号（UUID） |
| Quantity | int | 购买数量 |
| TotalPrice | float64 | 总价 |
| Status | varchar(32) | unpaid → paid / failed / cancelled |

### Ticket

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | uint | 用户 |
| EventID | uint | 活动 |
| ShowID | uint | 场次 |
| TicketTypeID | uint | 票种 |
| OrderNo | varchar(64) | 关联订单号 |
| Quantity | int | 数量 |
| TotalPrice | float64 | 总价 |
| Status | varchar(32) | unused / used / refunded |
| QRCode | text | 验票码 |
| RealName | varchar(64) | 实名 |
| IDCard | varchar(32) | 身份证 |

### Redis Key 设计

```
# 库存（Lua 原子扣减）
stock:ticket:{ticketTypeId}

# 分布式锁
lock:order:{userId}:{eventId}:{showId}:{ticketTypeId}
lock:payment:{orderNo}

# 幂等
idem:{userId}:{eventId}:{showId}:{ticketTypeId}:{requestId}
idem:pay:{userId}:{orderNo}:{requestId}

# 限流
limit:buy:user:{userId}
bucket:ticket:{ticketTypeId}
```
