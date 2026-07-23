package logic

import (
	"context"
	"strconv"

	"Ticket/app/inventory/inventory-rpc/internal/svc"
	"Ticket/app/inventory/inventory-rpc/inventory"

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

	key := stockKey(in.TicketTypeId)

	value, err := l.svcCtx.Redis.Get(l.ctx, key)

	if err == nil {
		stock, _ := strconv.Atoi(value)
		return &inventory.GetStockResp{
			Success:      true,
			Message:      "success",
			TicketTypeId: in.TicketTypeId,
			Stock:        int32(stock),
		}, nil
	}

	return &inventory.GetStockResp{
		Success:      false,
		Message:      "库存未初始化",
		TicketTypeId: in.TicketTypeId,
	}, nil
}
