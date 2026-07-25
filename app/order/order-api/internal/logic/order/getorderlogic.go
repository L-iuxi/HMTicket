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
	if err != nil {
		return nil, err
	}

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
