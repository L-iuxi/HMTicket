# 测试手册

## 自动化测试（Go）

| 工具 | 测试内容 | 验证点 |
|------|----------|--------|
| `go run test/idem/main.go` | 幂等 + 分布式锁 | 同 requestId 同结果、5 并发锁串行 |
| `go run test/rate/main.go` | 三层限流 | 用户窗口 5/s、令牌桶、支付限流 |
| `go run test/e2e/main.go` | 端到端全链路 | 创建活动 → 买票 → MQ 建单 → 支付 → 出票 |

每个自动注册独立用户，可重复运行。

## 压测工具（Go）

```bash
go run test/bench/main.go -c 100 -d 10s
```

参数: `-c` 并发数 `-d` 持续时间 `-token` JWT（可选）

每个 goroutine 生成随机 requestId 打 buy 接口，每秒报告 QPS。

## 手动 curl

| 接口 | 方法 | 地址 | 说明 |
|------|------|------|------|
| 注册 | POST | `:8888/api/user/register` | `{username,password,email,phone,gender}` |
| 登录 | POST | `:8888/api/user/login` | `{account,password}` 返回 token |
| 个人信息 | GET | `:8888/api/user/profile` | Header: `Authorization: Bearer <token>` |
| 创建活动 | POST | `:8889/admin/event` | `{title,description,location,...}` |
| 创建场次 | POST | `:8889/admin/show` | `{eventId,name,showTime,endTime,sortOrder}` |
| 创建票种 | POST | `:8889/admin/ticket-type` | `{eventId,showId,name,price,stock,maxPerUser}` |
| 发布活动 | PUT | `:8889/admin/event/status` | `{eventId,status:"selling"}` |
| 查看活动 | GET | `:8890/event/{id}` | 公开 |
| 查看票种 | GET | `:8890/ticket-type/{id}` | 公开 |
| 初始化库存 | PUT | `:8891/api/v1/admin/inventory/stock` | `{ticketTypeId,stock}` |
| 查库存 | GET | `:8891/api/v1/inventory/stock/{id}` | |
| 买票 | POST | `:8894/api/v1/order/buy` | `{eventId,showId,ticketTypeId,quantity,requestId}` |
| 订单列表 | GET | `:8894/api/v1/order/list` | 需 token |
| 订单详情 | GET | `:8894/api/v1/order/{orderNo}` | 需 token |
| 取消订单 | POST | `:8894/api/v1/order/cancel` | `{orderNo}` |
| 支付 | POST | `:8894/api/v1/order/pay` | `{orderNo,requestId}` |
| 我的票 | GET | `:8892/api/ticket/list` | 需 token |
| 验票 | POST | `:8892/api/ticket/check` | `{ticketId}` |
| 转赠 | POST | `:8892/api/ticket/transfer` | `{ticketId,toUserId}` |
