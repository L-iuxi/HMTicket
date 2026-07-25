package logic

import (
	"context"

	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/app/order/order-rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户订单列表
func (l *GetOrderListLogic) GetOrderList(in *order.GetOrderListReq) (*order.GetOrderListResp, error) {

	var orders []db.Order

	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).Order("created_at desc").Find(&orders).Error

	if err != nil {
		return nil, err
	}

	resp := &order.GetOrderListResp{}

	for _, orde := range orders {

		resp.Orders = append(resp.Orders, &order.OrderInfo{
			UserId:       uint32(orde.UserID),
			EventId:      orde.EventID,
			ShowId:       uint32(orde.ShowID),
			TicketTypeId: uint32(orde.TicketTypeID),
			OrderNo:      orde.OrderNo,
			Quantity:     int32(orde.Quantity),
			TotalPrice:   orde.TotalPrice,
			Status:       orde.Status,
		})
	}

	return resp, nil
}
