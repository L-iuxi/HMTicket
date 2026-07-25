// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"Ticket/app/event/event-api/internal/svc"
	"Ticket/app/event/event-api/internal/types"
	"Ticket/app/event/event-rpc/event"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTicketTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTicketTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTicketTypeLogic {
	return &GetTicketTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTicketTypeLogic) GetTicketType(req *types.GetTicketTypeReq) (*types.GetTicketTypeResp, error) {

	t, err := l.svcCtx.EventRpc.GetTicketType(l.ctx, &event.GetTicketTypeReq{
		TicketTypeId: uint64(req.TicketTypeId),
	})
	if err != nil {
		return nil, err
	}

	return &types.GetTicketTypeResp{
		TicketTypeID: int64(t.TicketType.Id),
		ShowID:       int64(t.TicketType.ShowId),
		Name:         t.TicketType.Name,
		Price:        t.TicketType.Price,
		Stock:        t.TicketType.Stock,
	}, nil
}
