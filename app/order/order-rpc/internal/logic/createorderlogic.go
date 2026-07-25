package logic

import (
	"context"

	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/app/order/order-rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateOrder 同步创建订单（供 RPC 调用，MQ 异步模式下不再是主路径）
// 分布式锁已上提到 API 层
func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
	//orderNo := uuid.NewString()

	orde := db.Order{
		UserID:       uint(in.UserId),
		EventID:      in.EventId,
		ShowID:       uint(in.ShowId),
		TicketTypeID: uint(in.TicketTypeId),
		OrderNo:      in.OrderNo,
		Quantity:     int(in.Quantity),
		TotalPrice:   in.TotalPrice,
		Status:       "unpaid",
	}

	if err := l.svcCtx.DB.Create(&orde).Error; err != nil {
		return nil, err
	}

	return &order.CreateOrderResp{
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
