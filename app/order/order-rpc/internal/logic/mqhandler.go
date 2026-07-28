package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/common/mq"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// StartOrderConsumer 监听 order.create.queue，异步建订单
func StartOrderConsumer(svcCtx *svc.ServiceContext) error {
	if svcCtx.MQ == nil {
		logx.Error("[order-rpc] MQ 未连接，跳过消费者启动")
		return nil
	}

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

	if msg.IdemKey != "" && svcCtx.Idempotent != nil {
		val, _ := svcCtx.Idempotent.Get(ctx, msg.IdemKey)
		if strings.HasPrefix(val, "ok") {
			return nil //已幂等成功，由于各种原因来到队列
		}
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

	// 标记幂等成功
	if msg.IdemKey != "" && svcCtx.Idempotent != nil {
		var idemErr error
		for i := 0; i < 3; i++ {
			idemErr = svcCtx.Idempotent.Success(ctx, msg.IdemKey, msg.OrderNo, 15*time.Minute)
			if idemErr == nil {
				break
			}
			logx.Errorf("[MQ] 幂等标记失败(第%d次): idemKey=%s, err=%v", i+1, msg.IdemKey, idemErr)
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
		}
		if idemErr != nil {
			return idemErr
		}
	}

	// 发超时消息
	timeoutMsg := mq.TimeOutCancelMessage{
		OrderNo:      msg.OrderNo,
		Quantity:     msg.Quantity,
		TicketTypeID: msg.TicketTypeID,
	}
	timeoutBody, err := json.Marshal(timeoutMsg)
	if err != nil {
		logx.Errorf("[MQ] 超时消息序列化失败")
		return nil
	}

	var timeoutErr error
	for i := 0; i < 3; i++ {
		timeoutErr = svcCtx.MQ.Ch.PublishWithContext(ctx,
			mq.TicketExchange,
			mq.RoutingTimeOutCancel,
			false, false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        timeoutBody,
				Expiration:  "900000",
			})
		if timeoutErr == nil {
			break
		}
		logx.Errorf("[MQ] 超时消息发送失败(第%d次): orderNo=%s, err=%v", i+1, msg.OrderNo, timeoutErr)
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	if timeoutErr != nil {
		return timeoutErr
	}

	logx.Infof("[MQ] 订单创建成功: orderNo=%s, userId=%d, ticketTypeId=%d",
		msg.OrderNo, msg.UserID, msg.TicketTypeID)

	return nil
}

// StartDeadConsumer 监听 dead_queue。
// 两种消息会到这里：1) 建订单 3 次重试失败  2) 超时取消 TTL 到期
func StartDeadConsumer(svcCtx *svc.ServiceContext) error {
	if svcCtx.MQ == nil {
		logx.Error("[order-rpc] MQ 未连接，跳过死信消费者启动")
		return nil
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
		// 尝试按 CreateOrderMessage 解析（有 idemKey 说明是建订单失败）
		var createMsg mq.CreateOrderMessage
		if err := json.Unmarshal(body, &createMsg); err == nil && createMsg.IdemKey != "" {
			return handleDeadOrderCreate(ctx, svcCtx, &createMsg)
		}

		// 否则是超时取消消息
		return handleTimeOutCancel(ctx, svcCtx, body)
	}

	logx.Info("[order-rpc] dead consumer started, listening on dead_queue")
	return svcCtx.MQ.Consume(mq.DeadQueue, handler)
}

// 建订单死信
func handleDeadOrderCreate(ctx context.Context, svcCtx *svc.ServiceContext, msg *mq.CreateOrderMessage) error {
	logx.Errorf("[DEAD] 建订单消息进死信: orderNo=%s, idemKey=%s", msg.OrderNo, msg.IdemKey)

	//  先查 MySQL
	var existingOrder db.Order
	err := svcCtx.DB.Where("order_no = ?", msg.OrderNo).First(&existingOrder).Error

	if err == nil {
		// 订单存在，补幂等标记
		logx.Infof("[DEAD] 订单已存在，补幂等标记: orderNo=%s, idemKey=%s", msg.OrderNo, msg.IdemKey)
		if msg.IdemKey != "" && svcCtx.Idempotent != nil {
			if err := svcCtx.Idempotent.Success(ctx, msg.IdemKey, msg.OrderNo, 15*time.Minute); err != nil {
				logx.Errorf("[DEAD] 补幂等标记失败: idemKey=%s, err=%v", msg.IdemKey, err)
			}
		}

		// 补发超时消息（
		timeoutMsg := mq.TimeOutCancelMessage{
			OrderNo:      msg.OrderNo,
			Quantity:     msg.Quantity,
			TicketTypeID: msg.TicketTypeID,
		}
		timeoutBody, _ := json.Marshal(timeoutMsg)
		for i := 0; i < 3; i++ {
			err := svcCtx.MQ.Ch.PublishWithContext(ctx,
				mq.TicketExchange,
				mq.RoutingTimeOutCancel,
				false, false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        timeoutBody,
					Expiration:  "900000",
				})
			if err == nil {
				logx.Infof("[DEAD] 超时消息补发成功: orderNo=%s", msg.OrderNo)
				break
			}
			logx.Errorf("[DEAD] 超时消息补发失败(第%d次): orderNo=%s, err=%v", i+1, msg.OrderNo, err)
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
		}

		return nil
	}

	// 订单不存在
	logx.Errorf("[DEAD] 订单不存在，确认创建失败: orderNo=%s", msg.OrderNo)

	if msg.IdemKey != "" && svcCtx.Idempotent != nil {
		if err := svcCtx.Idempotent.Failed(ctx, msg.IdemKey, 10*time.Minute); err != nil {
			logx.Errorf("[DEAD] 标记幂等失败: idemKey=%s, err=%v", msg.IdemKey, err)
		}
	}

	// 防重复回滚：SETNX 标记，只允许执行一次
	releaseKey := "release:dead:" + msg.OrderNo
	ok, setErr := svcCtx.Redis.Client().SetNX(ctx, releaseKey, "1", 30*time.Minute).Result()
	if setErr != nil {
		logx.Errorf("[DEAD] 释放标记设置失败: orderNo=%s, err=%v", msg.OrderNo, setErr)
		return setErr
	}
	if !ok {
		logx.Infof("[DEAD] 库存已释放过，跳过: orderNo=%s", msg.OrderNo)
		return nil
	}

	// 回滚库存
	stockKey := fmt.Sprintf("stock:ticket:%d", msg.TicketTypeID)
	for i := 0; i < 3; i++ {
		_, err := svcCtx.Redis.IncrBy(ctx, stockKey, int64(msg.Quantity))
		if err == nil {
			logx.Infof("[DEAD] 库存回滚成功: ticketTypeId=%d, quantity=%d", msg.TicketTypeID, msg.Quantity)
			return nil
		}
		logx.Errorf("[DEAD] 库存回滚失败(第%d次): ticketTypeId=%d, err=%v", i+1, msg.TicketTypeID, err)
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}

	// 三次回滚全失败 → 删掉标记（允许下次重试再尝试）+ 写补偿记录
	svcCtx.Redis.Del(ctx, releaseKey)

	reason := "建订单3次重试失败进入死信，回滚库存3次INCRBY也失败"
	if err := svcCtx.Redis.RecordCompensation(ctx, uint64(msg.TicketTypeID), msg.Quantity, reason, msg.OrderNo); err != nil {
		logx.Errorf("[DEAD] 写入补偿记录失败: ticketTypeId=%d, err=%v", msg.TicketTypeID, err)
		return nil
	}
	logx.Errorf("[DEAD][COMPENSATE] 补偿记录已写入: ticketTypeId=%d, quantity=%d", msg.TicketTypeID, msg.Quantity)

	return nil
}

