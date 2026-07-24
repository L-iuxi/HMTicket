package logic

import (
	"context"
	"errors"

	db "Ticket/app/ticket/model"
	"Ticket/app/ticket/ticket-rpc/internal/svc"
	"Ticket/app/ticket/ticket-rpc/ticket"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type VerifyTicketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTicketLogic {
	return &VerifyTicketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 验票（扫码）
func (l *VerifyTicketLogic) VerifyTicket(in *ticket.VerifyTicketReq) (*ticket.CommonResp, error) {
	var t db.Ticket
	if err := l.svcCtx.DB.First(&t, in.TicketId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &ticket.CommonResp{
				Success: false,
				Message: "票不存在",
			}, nil
		}
		return nil, err
	}

	if t.Status != "unused" {
		return &ticket.CommonResp{
			Success: false,
			Message: "票已被使用或无效",
		}, nil
	}

	t.Status = "used"
	if err := l.svcCtx.DB.Save(&t).Error; err != nil {
		return nil, err
	}

	return &ticket.CommonResp{
		Success: true,
		Message: "验票通过",
	}, nil
}
