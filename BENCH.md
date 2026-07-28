# 压测报告

## 测试环境

| 项目 | 配置 |
|------|------|
| CPU | Intel i7-12700H (14 核 20 线程) |
| 内存 | 16 GB |
| OS | Ubuntu 22.04 |
| Go | 1.24 |
| 部署 | 所有服务 + MySQL + Redis + RabbitMQ + etcd 同机运行 |
| 工具 | `go run test/bench/main.go` |

## 测试结果

| 场景 | 配置 | 总请求 | QPS | 卖出 | 超卖 |
|------|------|--------|-----|------|------|
| 200 人抢 50 张 | 单实例 | 902,860 | 30,091 | 50 | 0 |
| 500 人抢 100 张 | 单实例 | 917,090 | 30,566 | 100 | 0 |
| 500 人抢 200 张 | 单实例 | 893,920 | 29,793 | 200 | 0 |
| 3×300 人抢 50 张 | Nginx + 3 实例 | 928,124 | 30,909 | 150 | 0 |
| 3×300 人抢 50 张 | JWT 缓存 + keepalive | 928,556 | 30,952 | 150 | 0 |

## 保护层开销

raw 模式 90 万请求中，保护层拦截分布：

| 保护层 | 技术 | 拦截量 |
|--------|------|--------|
| 用户级限流 | Redis Lua 固定窗口 5次/秒/人 | ~45 万 |
| 票种令牌桶 | Redis Lua 令牌桶 200 tokens + 100/s | ~44 万 |
| 分布式锁 | Redis SETNX 用户+票种粒度 | ~1 万 |
| 幂等 | requestId Redis idem key | 防重不防多 |

## CPU Profile 分析

30 秒采样，pprof 数据：

| 组件 | CPU 占比 |
|------|---------|
| net/http conn.serve | 48% |
| JWT Authorize (HMAC-SHA256) | 38% |
| TraceHandler (OpenTelemetry) | 24% |
| Redis go-redis | 11% |
| JSON 编解码 | 6% |

瓶颈在 HTTP 框架和中间件，不在业务逻辑。

## 优化记录

| 优化 | 效果 |
|------|------|
| Redis Pipeline（限流+令牌桶+锁合并） | 减少 2 次 Redis 往返，延迟降 0.09ms |
| JWT 内存缓存（sync.Map，1 分钟 TTL） | 同用户重复请求跳过 HMAC |
| Nginx keepalive（HTTP/1.1 + 64 连接池） | 单机瓶颈不在 Nginx |
| 关 Telemetry | 无明显提升（瓶颈在 net/http 本身） |

## 结论

- 单机 **30k QPS**，go-zero + net/http 栈物理上限
- 0 超卖，所有场景验证通过
- 架构支持水平扩展，加机器 N × 30k
- 多实例 + Nginx 配置就位，部署即用
