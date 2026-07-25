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

type UpdateEventStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateEventStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEventStatusLogic {
	return &UpdateEventStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改活动状态（draft、ready、selling、closed）

func (l *UpdateEventStatusLogic) UpdateEventStatus(in *event.UpdateEventStatusReq) (*event.CommonResp, error) {
	// 查询活动
	var even db.Event
	if err := l.svcCtx.DB.First(&even, in.EventId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.EVENT_NOT_FOUND)
		}
		return nil, err
	}

	// 状态是否合法
	validStatus := map[string]bool{
		"draft":   true,
		"ready":   true,
		"selling": true,
		"closed":  true,
	}

	if !validStatus[in.Status] {
		return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
	}

	// 校验状态流转
	switch even.Status {

	case db.EventStatusDraft:
		// draft -> ready
		if in.Status != db.EventStatusReady {
			return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
		}

	case db.EventStatusReady:
		// ready -> selling
		// 是否有场次
		var showCount int64
		err := l.svcCtx.DB.Model(&db.Show{}).Where("event_id = ?", even.ID).Count(&showCount).Error
		if err != nil {
			return nil, err
		}

		if showCount == 0 {
			return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
		}

		// 是否有场次
		err = l.svcCtx.DB.Model(&db.Show{}).Where("event_id = ?", even.ID).Count(&showCount).Error
		if err != nil {
			return nil, err
		}

		if showCount == 0 {
			return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
		}
		if in.Status != db.EventStatusSelling {
			return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
		}

	case db.EventStatusSelling:
		// selling -> closed
		if in.Status != db.EventStatusClosed {
			return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
		}

	case db.EventStatusClosed:
		return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)

	default:
		return nil, xerr.NewErrCode(xerr.EVENT_STATUS_INVALID)
	}

	// 更新状态
	if err := l.svcCtx.DB.Model(&even).Update("status", in.Status).Error; err != nil {
		return nil, err
	}

	return &event.CommonResp{}, nil
}
