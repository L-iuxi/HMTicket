## 测试手册
此处测试为基本接口测试，
限流测试——ratelimit_test.sh
幂等、分布式锁测试——idempotent_test.sh
```
go run test/idem/main.go
go run test/rate/main.go
go run test/e2e/main.go
go run test/bench/main.go -c 200 -n 50

```
### 前置安装
需要redis-server，mysql，etcd，在终端执行如下命令安装：
```
sudo apt install mysql-server -y//mysql密码配置可修改
sudo apt install redis-server -y
sudo apt install etcd-server etcd-client -y
```
### 启动服务

  make start-rpc   # 先起 RPC（等注册到 etcd）
  make start-api   # 再起 API
  make start       # 一键全起
  make stop        # 杀掉全部
  make build       # 检查编译

### user包
**注册：**
curl -X POST http://localhost:8888/api/user/register \
-H "Content-Type: application/json" \
-d '{
  "username":"dxyxx",
  "password":"lyx188288",
  "email":"3301919119@qq.com",
  "phone":"13811111111",
  "gender":1
}'

1. 电话，用户名，邮箱重复注册返回提示

**登陆：**
curl -X POST http://localhost:8888/api/user/login \
-H "Content-Type: application/json" \
-d '{
  "account":"admin",
  "password":"123456"
}'
1. 账号密码错误返回提示
   
**查询个人信息**
curl http://localhost:8888/api/user/profile \
-H "Authorization: Bearer 你的token"

**修改个人信息**
curl -X PUT http://localhost:8888/api/user/profile \
-H "Authorization: Bearer token" \
-H "Content-Type: application/json" \
-d '{
    "phone":"13912345678"
}'

### admin包
**管理员添创建活动**
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
1. 同名称同活动校验
2. 输入校验
**管理员创建场次**
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
**管理员修改活动**
  Ticket git:(main) ✗ curl -X PUT http://localhost:8889/admin/event \
-H "Authorization: Bearer token" \
-H "Content-Type: application/json" \
-d '{
    "eventId":2,
    "endTime":"2026-02-20 22:29:00"
}'
**管理员修改活动状态**
curl -X PUT http://localhost:8889/admin/event/status \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "eventId":1,
    "status":"online"
}'  
1. 对输入状态有判断 
**管理员修改场次**
curl -X PUT http://localhost:8889/admin/show \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "showId":1,
    "name":"第一场（加场）",
    "sortOrder":2
}'

**管理员创建票种**
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
**管理员修改票种**
curl -X PUT http://localhost:8889/admin/ticket-type \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "ticketTypeId":1,
    "price":1380,
    "stock":80,
    "maxPerUser":4
}'

### event包
**查看某活动**
curl http://localhost:8890/event/1 
**查看某场次**
curl http://localhost:8890/show/1
**查询某票种**
curl http://localhost:8890/ticket-type/1

### inventory包
**查看库存**
curl http://localhost:8891/api/v1/inventory/stock/1
**修改库存**

curl -X PUT http://localhost:8891/api/v1/admin/inventory/stock \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "ticketTypeId":1,
    "stock":200
}'
### payment包
**支付订单**
curl -X POST \
    http://127.0.0.1:8894/api/v1/order/pay \
    -H "Authorization: Bearer Token" \
    -H "Content-Type: application/json" \
    -d '{
      "orderNo":"",
      "requestId":"1"                                                     
  }'

### order包
**买票**
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
**查看某人全部订单**
➜  ~ curl -X GET \
http://127.0.0.1:8894/api/v1/order/list \
-H "Authorization: Bearer Token"

**查看某个订单**
➜  ~ curl -X GET \
http://127.0.0.1:8894/api/v1/order/3740c823-ca42-4400-9af0-cc7971848b98 \
-H "Authorization: Bearer Token"

**取消订单**
➜  ~ curl -X POST \                      
http://127.0.0.1:8894/api/v1/order/cancel \
-H "Authorization: Bearer Token" \
-H "Content-Type: application/json" \
-d '{
    "orderNo":""
}'
**管理员删除订单**
➜  ~ curl -X DELETE \
http://127.0.0.1:8894/api/v1/admin/order/OrderNo \
-H "Authorization: Bearer Token"
**管理员修改订单**
➜  ~ curl -X PUT \                       
http://127.0.0.1:8894/api/v1/admin/order/update \
-H "Authorization: Bearer Token" \
-H "Content-Type: application/json" \
-d '{
    "orderNo":"",
    "quantity":3
}'

### ticket包
**查看我的所有票**
➜  ~ curl -X GET \
http://127.0.0.1:8892/api/ticket/list \
-H "Authorization: Bearer Token"
**查看票详情**
➜  ~ curl -X GET \
http://127.0.0.1:8892/api/ticket/detail/2 \
-H "Authorization: Bearer Token"
**验票**
➜  ~ curl -X POST \
http://127.0.0.1:8892/api/ticket/check \
-H "Authorization: Bearer Token" \
-H "Content-Type: application/json" \
-d '
{
    "ticketId":1
}'
**转赠**
