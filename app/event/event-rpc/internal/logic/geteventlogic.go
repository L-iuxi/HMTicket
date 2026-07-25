package logic

import (
	"context"
	"errors"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	db "Ticket/app/event/model"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventLogic {
	return &GetEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取活动详情
func (l *GetEventLogic) GetEvent(in *event.GetEventReq) (*event.GetEventResp, error) {

	var e db.Event

	if err := l.svcCtx.DB.
		First(&e, in.EventId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.EVENT_NOT_FOUND)
		}
		return nil, err
	}

	return &event.GetEventResp{
		Event: &event.EventInfo{
			Id:          uint64(e.ID),
			Title:       e.Title,
			Description: e.Description,
			Location:    e.Location,
			CoverImage:  e.CoverImage,
			StartTime:   e.StartTime.Unix(),
			EndTime:     e.EndTime.Unix(),
			Status:      e.Status,
			TotalStock:  int32(e.TotalStock),
		},
	}, nil
}
