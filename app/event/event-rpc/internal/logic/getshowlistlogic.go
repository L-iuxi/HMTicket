package logic

import (
	"context"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	db "Ticket/app/event/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetShowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShowListLogic {
	return &GetShowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取活动下所有场次
func (l *GetShowListLogic) GetShowList(in *event.GetShowListReq) (*event.GetShowListResp, error) {

	var shows []db.Show

	if err := l.svcCtx.DB.
		Where("event_id = ?", in.EventId).
		Order("show_time asc").
		Find(&shows).Error; err != nil {
		return nil, err
	}

	resp := &event.GetShowListResp{}

	for _, s := range shows {

		resp.Shows = append(resp.Shows, &event.ShowInfo{
			Id:        uint64(s.ID),
			EventId:   uint64(s.EventID),
			Name:      s.Name,
			ShowTime:  s.ShowTime.Format("2006-01-02 15:04:05"),
			EndTime:   s.EndTime.Format("2006-01-02 15:04:05"),
			Status:    s.Status,
			SoldCount: int32(s.SoldCount),
			SortOrder: int32(s.SortOrder),
			Venue:     s.Venue,
		})
	}

	return resp, nil
}
