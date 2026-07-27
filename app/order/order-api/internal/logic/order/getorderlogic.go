package order

import (
	"context"
	"strings"

	"Ticket/app/order/order-api/internal/svc"
	"Ticket/app/order/order-api/internal/types"
	"Ticket/app/order/order-rpc/order"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderReq) (*types.GetOrderResp, error) {
	resp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderReq{
		OrderNo: req.OrderNo,
	})
	if err == nil {
		return &types.GetOrderResp{
			UserID:       uint64(resp.Order.UserId),
			EventID:      resp.Order.EventId,
			ShowID:       uint64(resp.Order.ShowId),
			TicketTypeID: uint64(resp.Order.TicketTypeId),
			OrderNo:      resp.Order.OrderNo,
			Quantity:     resp.Order.Quantity,
			TotalPrice:   resp.Order.TotalPrice,
			Status:       resp.Order.Status,
		}, nil
	}

	// MySQL 没查到，可能是 MQ 还没消费完。查 Redis 回退。
	if st, ok := status.FromError(err); ok {
		// gRPC error 格式 "ErrCode:600001, ErrMsg:订单不存在"
		msg := st.Message()
		if strings.Contains(msg, "600001") || strings.Contains(msg, "订单不存在") {
			return l.checkCreatingStatus(req.OrderNo)
		}
	}

	return nil, err
}

// checkCreatingStatus 通过 Redis 查询异步建订单的实时状态
func (l *GetOrderLogic) checkCreatingStatus(orderNo string) (*types.GetOrderResp, error) {
	creatingKey := "order:creating:" + orderNo

	// 查 order:creating:{orderNo} → 拿到 idemKey
	idemKey, err := l.svcCtx.Redis.Get(l.ctx, creatingKey)
	if err != nil || idemKey == "" {
		// 映射不存在，订单号无效或已过期
		return nil, xerr.NewErrCode(xerr.ORDER_NOT_FOUND)
	}

	// 查 idemKey 状态
	val, err := l.svcCtx.Idempotent.Get(l.ctx, idemKey)
	if err != nil || val == "" {
		// idemKey 不存在，但映射还在（异常），返回 404
		return nil, xerr.NewErrCode(xerr.ORDER_NOT_FOUND)
	}

	switch {
	case val == "pending":
		return &types.GetOrderResp{
			OrderNo: orderNo,
			Status:  "creating",
		}, nil

	case strings.HasPrefix(val, "ok:"):
		// 订单应该已创建，但 MySQL 没查到。重试一次 RPC
		resp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderReq{OrderNo: orderNo})
		if err == nil {
			return &types.GetOrderResp{
				UserID:       uint64(resp.Order.UserId),
				EventID:      resp.Order.EventId,
				ShowID:       uint64(resp.Order.ShowId),
				TicketTypeID: uint64(resp.Order.TicketTypeId),
				OrderNo:      resp.Order.OrderNo,
				Quantity:     resp.Order.Quantity,
				TotalPrice:   resp.Order.TotalPrice,
				Status:       resp.Order.Status,
			}, nil
		}
		// 还是查不到，返回 creating（可能主从延迟）
		return &types.GetOrderResp{
			OrderNo: orderNo,
			Status:  "creating",
		}, nil

	case val == "failed":
		return &types.GetOrderResp{
			OrderNo: orderNo,
			Status:  "failed",
		}, nil

	default:
		return nil, xerr.NewErrCode(xerr.ORDER_NOT_FOUND)
	}
}
