# 通用 Dockerfile — 一个文件打所有微服务
# 用法:
#   docker build --build-arg APP_PATH=app/order/order-api --build-arg CONF_FILE=etc/order-api.yaml -t ticket:order-api .
#   docker build --build-arg APP_PATH=app/order/order-rpc --build-arg CONF_FILE=etc/order.yaml       -t ticket:order-rpc .

# ============================================================
# 阶段 1：编译
# ============================================================
FROM golang:1.25-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 先拷依赖描述文件，利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

# 再拷源码
COPY . .

ARG APP_PATH
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./${APP_PATH}

# ============================================================
# 阶段 2：运行
# ============================================================
FROM alpine:3.21

RUN apk add --no-cache tzdata ca-certificates
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/server .

# ARG 不能跨 FROM，阶段 2 必须重新声明
ARG APP_PATH=app/order/order-api
ARG CONF_FILE=etc/order-api.yaml
COPY ${APP_PATH}/${CONF_FILE} ./config.yaml

EXPOSE 8080

ENTRYPOINT ["sh", "-c", "sed -i \"s/127.0.0.1:3306/${MYSQL_HOST:-127.0.0.1}:3306/g; s/127.0.0.1:2379/${ETCD_HOST:-127.0.0.1}:2379/g; s/127.0.0.1:26379/${REDIS_S1_HOST:-127.0.0.1}:26379/g; s/127.0.0.1:26380/${REDIS_S2_HOST:-127.0.0.1}:26380/g; s/127.0.0.1:26381/${REDIS_S3_HOST:-127.0.0.1}:26381/g; s/127.0.0.1:4317/${OTEL_HOST:-127.0.0.1}:4317/g\" config.yaml && exec ./server -f config.yaml"]
