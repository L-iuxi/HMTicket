// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"context"

	"Ticket/app/admin/internal/svc"
	"Ticket/app/admin/internal/types"
	"Ticket/app/event/event-rpc/event"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTicketTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTicketTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTicketTypeLogic {
	return &UpdateTicketTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTicketTypeLogic) UpdateTicketType(req *types.UpdateTicketTypeReq) (resp *types.UpdateTicketTypeResp, err error) {
	_, err = l.svcCtx.EventRpc.UpdateTicketType(l.ctx, &event.UpdateTicketTypeReq{
		TicketTypeId: uint64(req.TicketTypeID),
		Name:         req.Name,
		Price:        req.Price,
		Stock:        req.Stock,
		MaxPerUser:   req.MaxPerUser,
		SortOrder:    int32(req.SortOrder),
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateTicketTypeResp{}, nil
}
