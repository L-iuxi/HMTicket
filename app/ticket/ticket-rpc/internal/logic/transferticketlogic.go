package logic

import (
	"context"

	"Ticket/app/ticket/ticket-rpc/internal/svc"
	"Ticket/app/ticket/ticket-rpc/ticket"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferTicketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferTicketLogic {
	return &TransferTicketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 转赠
func (l *TransferTicketLogic) TransferTicket(in *ticket.TransferTicketReq) (*ticket.CommonResp, error) {
	// todo: add your logic here and delete this line

	return &ticket.CommonResp{}, nil
}
