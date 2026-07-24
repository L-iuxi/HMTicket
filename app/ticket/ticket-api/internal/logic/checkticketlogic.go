// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"Ticket/app/ticket/ticket-api/internal/svc"
	"Ticket/app/ticket/ticket-api/internal/types"
	"Ticket/app/ticket/ticket-rpc/ticket"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTicketLogic {
	return &CheckTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckTicketLogic) CheckTicket(req *types.CheckTicketReq) (*types.CommonResp, error) {

	resp, err := l.svcCtx.TicketRpc.VerifyTicket(l.ctx, &ticket.VerifyTicketReq{
		TicketId: req.TicketId,
	})
	if err != nil {
		return nil, err
	}

	code := int64(0)
	if !resp.Success {
		code = 1
	}

	return &types.CommonResp{
		Code: code,
		Msg:  resp.Message,
	}, nil
}
