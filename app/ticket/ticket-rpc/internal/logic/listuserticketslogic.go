package logic

import (
	"context"

	db "Ticket/app/ticket/model"
	"Ticket/app/ticket/ticket-rpc/internal/svc"
	"Ticket/app/ticket/ticket-rpc/ticket"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserTicketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserTicketsLogic {
	return &ListUserTicketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 我的票（分页）
func (l *ListUserTicketsLogic) ListUserTickets(in *ticket.ListUserTicketsReq) (*ticket.ListUserTicketsResp, error) {

	var (
		tickets []db.Ticket
		total   int64
	)

	query := l.svcCtx.DB.Model(&db.Ticket{}).
		Where("user_id = ?", in.UserId)

	if in.Status != "" {
		query = query.Where("status = ?", in.Status)
	}

	// 总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 默认分页
	page := int(in.Page)
	pageSize := int(in.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at desc").
		Find(&tickets).Error; err != nil {
		return nil, err
	}

	resp := &ticket.ListUserTicketsResp{
		Total: total,
	}

	for _, t := range tickets {
		resp.Tickets = append(resp.Tickets, &ticket.TicketInfo{
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
		})
	}

	return resp, nil
}
