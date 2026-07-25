package logic

import (
	"context"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	db "Ticket/app/event/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventListLogic {
	return &GetEventListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEventListLogic) GetEventList(in *event.GetEventListReq) (*event.GetEventListResp, error) {

	var events []db.Event

	if err := l.svcCtx.DB.Order("start_time asc").Find(&events).Error; err != nil {
		return nil, err
	}

	resp := &event.GetEventListResp{}

	for _, e := range events {

		resp.Events = append(resp.Events, &event.EventInfo{
			Id:          uint64(e.ID),
			Title:       e.Title,
			Description: e.Description,
			Location:    e.Location,
			CoverImage:  e.CoverImage,
			StartTime:   e.StartTime.Unix(),
			EndTime:     e.EndTime.Unix(),
			Status:      e.Status,
			TotalStock:  int32(e.TotalStock),
		})
	}

	return resp, nil
}
