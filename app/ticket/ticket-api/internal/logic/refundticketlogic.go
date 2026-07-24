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

type RefundTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundTicketLogic {
	return &RefundTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundTicketLogic) RefundTicket(req *types.RefundTicketReq) (*types.CommonResp, error) {

	resp, err := l.svcCtx.TicketRpc.RefundTicket(l.ctx, &ticket.RefundTicketReq{
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
