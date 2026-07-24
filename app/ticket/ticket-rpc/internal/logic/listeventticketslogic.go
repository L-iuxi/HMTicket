package logic

import (
	"context"

	"Ticket/app/ticket/ticket-rpc/internal/svc"
	"Ticket/app/ticket/ticket-rpc/ticket"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEventTicketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListEventTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEventTicketsLogic {
	return &ListEventTicketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员查看某活动所有票
func (l *ListEventTicketsLogic) ListEventTickets(in *ticket.ListEventTicketsReq) (*ticket.ListEventTicketsResp, error) {
	// todo: add your logic here and delete this line

	return &ticket.ListEventTicketsResp{}, nil
}
