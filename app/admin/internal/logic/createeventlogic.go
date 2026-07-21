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

type CreateEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEventLogic {
	return &CreateEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEventLogic) CreateEvent(req *types.CreateEventReq) (resp *types.CreateEventResp, err error) {
	_, err = l.svcCtx.EventRpc.CreateEvent(l.ctx, &event.CreateEventReq{
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

	return &types.CreateEventResp{}, nil
}