// handleTimeOutCancel 超时取消订单
func handleTimeOutCancel(ctx context.Context, svcCtx *svc.ServiceContext, body []byte) error {
	var msg mq.TimeOutCancelMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logx.Errorf("[DEAD] 超时消息反序列化失败: %v, body=%s", err, string(body))
		return nil // 格式不对
	}

	var order db.Order
	if err := svcCtx.DB.Where("order_no = ?", msg.OrderNo).First(&order).Error; err != nil {
		logx.Errorf("[DEAD] 订单不存在: %s", msg.OrderNo)
		return nil // 订单已经没了
	}

	if order.Status != "unpaid" {
		return nil
	}

	// 查订单状态补偿记录
	compKey := "compensate:order:" + msg.OrderNo
	if exists, _ := svcCtx.Redis.Exists(ctx, compKey); exists {
		logx.Infof("[DEAD] 订单已支付（补偿记录存在），跳过回滚: orderNo=%s", msg.OrderNo)
		svcCtx.DB.Model(&order).Update("status", db.OrderPaid)
		svcCtx.Redis.RemoveOrderStatusCompensation(ctx, msg.OrderNo)
		return nil
	}

	// 防重复回滚：SETNX 标记，只允许执行一次
	releaseKey := "release:timeout:" + msg.OrderNo
	ok, err := svcCtx.Redis.Client().SetNX(ctx, releaseKey, "1", 30*time.Minute).Result()
	if err != nil {
		logx.Errorf("[DEAD] 释放标记设置失败: orderNo=%s, err=%v", msg.OrderNo, err)
		return err // NACK 重试
	}
	if !ok {
		// 已释放过
		logx.Infof("[DEAD] 库存已释放过，跳过: orderNo=%s", msg.OrderNo)
		svcCtx.DB.Model(&order).Update("status", db.OrderCanceled)
		return nil
	}

	// 回滚库存
	stockKey := fmt.Sprintf("stock:ticket:%d", msg.TicketTypeID)
	for i := 0; i < 3; i++ {
		_, err := svcCtx.Redis.IncrBy(ctx, stockKey, int64(msg.Quantity))
		if err == nil {
			break
		}
		logx.Errorf("[DEAD] 超时取消回滚库存失败(第%d次): orderNo=%s, err=%v", i+1, msg.OrderNo, err)
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}

	// 标记订单已取消
	svcCtx.DB.Model(&order).Update("status", db.OrderCanceled)

	logx.Infof("[DEAD] 订单超时未支付已取消: %s", msg.OrderNo)
	return nil
}
