// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	"Ticket/app/order/order-api/internal/common"
	"Ticket/app/order/order-api/internal/svc"
	"Ticket/app/order/order-api/internal/types"
	"Ticket/app/order/order-rpc/order"
	"Ticket/app/order/order-rpc/orderclient"
	"Ticket/app/payment/pay-rpc/paymentclient"
	"Ticket/app/ticket/ticket-rpc/ticketclient"
	"Ticket/common/xerr"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayorderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayorderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayorderLogic {
	return &PayorderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayorderLogic) Payorder(req *types.PayOrderReq) (resp *types.PayOrderResp, err error) {
	userID := common.GetUserID(l.ctx)

	if req.RequestId == "" {
		return &types.PayOrderResp{Success: false, Message: "缺少请求ID"}, nil
	}

	// 幂等检查在前——重复请求直接返回，不消耗限流配额
	idemKey := redis.IdempotentKey("pay", userID, req.OrderNo, req.RequestId)

	val, err := l.svcCtx.Idempotent.Get(l.ctx, idemKey)
	if err != nil {
		return nil, fmt.Errorf("幂等检查失败: %w", err)
	}
	if strings.HasPrefix(val, "ok:") {
		return &types.PayOrderResp{Success: true, Message: "支付成功！已出票^^"}, nil
	}
	if val == "pending" {
		return &types.PayOrderResp{Success: false, Message: "支付处理中，请稍后"}, nil
	}

	if err := l.svcCtx.Idempotent.Start(l.ctx, idemKey, 5*time.Minute); err != nil {
		if err == redis.ErrDuplicateRequest {
			return &types.PayOrderResp{Success: false, Message: "请勿重复提交"}, nil
		}
		return nil, fmt.Errorf("幂等标记失败: %w", err)
	}

	// 用户级限流
	limitkey := fmt.Sprintf("limit:pay:user:%d", userID)
	ok, err := l.svcCtx.RateLimiter.Allow(l.ctx, limitkey, 5, time.Second)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, xerr.NewErrMsg("请求过于频繁请稍候重试")
	}

	// 订单级限流
	limitkey = fmt.Sprintf("limit:pay:order:%s", req.OrderNo)
	ok, err = l.svcCtx.RateLimiter.Allow(l.ctx, limitkey, 1, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, xerr.NewErrMsg("请求过于频繁请稍候重试")
	}

	orderResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderReq{OrderNo: req.OrderNo})
	if err != nil || orderResp == nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return &types.PayOrderResp{Success: false, Message: "订单不存在"}, nil
	}

	if orderResp.Order.UserId != uint32(userID) {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return &types.PayOrderResp{Success: false, Message: "只能支付自己的订单"}, nil
	}

	if orderResp.Order.Status != "unpaid" {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return &types.PayOrderResp{Success: false, Message: "订单状态出错"}, nil
	}

	// 支付
	payResp, err := l.svcCtx.PaymentRpc.Pay(l.ctx, &paymentclient.PayReq{
		OrderNo:    req.OrderNo,
		Amount:     orderResp.Order.TotalPrice,
		TotalPrice: orderResp.Order.TotalPrice,
	})
	if err != nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		l.releaseStock(uint64(orderResp.Order.TicketTypeId), orderResp.Order.Quantity)
		l.updateOrderStatus(req.OrderNo, "failed")
		return nil, fmt.Errorf("支付失败: %w", err)
	}

	if !payResp.Success {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		l.releaseStock(uint64(orderResp.Order.TicketTypeId), orderResp.Order.Quantity)
		l.updateOrderStatus(req.OrderNo, "failed")
		return &types.PayOrderResp{Success: false, Message: "支付失败"}, nil
	}

	// 出票
	_, err = l.svcCtx.TicketRpc.CreateTicket(l.ctx, &ticketclient.CreateTicketReq{
		UserId:       uint64(userID),
		EventId:      orderResp.Order.EventId,
		ShowId:       uint64(orderResp.Order.ShowId),
		TicketTypeId: uint64(orderResp.Order.TicketTypeId),
		OrderNo:      req.OrderNo,
		Quantity:     orderResp.Order.Quantity,
		TotalPrice:   orderResp.Order.TotalPrice,
	})
	if err != nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		l.Errorf("CRITICAL: 出票失败，订单=%s, err=%v", req.OrderNo, err)
		l.updateOrderStatus(req.OrderNo, "ticket_failed")
		return nil, fmt.Errorf("出票失败，请联系客服，订单号: %s", req.OrderNo)
	}

	// 订单标记已支付
	l.updateOrderStatus(req.OrderNo, "paid")

	// 标记幂等成功
	l.svcCtx.Idempotent.Success(l.ctx, idemKey, "paid", 15*time.Minute)

	return &types.PayOrderResp{Success: true, Message: "支付成功！已出票^^"}, nil
}

// releaseStock 回滚库存（best-effort）
func (l *PayorderLogic) releaseStock(ticketTypeID uint64, quantity int32) {
	_, err := l.svcCtx.InventoryRpc.ReleaseStock(l.ctx, &inventoryclient.ReleaseStockReq{
		TicketTypeId: ticketTypeID,
		Quantity:     quantity,
	})
	if err != nil {
		l.Errorf("回滚库存失败: ticketTypeId=%d, quantity=%d, err=%v", ticketTypeID, quantity, err)
	}
}

// updateOrderStatus 修改订单状态，带重试
func (l *PayorderLogic) updateOrderStatus(orderNo, status string) {
	var lastErr error
	for range 3 {
		_, err := l.svcCtx.OrderRpc.UpdateOrderStatus(l.ctx, &orderclient.UpdateStatusReq{
			OrderNo: orderNo,
			Status:  status,
		})
		if err == nil {
			return
		}
		lastErr = err
	}
	l.Errorf("更新订单状态失败(已重试3次): orderNo=%s, status=%s, err=%v", orderNo, status, lastErr)
}
