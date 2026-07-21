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

type CreateShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShowLogic {
	return &CreateShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShowLogic) CreateShow(req *types.CreateShowReq) (resp *types.CreateShowResp, err error) {
	_, err = l.svcCtx.EventRpc.CreateShow(l.ctx, &event.CreateShowReq{
		EventId:   uint64(req.EventID),
		Name:      req.Name,
		ShowTime:  req.ShowTime,
		EndTime:   req.EndTime,
		SortOrder: int32(req.SortOrder),
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateShowResp{}, nil
}
