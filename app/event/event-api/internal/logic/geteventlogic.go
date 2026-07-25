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

type GetEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventLogic {
	return &GetEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventLogic) GetEvent(req *types.GetEventReq) (*types.GetEventResp, error) {

	// 查 Event
	even, err := l.svcCtx.EventRpc.GetEvent(l.ctx, &event.GetEventReq{
		EventId: uint64(req.EventId),
	})
	if err != nil {
		return nil, err
	}

	//  查 ShowList
	shows, err := l.svcCtx.EventRpc.GetShowList(l.ctx, &event.GetShowListReq{
		EventId: uint64(req.EventId),
	})
	if err != nil {
		return nil, err
	}

	resp := &types.GetEventResp{
		EventID:    int64(even.Event.Id),
		Title:      even.Event.Title,
		Desc:       even.Event.Description,
		Status:     even.Event.Status,
		CoverImage: even.Event.CoverImage,
		StartTime:  even.Event.StartTime,
		EndTime:    even.Event.EndTime,
		Location:   even.Event.Location,
		TotalStock: even.Event.TotalStock,
	}

	//  拼 ShowDetail
	for _, s := range shows.Shows {
		resp.Shows = append(resp.Shows, types.ShowDetail{
			ShowID:   int64(s.Id),
			Name:     s.Name,
			ShowTime: s.ShowTime,
			Venue:    s.Venue,
		})
	}

	return resp, nil
}
