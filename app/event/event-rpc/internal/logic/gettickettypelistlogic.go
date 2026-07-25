package logic

import (
	"context"
	"fmt"
	"strconv"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	db "Ticket/app/event/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTicketTypeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTicketTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTicketTypeListLogic {
	return &GetTicketTypeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取场次下所有票种
func (l *GetTicketTypeListLogic) GetTicketTypeList(in *event.GetTicketTypeListReq) (*event.GetTicketTypeListResp, error) {

	var ticketTypes []db.TicketType

	if err := l.svcCtx.DB.
		Where("show_id = ?", in.ShowId).
		Order("sort_order asc").
		Find(&ticketTypes).Error; err != nil {
		return nil, err
	}

	resp := &event.GetTicketTypeListResp{}

	for _, t := range ticketTypes {

		stock := t.Stock

		key := fmt.Sprintf("ticket:stock:%d", t.ID)

		value, err := l.svcCtx.Redis.Get(l.ctx, key)
		if err == nil {

			if num, convErr := strconv.Atoi(value); convErr == nil {
				stock = int32(num)
			}
		}

		resp.TicketTypes = append(resp.TicketTypes, &event.TicketTypeInfo{
			Id:         uint64(t.ID),
			EventId:    uint64(t.EventID),
			ShowId:     uint64(t.ShowID),
			Name:       t.Name,
			Price:      t.Price,
			Stock:      int32(stock),
			MaxPerUser: int32(t.MaxPerUser),
			SortOrder:  int32(t.SortOrder),
		})
	}

	return resp, nil
}
