.PHONY: start-infra start-rpc start-api start-all stop build clean

ROOT := $(shell pwd)

# ======== 基础设施 ========

start-infra:
	@echo "start etcd..."
	@etcd > /dev/null 2>&1 &
	@sleep 1
	@echo "start redis..."
	@redis-server --daemonize yes 2>/dev/null; true
	@echo "infra ready"

# ======== 启动 RPC ========

start-rpc:
	@echo "start event-rpc..."
	cd $(ROOT)/app/event/event-rpc && go run event.go &
	@sleep 1
	@echo "start inventory-rpc..."
	cd $(ROOT)/app/inventory/inventory-rpc && go run inventory.go &
	@sleep 1
	@echo "start order-rpc..."
	cd $(ROOT)/app/order/order-rpc && go run order.go &
	@sleep 1
	@echo "start payment-rpc..."
	cd $(ROOT)/app/payment/pay-rpc && go run payment.go &
	@sleep 1
	@echo "start ticket-rpc..."
	cd $(ROOT)/app/ticket/ticket-rpc && go run ticket.go &

# ======== 启动 API ========

start-api:
	@echo "start user-api..."
	cd $(ROOT)/app/user && go run user.go &
	@sleep 0.5
	@echo "start event-api..."
	cd $(ROOT)/app/event/event-api && go run event.go &
	@sleep 0.5
	@echo "start inventory-api..."
	cd $(ROOT)/app/inventory/inventory-api && go run inventory.go &
	@sleep 0.5
	@echo "start ticket-api..."
	cd $(ROOT)/app/ticket/ticket-api && go run ticket.go &
	@sleep 0.5
	@echo "start order-api..."
	cd $(ROOT)/app/order/order-api && go run order.go &
	@sleep 0.5
	@echo "start admin..."
	cd $(ROOT)/app/admin && go run admin.go &

# ======== 一键 ========

start: start-infra
	@sleep 2
	@$(MAKE) start-rpc
	@sleep 6
	@$(MAKE) start-api
	@echo "=== 全部启动完成 ==="

# ======== 停止 ========

stop:
	@echo "killing all..."
	@fuser -k 8888/tcp 8889/tcp 8890/tcp 8891/tcp 8892/tcp 8894/tcp 8080/tcp 8081/tcp 8082/tcp 8083/tcp 8084/tcp 2>/dev/null; true
	@echo "done"

# ======== 其他 ========

build:
	go build -o /dev/null ./...

clean: stop
