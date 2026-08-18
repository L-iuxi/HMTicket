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

type GetEventListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventListLogic {
	return &GetEventListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventListLogic) GetEventList(req *types.GetEventListReq) (*types.GetEventListResp, error) {
	list, err := l.svcCtx.EventRpc.GetEventList(l.ctx, &event.GetEventListReq{})
	if err != nil {
		return nil, err
	}

	resp := &types.GetEventListResp{}
	for _, e := range list.Events {
		resp.Events = append(resp.Events, types.EventListItem{
			EventID:     int64(e.Id),
			Title:       e.Title,
			Description: e.Description,
			Location:    e.Location,
			CoverImage:  e.CoverImage,
			StartTime:   e.StartTime,
			EndTime:     e.EndTime,
			Status:      e.Status,
			TotalStock:  e.TotalStock,
		})
	}

	return resp, nil
}
