## 测试手册

### 前置安装
需要redis-server，mysql，etcd，rabbitMQ,nginx在终端执行如下命令安装：
```bash
# 安装mysql
sudo apt install mysql-server -y//mysql密码配置可修改

# 安装redis
sudo apt install redis-server -y# 

# 安装etcd
sudo apt install etcd-server etcd-client -y

# 启动rabbitmq镜像
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=123456 \
  rabbitmq:3-management

# Nginx
# 安装
sudo apt install nginx
# 把项目配置软链过去
sudo ln -sf /home/Ticket/nginx/nginx.conf /etc/nginx/nginx.conf
# 启动
sudo nginx
# 重载（修改配置后）
sudo nginx -s reload
# 停止
sudo nginx -s stop
```

### 启动服务
```bash
  make start-rpc   # 先起 RPC（等注册到 etcd）
  make start-api   # 再起 API
  make start       # 一键全起
  make stop        # 杀掉全部
  make build       # 检查编译
```
### Jaeger 链路追踪

```bash
docker run -d --name jaeger -p 4317:4317 -p 16686:16686 jaegertracing/all-in-one:latest
```

浏览器打开 `http://127.0.0.1:16686` 查看调用链。
### 压测

跑压测需要一个管理员账号

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
**压力测试**
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


测试结果报告见[bench.md](bench.md)。

### 普通测试
*以下测试在本机测试端口为各服务端口，如果在容器内测试统一端口为nginx端口：8090*
### user包
**注册：**
```bash
curl -X POST http://localhost:8888/api/user/register \
-H "Content-Type: application/json" \
-d '{
  "username":"dxyxx",
  "password":"lyx188288",
  "email":"3301919119@qq.com",
  "phone":"13811111111",
  "gender":1
}'
```

1. 电话，用户名，邮箱重复注册返回提示
```bash
**登陆：**
curl -X POST http://localhost:8888/api/user/login \
-H "Content-Type: application/json" \
-d '{
  "account":"admin",
  "password":"123456"
}'
```
1. 账号密码错误返回提示

**查询个人信息**
```bash
curl http://localhost:8888/api/user/profile \
-H "Authorization: Bearer 你的token"
```

