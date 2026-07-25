package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/common/mq"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// StartOrderConsumer 启动 MQ 消费者，监听 order.create 队列
func StartOrderConsumer(svcCtx *svc.ServiceContext) error {
	if svcCtx.MQ == nil {
		logx.Error("[order-rpc] MQ 未连接，跳过消费者启动")
		return nil
	}

	// 初始化 exchange + queue + bind
	if err := svcCtx.MQ.InitExchange(); err != nil {
		return fmt.Errorf("init exchange: %w", err)
	}
	if err := svcCtx.MQ.InitDeadExchange(); err != nil {
		return fmt.Errorf("init dead exchange: %w", err)
	}

	queueName := "order.create.queue"
	if err := svcCtx.MQ.InitQueue(queueName); err != nil {
		return fmt.Errorf("init queue: %w", err)
	}
	if err := svcCtx.MQ.InitDeadQueue(); err != nil {
		return fmt.Errorf("init dead queue: %w", err)
	}
	if err := svcCtx.MQ.InitQueue("order.timeout.queue"); err != nil {
		return fmt.Errorf("init timeout queue: %w", err)
	}
	if err := svcCtx.MQ.BindQueue(queueName, mq.RoutingOrderCreate); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}
	if err := svcCtx.MQ.BindQueue("order.timeout.queue", mq.RoutingTimeOutCancel); err != nil {
		return fmt.Errorf("bind timeout queue: %w", err)
	}
	if err := svcCtx.MQ.BindDeadQueue(); err != nil {
		return fmt.Errorf("bind dead queue: %w", err)
	}

	handler := func(ctx context.Context, body []byte) error {
		return handleCreateOrder(ctx, svcCtx, body)
	}

	logx.Info("[order-rpc] MQ consumer started, listening on order.create.queue")
	return svcCtx.MQ.Consume(queueName, handler)
}

func handleCreateOrder(ctx context.Context, svcCtx *svc.ServiceContext, body []byte) error {
	var msg mq.CreateOrderMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logx.Errorf("[MQ] 消息反序列化失败: %v, body=%s", err, string(body))
		return err
	}

	order := db.Order{
		UserID:       uint(msg.UserID),
		EventID:      msg.EventID,
		ShowID:       uint(msg.ShowID),
		TicketTypeID: uint(msg.TicketTypeID),
		OrderNo:      msg.OrderNo,
		Quantity:     int(msg.Quantity),
		TotalPrice:   msg.TotalPrice,
		Status:       "unpaid",
	}

	if err := svcCtx.DB.Create(&order).Error; err != nil {
		logx.Errorf("[MQ] 创建订单失败: orderNo=%s, err=%v", msg.OrderNo, err)
		return err
	}

	// 标记幂等成功，TTL 15 分钟，15min之后幂等建过期
	if msg.IdemKey != "" && svcCtx.Idempotent != nil {
		if err := svcCtx.Idempotent.Success(ctx, msg.IdemKey, msg.OrderNo, 15*60*time.Second); err != nil {
			logx.Errorf("[MQ] 幂等标记失败: idemKey=%s, err=%v", msg.IdemKey, err)
		}
	}

	//加入支付队列
	timeoutMsg := mq.TimeOutCancelMessage{
		OrderNo:      msg.OrderNo,
		Quantity:     msg.Quantity,
		TicketTypeID: msg.TicketTypeID,
	}
	body, err := json.Marshal(timeoutMsg)
	if err != nil {
		logx.Errorf("[MQ] 消息序列化失败")
		return err
	}
	if err := svcCtx.MQ.Ch.PublishWithContext(ctx,
		mq.TicketExchange,
		mq.RoutingTimeOutCancel,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Expiration:  "900000", //15min
		}); err != nil {
		logx.Errorf("[MQ] 超时消息发送失败: orderNo=%s, err=%v", msg.OrderNo, err)
	}
	logx.Infof("[MQ] 订单创建成功: orderNo=%s, userId=%d, ticketTypeId=%d",
		msg.OrderNo, msg.UserID, msg.TicketTypeID)

	return nil
}

// StartTimeoutConsumer 启动延迟超时消费者（死信队列）
// 订单 15 分钟后进入死信 → 消费者检查是否未支付 → 取消订单+回滚库存
func StartTimeoutConsumer(svcCtx *svc.ServiceContext) error {

	if svcCtx.MQ == nil {
		logx.Error("[order-rpc] MQ 未连接，跳过消费者启动")
		return nil
	}

	// 初始化 exchange + queue + bind
	if err := svcCtx.MQ.InitExchange(); err != nil {
		return fmt.Errorf("init exchange: %w", err)
	}
	if err := svcCtx.MQ.InitDeadExchange(); err != nil {
		return fmt.Errorf("init dead exchange: %w", err)
	}

	if err := svcCtx.MQ.InitDeadQueue(); err != nil {
		return fmt.Errorf("init dead queue: %w", err)
	}

	if err := svcCtx.MQ.BindDeadQueue(); err != nil {
		return fmt.Errorf("bind dead queue: %w", err)
	}

	handler := func(ctx context.Context, body []byte) error {
		return handleTimeOutCancel(ctx, svcCtx, body)
	}

	logx.Info("[order-rpc] timeout consumer started, listening on dead_queue")
	return svcCtx.MQ.Consume(mq.DeadQueue, handler)
}

func handleTimeOutCancel(ctx context.Context, svcCtx *svc.ServiceContext, body []byte) error {
	var msg mq.TimeOutCancelMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logx.Errorf("[MQ] 消息反序列化失败: %v, body=%s", err, string(body))
		return err
	}
	var order db.Order
	if err := svcCtx.DB.Where("order_no = ?", msg.OrderNo).First(&order).Error; err != nil {
		logx.Errorf("[MQ] 订单不存在: %s", msg.OrderNo)
		return err
	}

	if order.Status != "unpaid" {
		return nil
	}
	//回滚库存
	_, err := svcCtx.InventoryRpc.ReleaseStock(ctx, &inventoryclient.ReleaseStockReq{
		TicketTypeId: uint64(order.TicketTypeID),
		Quantity:     int32(order.Quantity),
	})
	if err != nil {
		logx.Errorf("超时取消回滚库存失败: orderNo=%s, ticketTypeId=%d, err=%v",
			order.OrderNo, order.TicketTypeID, err)
		return err
	}
	//订单标记为已取消
	svcCtx.DB.Model(&order).Update("status", db.OrderCanceled)

	logx.Infof("[MQ] 订单超时未支付已取消: %s", msg.OrderNo)
	return nil
}
