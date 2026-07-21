// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package inventory

import (
	"context"

	"Ticket/app/inventory/inventory-api/internal/svc"
	"Ticket/app/inventory/inventory-api/internal/types"

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

func (l *UpdateStockLogic) UpdateStock(req *types.UpdateStockReq) (resp *types.UpdateStockResp, err error) {
	// todo: add your logic here and delete this line

	return
}
