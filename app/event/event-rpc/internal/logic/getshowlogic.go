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

type GetShowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShowLogic {
	return &GetShowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取场次详情
func (l *GetShowLogic) GetShow(in *event.GetShowReq) (*event.GetShowResp, error) {

	var s db.Show

	if err := l.svcCtx.DB.
		First(&s, in.ShowId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.SHOW_NOT_FOUND)
		}
		return nil, err
	}

	return &event.GetShowResp{
		Show: &event.ShowInfo{
			Id:        uint64(s.ID),
			EventId:   uint64(s.EventID),
			Name:      s.Name,
			ShowTime:  s.ShowTime.Format("2006-01-02 15:04:05"),
			EndTime:   s.EndTime.Format("2006-01-02 15:04:05"),
			Status:    s.Status,
			SoldCount: int32(s.SoldCount),
			SortOrder: int32(s.SortOrder),
			Venue:     s.Venue,
		},
	}, nil
}
