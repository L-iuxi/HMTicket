
package logic

import (
	"context"

	"Ticket/app/inventory/inventory-rpc/internal/svc"
	"Ticket/app/inventory/inventory-rpc/inventory"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStockLogic {
	return &UpdateStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改库存
func (l *UpdateStockLogic) UpdateStock(in *inventory.UpdateStockReq) (*inventory.CommonResp, error) {

	if err := l.svcCtx.Redis.Set(
		l.ctx,
		redis.StockKey(in.TicketTypeId),
		in.Stock,
		0,
	); err != nil {
		return nil, err
	}

	return &inventory.CommonResp{
		Success: true,
		Message: "修改库存成功",
	}, nil
}
