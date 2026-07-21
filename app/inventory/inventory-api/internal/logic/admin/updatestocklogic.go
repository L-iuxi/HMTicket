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

type UpdateStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStockLogic {
	return &UpdateStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStockLogic) UpdateStock(req *types.UpdateStockReq) (*types.UpdateStockResp, error) {

	resp, err := l.svcCtx.InventoryRpc.UpdateStock(l.ctx, &inventory.UpdateStockReq{
		TicketTypeId: req.TicketTypeID,
		Stock:        req.Stock,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateStockResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}
