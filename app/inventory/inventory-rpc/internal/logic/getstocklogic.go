
package logic

import (
	"context"
	"strconv"

	"Ticket/app/inventory/inventory-rpc/internal/svc"
	"Ticket/app/inventory/inventory-rpc/inventory"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStockLogic {
	return &GetStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询库存
func (l *GetStockLogic) GetStock(in *inventory.GetStockReq) (*inventory.GetStockResp, error) {

	key := redis.StockKey(in.TicketTypeId)

	value, err := l.svcCtx.Redis.Get(l.ctx, key)
	if err != nil {
		return &inventory.GetStockResp{
			Success:      false,
			Message:      "库存未初始化",
			TicketTypeId: in.TicketTypeId,
		}, nil
	}

	stock, err := strconv.Atoi(value)
	if err != nil {
		return &inventory.GetStockResp{
			Success:      false,
			Message:      "库存数据异常",
			TicketTypeId: in.TicketTypeId,
		}, nil
	}

	return &inventory.GetStockResp{
		Success:      true,
		Message:      "success",
		TicketTypeId: in.TicketTypeId,
		Stock:        int32(stock),
	}, nil
}
