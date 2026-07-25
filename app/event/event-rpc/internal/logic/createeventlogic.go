package logic

import (
	"context"
	"strconv"
	"time"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	db "Ticket/app/event/model"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEventLogic {
	return &CreateEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建活动
func (l *CreateEventLogic) CreateEvent(in *event.CreateEventReq) (*event.CommonResp, error) {
	startTime, err := time.Parse("2006-01-02 15:04:05", in.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", in.EndTime)
	if err != nil {
		return nil, err
	}

	if in.Title == "" {
		return nil, xerr.NewErrCode(xerr.EVENT_NAME_EMPTY)
	}
	if in.Location == "" {
		return nil, xerr.NewErrCode(xerr.EVENT_LOCATION_EMPTY)
	}
	if in.StartTime == "" {
		return nil, xerr.NewErrCode(xerr.EVENT_START_TIME_EMPTY)
	}
	if in.EndTime == "" {
		return nil, xerr.NewErrCode(xerr.EVENT_END_TIME_EMPTY)
	}
	if !endTime.After(startTime) {
		return nil, xerr.NewErrCode(xerr.EVENT_TIME_INVALID)
	}

	// 检查活动是否重复（同名+同时间）
	var count int64
	l.svcCtx.DB.Model(&db.Event{}).
		Where("title = ? AND start_time = ?", in.Title, startTime).
		Count(&count)
	if count > 0 {
		return nil, xerr.NewErrCode(xerr.EVENT_ALREADY_EXISTS)
	}

	e := db.Event{
		Title:       in.Title,
		Description: in.Description,
		Location:    in.Location,
		CoverImage:  in.CoverImage,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "draft",
		TotalStock:  0,
	}

	if err := l.svcCtx.DB.Create(&e).Error; err != nil {
		return nil, err
	}

	return &event.CommonResp{
		Success: true,
		Message: strconv.FormatUint(uint64(e.ID), 10),
	}, nil
}
