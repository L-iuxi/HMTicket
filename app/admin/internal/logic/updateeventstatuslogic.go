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

type UpdateEventStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEventStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEventStatusLogic {
	return &UpdateEventStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEventStatusLogic) UpdateEventStatus(req *types.UpdateEventStatusReq) (resp *types.UpdateEventStatusResp, err error) {
	_, err = l.svcCtx.EventRpc.UpdateEventStatus(l.ctx, &event.UpdateEventStatusReq{
		EventId: uint64(req.EventID),
		Status:  req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateEventStatusResp{}, nil
}
