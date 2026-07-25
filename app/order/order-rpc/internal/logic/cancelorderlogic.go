package logic

import (
	"context"

	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	db "Ticket/app/order/model"
	"Ticket/app/order/order-rpc/internal/svc"
	"Ticket/app/order/order-rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消订单
func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderReq) (*order.CancelOrderResp, error) {

	var orde db.Order

	if err := l.svcCtx.DB.
		Where("order_no = ?", in.OrderNo).
		First(&orde).Error; err != nil {
		return nil, err
	}

	// 只有未支付订单取消时才回滚库存
	if orde.Status == db.OrderUnpaid {
		_, err := l.svcCtx.InventoryRpc.ReleaseStock(l.ctx, &inventoryclient.ReleaseStockReq{
			TicketTypeId: uint64(orde.TicketTypeID),
			Quantity:     int32(orde.Quantity),
		})
		if err != nil {
			logx.Errorf("取消订单回滚库存失败: orderNo=%s, err=%v", in.OrderNo, err)
		}
	}

	orde.Status = db.OrderCanceled

	if err := l.svcCtx.DB.Save(&orde).Error; err != nil {
		return nil, err
	}

	return &order.CancelOrderResp{
		Success: true,
	}, nil
}
