package logic

import (
	"context"
	"errors"

	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/app/order/order-rpc/order"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取订单详情
func (l *GetOrderLogic) GetOrder(in *order.GetOrderReq) (*order.GetOrderResp, error) {

	var orde db.Order

	err := l.svcCtx.DB.Where("order_no = ?", in.OrderNo).First(&orde).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ORDER_NOT_FOUND)
		}
		return nil, err
	}

	return &order.GetOrderResp{
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
