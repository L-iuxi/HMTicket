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

type GetShowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShowListLogic {
	return &GetShowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShowListLogic) GetShowList(req *types.GetShowListReq) (*types.GetShowListResp, error) {
	list, err := l.svcCtx.EventRpc.GetShowList(l.ctx, &event.GetShowListReq{
		EventId: uint64(req.EventId),
	})
	if err != nil {
		return nil, err
	}

	resp := &types.GetShowListResp{}
	for _, s := range list.Shows {
		resp.Shows = append(resp.Shows, types.ShowListItem{
			ShowID:    int64(s.Id),
			EventID:   int64(s.EventId),
			Name:      s.Name,
			ShowTime:  s.ShowTime,
			EndTime:   s.EndTime,
			Status:    s.Status,
			SoldCount: s.SoldCount,
			SortOrder: s.SortOrder,
			Venue:     s.Venue,
		})
	}

	return resp, nil
}
