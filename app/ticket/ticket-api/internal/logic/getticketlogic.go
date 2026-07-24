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

type GetTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTicketLogic {
	return &GetTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTicketLogic) GetTicket(req *types.GetTicketReq) (*types.GetTicketResp, error) {

	resp, err := l.svcCtx.TicketRpc.GetTicket(l.ctx, &ticket.GetTicketReq{
		TicketId: req.TicketId,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetTicketResp{
		Code:           0,
		Msg:            "success",
		TicketId:       resp.Id,
		OrderNo:        resp.OrderNo,
		EventId:        resp.EventId,
		ShowId:         resp.ShowId,
		TicketTypeId:   resp.TicketTypeId,
		Quantity:       int(resp.Quantity),
		TotalPrice:     resp.TotalPrice,
		Status:         resp.Status,
		QrCode:         resp.QrCode,
		DiscountCode:   resp.DiscountCode,
		RealName:       resp.RealName,
		IdCard:         resp.IdCard,
		Phone:          resp.Phone,
		TransferredTo:  resp.TransferredTo,
		TransferStatus: resp.TransferStatus,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      resp.UpdatedAt,
	}, nil
}
