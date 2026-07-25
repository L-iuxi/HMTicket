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

type DeleteOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrderLogic {
	return &DeleteOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOrderLogic) DeleteOrder(req *types.DeleteOrderReq) (*types.DeleteOrderResp, error) {

	resp, err := l.svcCtx.OrderRpc.DeleteOrder(l.ctx, &order.DeleteOrderReq{
		OrderNo: req.OrderNo,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteOrderResp{
		Success: resp.Success,
		Message: "删除成功",
	}, nil
}
