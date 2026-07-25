package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	db "Ticket/app/event/model"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetTicketTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTicketTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTicketTypeLogic {
	return &GetTicketTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取票种详情（Order 创建订单主要调用）
func (l *GetTicketTypeLogic) GetTicketType(in *event.GetTicketTypeReq) (*event.GetTicketTypeResp, error) {

	var ticketType db.TicketType

	err := l.svcCtx.DB.
		First(&ticketType, in.TicketTypeId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.TICKET_TYPE_NOT_FOUND)
		}
		return nil, err
	}

	// 默认使用 MySQL 中的库存
	stock := ticketType.Stock

	// 如果 Redis 已初始化库存，则优先返回 Redis 库存（用于展示）
	key := fmt.Sprintf("ticket:stock:%d", ticketType.ID)

	value, err := l.svcCtx.Redis.Get(l.ctx, key)
	if err == nil {
		if s, convErr := strconv.Atoi(value); convErr == nil {
			stock = int32(s)
		}
	}

	return &event.GetTicketTypeResp{
		TicketType: &event.TicketTypeInfo{
			Id:         uint64(ticketType.ID),
			EventId:    uint64(ticketType.EventID),
			ShowId:     uint64(ticketType.ShowID),
			Name:       ticketType.Name,
			Price:      ticketType.Price,
			Stock:      int32(stock),
			MaxPerUser: int32(ticketType.MaxPerUser),
			SortOrder:  int32(ticketType.SortOrder),
		},
	}, nil
}
