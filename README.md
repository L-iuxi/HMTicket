# Ticket — 分布式抢票系统

基于 go-zero 微服务框架的票务平台。

## 技术栈

- **语言**: Go 1.25
- **框架**: go-zero v1.10 + gRPC + GORM
- **存储**: MySQL + Redis（Lua 原子库存扣减） + BadgerDB（HMETCD）
- **服务发现**: etcd

## 快速启动

```bash
# 前置：MySQL + Redis + etcd
make start
```

## 服务架构

```
用户端 API                管理端 API
user-api    :8888         admin-api   :8889
event-api   :8890
inventory-api :8891
ticket-api  :8892
order-api   :8894

RPC 服务
event-rpc     :8080       order-rpc     :8083
ticket-rpc   :8081        inventory-rpc :8084
payment-rpc  :8082
```

## 核心流程

```
POST /api/v1/order/buy
  → EventRpc.GetTicketType（查票价）
  → InventoryRpc.DeductStock（Redis Lua 原子扣库存）
  → OrderRpc.CreateOrder（MySQL 建订单）
  → PaymentRpc.Pay（支付）
  → TicketRpc.CreateTicket（出票）
```

支付失败自动回滚库存并取消订单，超时 15 分钟未支付定时取消。

## 数据模型

### User

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | varchar(6) | 用户 ID（6 位随机） |
| Username | varchar(64) | 用户名，唯一 |
| Password | varchar(255) | bcrypt 加密 |
| Email | varchar(128) | 邮箱 |
| Phone | varchar(20) | 手机号 |
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
| Status | varchar(32) | draft / ready / selling / closed |
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
| Stock | int32 | MySQL 库存（Redis 为真实库存来源） |
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
| Status | varchar(32) | reserved → paid / failed / cancelled |

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
| TransferStatus | varchar(32) | 转赠状态 |

### Redis 库存 Key

```
stock:ticket:{ticketTypeId}   →  整数（Lua 原子扣减）
```

---

继续完善中。
