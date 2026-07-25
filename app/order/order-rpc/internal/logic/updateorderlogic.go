package logic

import (
	"context"

	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/app/order/order-rpc/order"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改订单信息
func (l *UpdateOrderLogic) UpdateOrder(in *order.UpdateOrderReq) (*order.UpdateOrderResp, error) {

	var orde db.Order

	// 查询订单
	if err := l.svcCtx.DB.
		Where("order_no = ?", in.OrderNo).
		First(&orde).Error; err != nil {
		return nil, err
	}

	// 这里只允许未支付订单修改
	if orde.Status != "unpaid" {
		return nil, xerr.NewErrCode(xerr.ORDER_STATUS_INVALID)
	}

	// 更新数量
	orde.Quantity = int(in.Quantity)

	// TODO: 根据票价重新计算总价
	// order.TotalPrice = price * float64(order.Quantity)

	if err := l.svcCtx.DB.Save(&orde).Error; err != nil {
		return nil, err
	}

	return &order.UpdateOrderResp{
		Order: &order.OrderInfo{
			UserId:       uint32(orde.UserID),
			EventId:      orde.EventID,
			ShowId:       uint32(orde.ShowID),
			TicketTypeId: uint32(orde.TicketTypeID),
			OrderNo:      orde.OrderNo,
			Quantity:     int32(orde.Quantity),
			TotalPrice:   orde.TotalPrice,
			Status:       orde.Status,
		},
	}, nil
}
