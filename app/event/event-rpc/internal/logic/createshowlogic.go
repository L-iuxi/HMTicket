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

type CreateShowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShowLogic {
	return &CreateShowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建场次
func (l *CreateShowLogic) CreateShow(in *event.CreateShowReq) (*event.CommonResp, error) {
	var even db.Event
	//先检查活动是否存在
	if err := l.svcCtx.DB.First(&even, in.EventId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.EVENT_NOT_FOUND)
		}
		return nil, err
	}

	showTime, err := time.Parse("2006-01-02 15:04:05", in.ShowTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", in.EndTime)
	if err != nil {
		return nil, err
	}
	if in.Name == "" {
		return nil, xerr.NewErrCode(xerr.SHOW_NAME_EMPTY)
	}
	if !endTime.After(showTime) {
		return nil, xerr.NewErrCode(xerr.EVENT_TIME_INVALID)
	}

	if showTime.Before(even.StartTime) || endTime.After(even.EndTime) {
		return nil, xerr.NewErrCode(xerr.SHOW_TIME_INVALID)
	}

	show := &db.Show{
		EventID:   in.EventId,
		Name:      in.Name,
		ShowTime:  showTime,
		EndTime:   endTime,
		Status:    "draft",
		SoldCount: 0,
		SortOrder: in.SortOrder,
	}

	if err = l.svcCtx.DB.Create(show).Error; err != nil {
		return nil, err
	}

	return &event.CommonResp{}, nil
}
