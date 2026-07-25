package logic

import (
	"context"
	"errors"
	"time"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	db "Ticket/app/event/model"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEventLogic {
	return &UpdateEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改活动
func (l *UpdateEventLogic) UpdateEvent(in *event.UpdateEventReq) (*event.CommonResp, error) {
	var even db.Event
	if err := l.svcCtx.DB.First(&even, in.EventId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.EVENT_NOT_FOUND)
		}
		return nil, err
	}

	updates := make(map[string]interface{})

	if in.Title != "" {
		updates["title"] = in.Title
	}

	if in.Description != "" {
		updates["description"] = in.Description
	}

	if in.Location != "" {
		updates["location"] = in.Location
	}

	if in.CoverImage != "" {
		updates["cover_image"] = in.CoverImage
	}
	if even.Status != "draft" {
		return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
	}

	var startTime time.Time
	var endTime time.Time
	var parseErr error

	if in.StartTime != "" {
		startTime, parseErr = time.Parse("2006-01-02 15:04:05", in.StartTime)
		if parseErr != nil {
			return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
		}
		updates["start_time"] = startTime
	} else {
		startTime = even.StartTime
	}

	if in.EndTime != "" {
		endTime, parseErr = time.Parse("2006-01-02 15:04:05", in.EndTime)
		if parseErr != nil {
			return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
		}
		updates["end_time"] = endTime
	} else {
		endTime = even.EndTime
	}

	// 时间校验
	if in.StartTime != "" && in.EndTime != "" {
		if startTime.After(endTime) {
			return nil, xerr.NewErrCode(xerr.EVENT_TIME_INVALID)
		}
	} else if in.StartTime != "" {
		if startTime.After(even.EndTime) {
			return nil, xerr.NewErrCode(xerr.EVENT_TIME_INVALID)
		}
	} else if in.EndTime != "" {
		if even.StartTime.After(endTime) {
			return nil, xerr.NewErrCode(xerr.EVENT_TIME_INVALID)
		}
	}

	if len(updates) == 0 {
		return nil, xerr.NewErrCode(xerr.NO_DATA_TO_UPDATE)
	}

	if err := l.svcCtx.DB.Model(&even).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &event.CommonResp{}, nil
}
