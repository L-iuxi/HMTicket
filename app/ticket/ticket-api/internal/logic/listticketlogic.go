// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"Ticket/app/ticket/ticket-api/internal/common"
	"Ticket/app/ticket/ticket-api/internal/svc"
	"Ticket/app/ticket/ticket-api/internal/types"

	"Ticket/app/ticket/ticket-rpc/ticket"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTicketLogic {
	return &ListTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTicketLogic) ListTicket(req *types.ListTicketReq) (*types.ListTicketResp, error) {

	userId := common.GetUserID(l.ctx)

	resp, err := l.svcCtx.TicketRpc.ListUserTickets(l.ctx, &ticket.ListUserTicketsReq{
		UserId:   uint64(userId),
		Status:   req.Status,
		Page:     uint32(req.Page),
		PageSize: uint32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.TicketInfo, 0, len(resp.Tickets))

	for _, t := range resp.Tickets {
		list = append(list, types.TicketInfo{
			TicketId:       t.Id,
			OrderNo:        t.OrderNo,
			EventId:        t.EventId,
			ShowId:         t.ShowId,
			TicketTypeId:   t.TicketTypeId,
			Quantity:       int(t.Quantity),
			TotalPrice:     t.TotalPrice,
			Status:         t.Status,
			RealName:       t.RealName,
			Phone:          t.Phone,
			TransferStatus: t.TransferStatus,
			CreatedAt:      t.CreatedAt,
		})
	}

	return &types.ListTicketResp{
		Code:  0,
		Msg:   "success",
		Total: resp.Total,
		List:  list,
	}, nil
}
