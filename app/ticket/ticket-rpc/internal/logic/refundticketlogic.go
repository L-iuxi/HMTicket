package logic

import (
	"context"
	"errors"

	db "Ticket/app/ticket/model"
	"Ticket/app/ticket/ticket-rpc/internal/svc"
	"Ticket/app/ticket/ticket-rpc/ticket"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RefundTicketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundTicketLogic {
	return &RefundTicketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 退票
func (l *RefundTicketLogic) RefundTicket(in *ticket.RefundTicketReq) (*ticket.CommonResp, error) {

	var t db.Ticket

	if err := l.svcCtx.DB.First(&t, in.TicketId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.TICKET_NOT_FOUND)
		}
		return nil, err
	}

	// 已使用不能退
	if t.Status == "used" {
		return &ticket.CommonResp{
			Success: false,
			Message: "该票已使用，无法退票",
		}, nil
	}

	// 已退款
	if t.Status == "refunded" {
		return &ticket.CommonResp{
			Success: false,
			Message: "该票已退款",
		}, nil
	}

	t.Status = "refunded"

	if err := l.svcCtx.DB.Save(&t).Error; err != nil {
		return nil, err
	}

	return &ticket.CommonResp{
		Success: true,
		Message: "退票成功",
	}, nil
}
