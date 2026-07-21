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

type UpdateShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShowLogic {
	return &UpdateShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShowLogic) UpdateShow(req *types.UpdateShowReq) (resp *types.UpdateShowResp, err error) {
	_, err = l.svcCtx.EventRpc.UpdateShow(l.ctx, &event.UpdateShowReq{
		ShowId:    uint64(req.ShowID),
		Name:      req.Name,
		ShowTime:  req.ShowTime,
		EndTime:   req.EndTime,
		SortOrder: int32(req.SortOrder),
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateShowResp{}, nil
}