**修改个人信息**
```bash
curl -X PUT http://localhost:8888/api/user/profile \
-H "Authorization: Bearer token" \
-H "Content-Type: application/json" \
-d '{
    "phone":"13912345678"
}'
```
### admin包
**管理员添创建活动**
```bash
curl -X POST http://localhost:8889/admin/event \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "title":"五月天演唱会",
    "description":"2026巡演",
    "location":"上海体育馆",
    "coverImage":"xxx.jpg",
    "startTime":"2026-08-01 19:30:00",
    "endTime":"2026-08-01 22:30:00"
}'
```
1. 同名称同活动校验
2. 输入校验
**管理员创建场次**
```bash
curl -X POST http://localhost:8889/admin/show \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "eventId":1,
    "name":"第一场",
    "showTime":"2026-08-24 19:30:00",
    "endTime":"2026-08-24 22:30:00",
    "sortOrder":1
}'
```
**管理员修改活动**
```bash
  Ticket git:(main) ✗ curl -X PUT http://localhost:8889/admin/event \
-H "Authorization: Bearer token" \
-H "Content-Type: application/json" \
-d '{
    "eventId":2,
    "endTime":"2026-02-20 22:29:00"
}'
```
**管理员修改活动状态**
```bash
curl -X PUT http://localhost:8889/admin/event/status \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "eventId":1,
    "status":"online"
}'  
```
1. 对输入状态有判断 
**管理员修改场次**
```bash
curl -X PUT http://localhost:8889/admin/show \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "showId":1,
    "name":"第一场（加场）",
    "sortOrder":2
}'
```
**管理员创建票种**
```bash
curl -X POST http://localhost:8889/admin/ticket-type \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "eventId":1,
    "showId":1,
    "name":"VIP票",
    "price":1280,
    "stock":100,
    "maxPerUser":2,
    "sortOrder":1
}'
```
**管理员修改票种**
```bash
curl -X PUT http://localhost:8889/admin/ticket-type \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "ticketTypeId":1,
    "price":1380,
    "stock":80,
    "maxPerUser":4
}'
```
### event包
**查看某活动**
```bash
curl http://localhost:8890/event/1 
```
**查看某场次**
```bash
curl http://localhost:8890/show/1
```
**查询某票种**
```bash
curl http://localhost:8890/ticket-type/1
```
### inventory包
**查看库存**
```bash
curl http://localhost:8891/api/v1/inventory/stock/1
```
**修改库存**
```bash
curl -X PUT http://localhost:8891/api/v1/admin/inventory/stock \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "ticketTypeId":1,
    "stock":200
}'
```
### payment包
**支付订单**
```bash
curl -X POST \
    http://127.0.0.1:8894/api/v1/order/pay \
    -H "Authorization: Bearer Token" \
    -H "Content-Type: application/json" \
    -d '{
      "orderNo":"",
      "requestId":"1"                                                     
  }'
```
### order包
**买票**
```bash
curl -X POST \
http://127.0.0.1:8894/api/v1/order/buy \
-H "Authorization: Bearer Token" 
\-H "Content-Type: application/json" \
-d '{
    "eventId":1,
    "showId":1,
    "ticketTypeId":1,
    "quantity":1
    "requestId":1
}'
```
**查看某人全部订单**
```bash
curl -X GET \
http://127.0.0.1:8894/api/v1/order/list \
-H "Authorization: Bearer Token"
```
**查看某个订单**
```bash
curl -X GET \
http://127.0.0.1:8894/api/v1/order/3740c823-ca42-4400-9af0-cc7971848b98 \
-H "Authorization: Bearer Token"
```
**取消订单**
```bash
curl -X POST \                      
http://127.0.0.1:8894/api/v1/order/cancel \
-H "Authorization: Bearer Token" \
-H "Content-Type: application/json" \
-d '{
    "orderNo":""
}'
```
**管理员删除订单**
```bash
curl -X DELETE \
http://127.0.0.1:8894/api/v1/admin/order/OrderNo \
-H "Authorization: Bearer Token"
```
**管理员修改订单**
```bash
curl -X PUT \                       
http://127.0.0.1:8894/api/v1/admin/order/update \
-H "Authorization: Bearer Token" \
-H "Content-Type: application/json" \
-d '{
    "orderNo":"",
    "quantity":3
}'
```
### ticket包
**查看我的所有票**
```bash
curl -X GET \
http://127.0.0.1:8892/api/ticket/list \
-H "Authorization: Bearer Token"
```
**查看票详情**
```bash
curl -X GET \
http://127.0.0.1:8892/api/ticket/detail/2 \
-H "Authorization: Bearer Token"
```
**验票**
```bash
curl -X POST \
http://127.0.0.1:8892/api/ticket/check \
-H "Authorization: Bearer Token" \
-H "Content-Type: application/json" \
-d '
{
    "ticketId":1
}'

```
## 压测（bench）

### 用法

```bash
# 直连 order-api（单实例）
go run test/bench/main.go -c 200 -n 50

# raw 模式 — 去掉 100ms 重试间隔，打满 QPS
go run test/bench/main.go -c 500 -n 100 -raw

# 走 Nginx — -addr 指定 Nginx 地址
go run test/bench/main.go -c 500 -n 100 -raw -addr 127.0.0.1:8090

# 多客户端同时打 Nginx
go run test/bench/main.go -c 300 -n 50 -raw -addr 127.0.0.1:8090 &
go run test/bench/main.go -c 300 -n 50 -raw -addr 127.0.0.1:8090 &
go run test/bench/main.go -c 300 -n 50 -raw -addr 127.0.0.1:8090 &
wait
```

### 参数

| 参数 | 默认 | 说明 |
|------|------|------|
| `-c` | 200 | 并发用户数 |
| `-n` | 50 | 库存票数 |
| `-t` | 30s | 抢票超时 |
| `-raw` | false | 去 sleep，打满 QPS |
| `-addr` | 空 | Nginx 地址，如 `127.0.0.1:8090`。留空直连各服务端口 |

### 流程

自动完成：admin 登录 → 创建活动/场次/票种 → 注册用户 → 并发抢票 → 验证 0 超卖。

### 单机基准数据

```
500 人抢 100 张 | raw 模式 | 单实例直连
总请求 917,090 | QPS 30,566 | 卖出 100 | 超卖 0

3×300 人抢 50 张 | raw | Nginx + 3 实例
总请求 928,124 | QPS 30,909 | 卖出 150 | 超卖 0
```


