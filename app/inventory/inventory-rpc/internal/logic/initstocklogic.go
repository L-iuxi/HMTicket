package logic

import (
	"context"

	"Ticket/app/inventory/inventory-rpc/internal/svc"
	"Ticket/app/inventory/inventory-rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitStockLogic {
	return &InitStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 初始化库存
func (l *InitStockLogic) InitStock(in *inventory.InitStockReq) (*inventory.CommonResp, error) {

	if err := l.svcCtx.Redis.Set(l.ctx, stockKey(in.TicketTypeId), in.Stock, 0); err != nil {
		return nil, err
	}

	return &inventory.CommonResp{
		Success: true,
		Message: "初始化库存成功",
	}, nil
}
