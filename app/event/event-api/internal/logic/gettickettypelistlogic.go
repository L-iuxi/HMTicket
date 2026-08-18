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

type GetTicketTypeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTicketTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTicketTypeListLogic {
	return &GetTicketTypeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTicketTypeListLogic) GetTicketTypeList(req *types.GetTicketTypeListReq) (*types.GetTicketTypeListResp, error) {
	list, err := l.svcCtx.EventRpc.GetTicketTypeList(l.ctx, &event.GetTicketTypeListReq{
		ShowId: uint64(req.ShowId),
	})
	if err != nil {
		return nil, err
	}

	resp := &types.GetTicketTypeListResp{}
	for _, t := range list.TicketTypes {
		resp.TicketTypes = append(resp.TicketTypes, types.TicketTypeListItem{
			TicketTypeID: int64(t.Id),
			EventID:      int64(t.EventId),
			ShowID:       int64(t.ShowId),
			Name:         t.Name,
			Price:        t.Price,
			Stock:        t.Stock,
			MaxPerUser:   t.MaxPerUser,
			SortOrder:    t.SortOrder,
		})
	}

	return resp, nil
}
