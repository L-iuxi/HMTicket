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

type GetTicketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTicketLogic {
	return &GetTicketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询票详情
func (l *GetTicketLogic) GetTicket(in *ticket.GetTicketReq) (*ticket.TicketInfo, error) {

	var t db.Ticket

	if err := l.svcCtx.DB.First(&t, in.TicketId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.TICKET_NOT_FOUND)
		}
		return nil, err
	}

	return &ticket.TicketInfo{
		Id:             uint64(t.ID),
		UserId:         uint64(t.UserID),
		EventId:        uint64(t.EventID),
		ShowId:         uint64(t.ShowID),
		TicketTypeId:   uint64(t.TicketTypeID),
		OrderNo:        t.OrderNo,
		Quantity:       int32(t.Quantity),
		TotalPrice:     t.TotalPrice,
		Status:         t.Status,
		QrCode:         t.QRCode,
		DiscountCode:   t.DiscountCode,
		RealName:       t.RealName,
		IdCard:         t.IDCard,
		Phone:          t.Phone,
		TransferredTo:  uint64(t.TransferredTo),
		TransferStatus: t.TransferStatus,
		CreatedAt:      t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
