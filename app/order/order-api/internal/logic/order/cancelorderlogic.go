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

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderReq) (*types.CancelOrderResp, error) {

	resp, err := l.svcCtx.OrderRpc.CancelOrder(l.ctx, &order.CancelOrderReq{
		OrderNo: req.OrderNo,
	})
	if err != nil {
		return nil, err
	}

	return &types.CancelOrderResp{
		Success: resp.Success,
		Message: "取消成功",
	}, nil
}
