## Redis Lua 脚本

| 脚本 | 作用 |
|------|------|
| `DeductStockLua` | GET + 校验 + DECRBY 原子执行，保证 0 超卖 |
| `RateLimitLua` | 用户固定窗口限流 |
| `TokenBucketLua` | 票种令牌桶平滑限流 |
| `UnlockLua` | 比对 value 后释放，防误删 |
### Redis Key

```
stock:ticket:{ticketTypeId}                           # 库存
lock:order:{userId}:{eventId}:{showId}:{ticketTypeId} # 分布式锁
idem:{userId}:{eventId}:{showId}:{ticketTypeId}:{requestId} # 幂等
limit:buy:user:{userId}                               # 用户限流
bucket:ticket:{ticketTypeId}                          # 令牌桶
compensate:list                                       # 补偿记录列表