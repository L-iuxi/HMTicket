package logic

import (
	"context"

	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/app/order/order-rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改订单状态
func (l *UpdateOrderStatusLogic) UpdateOrderStatus(in *order.UpdateStatusReq) (*order.UpdateStatusResp, error) {

	var orde db.Order

	if err := l.svcCtx.DB.
		Where("order_no = ?", in.OrderNo).
		First(&orde).Error; err != nil {
		return nil, err
	}

	orde.Status = in.Status

	if err := l.svcCtx.DB.Save(&orde).Error; err != nil {
		return nil, err
	}

	return &order.UpdateStatusResp{
		Success: true,
	}, nil
}
