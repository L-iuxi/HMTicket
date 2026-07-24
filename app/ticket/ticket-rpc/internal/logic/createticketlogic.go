package logic

import (
	"context"

	db "Ticket/app/ticket/model"
	"Ticket/app/ticket/ticket-rpc/internal/svc"
	"Ticket/app/ticket/ticket-rpc/ticket"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTicketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTicketLogic {
	return &CreateTicketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 支付成功后出票

func (l *CreateTicketLogic) CreateTicket(in *ticket.CreateTicketReq) (*ticket.CommonResp, error) {

	// 参数校验
	if in.UserId == 0 || in.EventId == 0 || in.TicketTypeId == 0 {
		return &ticket.CommonResp{
			Success: false,
			Message: "参数错误",
		}, nil
	}

	t := db.Ticket{
		UserID:         uint(in.UserId),
		EventID:        uint(in.EventId),
		ShowID:         uint(in.ShowId),
		TicketTypeID:   uint(in.TicketTypeId),
		OrderNo:        in.OrderNo,
		Quantity:       int(in.Quantity),
		TotalPrice:     in.TotalPrice,
		Status:         "unused",
		QRCode:         uuid.NewString(),
		DiscountCode:   in.DiscountCode,
		RealName:       in.RealName,
		IDCard:         in.IdCard,
		Phone:          in.Phone,
		TransferStatus: "none",
	}

	if err := l.svcCtx.DB.Create(&t).Error; err != nil {
		return nil, err
	}

	return &ticket.CommonResp{
		Success: true,
		Message: "出票成功",
	}, nil
}
