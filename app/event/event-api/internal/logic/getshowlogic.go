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

type GetShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShowLogic {
	return &GetShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShowLogic) GetShow(req *types.GetShowReq) (*types.GetShowResp, error) {

	// 1. 查 Show
	showResp, err := l.svcCtx.EventRpc.GetShow(l.ctx, &event.GetShowReq{
		ShowId: uint64(req.ShowId),
	})
	if err != nil {
		return nil, err
	}

	// 2. 查 TicketTypeList
	tickets, err := l.svcCtx.EventRpc.GetTicketTypeList(l.ctx, &event.GetTicketTypeListReq{
		ShowId: uint64(req.ShowId),
	})
	if err != nil {
		return nil, err
	}

	resp := &types.GetShowResp{
		ShowID:   int64(showResp.Show.Id),
		Name:     showResp.Show.Name,
		ShowTime: showResp.Show.ShowTime,
		Venue:    showResp.Show.Venue,
	}

	// 3. 拼 TicketTypes
	for _, t := range tickets.TicketTypes {
		resp.TicketTypes = append(resp.TicketTypes, types.TicketTypeDetail{
			TicketTypeID: int64(t.Id),
			ShowID:       int64(t.ShowId),
			Name:         t.Name,
			Price:        t.Price,
			Stock:        t.Stock,
		})
	}

	return resp, nil
}
