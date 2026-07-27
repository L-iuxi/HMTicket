package logic

import (
	"context"

	"Ticket/app/inventory/inventory-rpc/internal/svc"
	"Ticket/app/inventory/inventory-rpc/inventory"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseStockLogic {
	return &ReleaseStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 回滚库存
func (l *ReleaseStockLogic) ReleaseStock(in *inventory.ReleaseStockReq) (*inventory.ReleaseStockResp, error) {

	remain, err := l.svcCtx.Redis.Client().IncrBy(
		l.ctx,
		redis.StockKey(in.TicketTypeId),
		int64(in.Quantity),
	).Result()

	if err != nil {
		return nil, err
	}

	return &inventory.ReleaseStockResp{
		Success:     true,
		Message:     "回滚库存成功",
		RemainStock: int32(remain),
	}, nil
}
