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

type UpdateShowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShowLogic {
	return &UpdateShowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改场次
func (l *UpdateShowLogic) UpdateShow(in *event.UpdateShowReq) (*event.CommonResp, error) {
	// 查询场次
	var show db.Show
	if err := l.svcCtx.DB.First(&show, in.ShowId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.SHOW_NOT_FOUND)
		}
		return nil, err
	}

	// 查询所属活动
	var even db.Event
	if err := l.svcCtx.DB.First(&even, show.EventID).Error; err != nil {
		return nil, xerr.NewErrCode(xerr.EVENT_NOT_FOUND)
	}

	// 已开始售卖或已关闭的活动禁止修改场次
	if even.Status == db.EventStatusSelling || even.Status == db.EventStatusClosed {
		return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
	}

	updates := make(map[string]interface{})

	if in.Name != "" {
		updates["name"] = in.Name
	}

	if in.Venue != "" {
		updates["venue"] = in.Venue
	}

	// 时间处理
	showTime := show.ShowTime
	if in.ShowTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", in.ShowTime)
		if err != nil {
			return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
		}

		// 场次时间必须位于活动时间内
		if t.Before(even.StartTime) || t.After(even.EndTime) {
			return nil, xerr.NewErrCode(xerr.SHOW_TIME_INVALID)
		}

		showTime = t
		updates["show_time"] = showTime
	}

	if len(updates) == 0 {
		return nil, xerr.NewErrCode(xerr.NO_DATA_TO_UPDATE)
	}

	if err := l.svcCtx.DB.Model(&show).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &event.CommonResp{}, nil
}
