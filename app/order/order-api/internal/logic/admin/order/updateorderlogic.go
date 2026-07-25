// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package order

import (
	"context"

	"Ticket/app/order/order-api/internal/svc"
	"Ticket/app/order/order-api/internal/types"
	"Ticket/app/order/order-rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderLogic) UpdateOrder(req *types.UpdateOrderReq) (*types.UpdateOrderResp, error) {

	resp, err := l.svcCtx.OrderRpc.UpdateOrder(l.ctx, &order.UpdateOrderReq{
		OrderNo:  req.OrderNo,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateOrderResp{
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
