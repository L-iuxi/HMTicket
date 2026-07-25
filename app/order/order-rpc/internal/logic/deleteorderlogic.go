package logic

import (
	"context"

	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/app/order/order-rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrderLogic {
	return &DeleteOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除订单
func (l *DeleteOrderLogic) DeleteOrder(in *order.DeleteOrderReq) (*order.DeleteOrderResp, error) {

	var orde db.Order

	if err := l.svcCtx.DB.
		Where("order_no = ?", in.OrderNo).
		First(&orde).Error; err != nil {
		return nil, err
	}

	if err := l.svcCtx.DB.Delete(&orde).Error; err != nil {
		return nil, err
	}

	return &order.DeleteOrderResp{
		Success: true,
	}, nil
}
