// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"context"

	"Ticket/app/admin/internal/svc"
	"Ticket/app/admin/internal/types"
	"Ticket/app/event/event-rpc/event"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEventLogic {
	return &UpdateEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEventLogic) UpdateEvent(req *types.UpdateEventReq) (resp *types.UpdateEventResp, err error) {
	_, err = l.svcCtx.EventRpc.UpdateEvent(l.ctx, &event.UpdateEventReq{
		EventId:     uint64(req.EventID),
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		CoverImage:  req.CoverImage,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateEventResp{}, nil
}
