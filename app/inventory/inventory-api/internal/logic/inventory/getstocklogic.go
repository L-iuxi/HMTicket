// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package inventory

import (
	"context"

	"Ticket/app/inventory/inventory-api/internal/svc"
	"Ticket/app/inventory/inventory-api/internal/types"
	"Ticket/app/inventory/inventory-rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStockLogic {
	return &GetStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStockLogic) GetStock(req *types.GetStockReq) (*types.GetStockResp, error) {

	resp, err := l.svcCtx.InventoryRpc.GetStock(l.ctx, &inventory.GetStockReq{
		TicketTypeId: req.TicketTypeID,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetStockResp{
		Success:      resp.Success,
		Message:      resp.Message,
		TicketTypeID: resp.TicketTypeId,
		Stock:        resp.Stock,
	}, nil
}
