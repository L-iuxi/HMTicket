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

type CreateTicketTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTicketTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTicketTypeLogic {
	return &CreateTicketTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTicketTypeLogic) CreateTicketType(req *types.CreateTicketTypeReq) (resp *types.CreateTicketTypeResp, err error) {
	_, err = l.svcCtx.EventRpc.CreateTicketType(l.ctx, &event.CreateTicketTypeReq{
		EventId:    uint64(req.EventID),
		ShowId:     uint64(req.ShowID),
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		MaxPerUser: req.MaxPerUser,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateTicketTypeResp{}, nil
}
