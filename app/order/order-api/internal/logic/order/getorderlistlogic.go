// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package order

import (
	"context"

	"Ticket/app/order/order-api/internal/common"
	"Ticket/app/order/order-api/internal/svc"
	"Ticket/app/order/order-api/internal/types"
	"Ticket/app/order/order-rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderListLogic) GetOrderList(req *types.GetOrderListReq) (*types.GetOrderListResp, error) {

	userId := common.GetUserID(l.ctx)

	resp, err := l.svcCtx.OrderRpc.GetOrderList(l.ctx, &order.GetOrderListReq{
		UserId: uint32(userId),
	})
	if err != nil {
		return nil, err
	}

	orders := make([]types.GetOrderInfo, 0, len(resp.Orders))

	for _, o := range resp.Orders {
		orders = append(orders, types.GetOrderInfo{
			UserID:       uint64(o.UserId),
			EventID:      o.EventId,
			ShowID:       uint64(o.ShowId),
			TicketTypeID: uint64(o.TicketTypeId),
			OrderNo:      o.OrderNo,
			Quantity:     o.Quantity,
			TotalPrice:   o.TotalPrice,
			Status:       o.Status,
		})
	}

	return &types.GetOrderListResp{
		Orders: orders,
	}, nil
}
