
### MySQL 核心表

- **User** — bcrypt 密码，user/admin 双角色
- **Event** — draft → ready → selling → closed
- **Show** — 场次
- **TicketType** — 票种（价格 / 库存 / 限购）
- **Order** — unpaid → paid / cancelled
- **Ticket** — 出票记录
## 高可用设计

| 层 | 方案 | 故障影响 |
|---|---|---|
| API 层 | Nginx upstream 3× order-api | 1 个实例挂，自动轮询另外两个 |
| RPC 层 | 5 服务各 2 实例 + etcd 客户端 lb | 1 个实例挂，etcd 踢掉，流量切另一个 |
| Redis | 1 Master + 1 Slave + 3 Sentinel | Master 挂，Sentinel 10s 内提 Slave 为新 Master |
| 熔断 | go-zero 内置 breaker，错误率 >50% 自动熔断 | 下游挂了秒级熔断，不拖垮上游 |
| 降载 | go-zero 自适应降载，CPU 过高时拒绝请求 | 过载时返回 503，保护服务不崩 |


## 中间件要点

- **Redis Pipeline** — 限流 + 令牌桶 + 幂等合并为 1 次 Redis 往返
- **RabbitMQ** — 扣库存后异步建订单，API 层只等 Redis（<1ms）
- **死信队列** — 15 分钟 TTL 自动取消未支付订单，回滚库存
- **补偿机制** — 回滚失败写补偿记录，后台定时自动执行
- **Nginx 多实例** — order-api 3 实例 upstream 轮询
- **JWT 缓存** — 已验证 token 缓存 1 分钟，跳过重复 HMAC 计算
- **Redis Sentinel** — go-redis FailoverClient 自动发现 Master，故障转移 10s 恢复
- **OpenTelemetry** — 全链路追踪，Jaeger 可视化
