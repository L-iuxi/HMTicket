package logic

import (
	"context"
	"errors"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	db "Ticket/app/event/model"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateTicketTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTicketTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTicketTypeLogic {
	return &UpdateTicketTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改票种
func (l *UpdateTicketTypeLogic) UpdateTicketType(in *event.UpdateTicketTypeReq) (*event.CommonResp, error) {

	// 查询票种
	var ticketType db.TicketType
	if err := l.svcCtx.DB.First(&ticketType, in.TicketTypeId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.TICKET_TYPE_NOT_FOUND)
		}
		return nil, err
	}

	// 查询所属活动
	var even db.Event
	if err := l.svcCtx.DB.First(&even, ticketType.EventID).Error; err != nil {
		return nil, xerr.NewErrCode(xerr.TICKET_TYPE_EVENT_NOT_FOUND)
	}

	// 已开售或已结束禁止修改
	if even.Status == db.EventStatusSelling ||
		even.Status == db.EventStatusClosed {
		return nil, xerr.NewErrCode(xerr.TICKET_TYPE_STATUS_INVALID)
	}

	updates := make(map[string]interface{})

	if in.Name != "" {
		updates["name"] = in.Name
	}

	if in.Price > 0 {
		updates["price"] = in.Price
	}

	if in.Stock >= 0 {
		updates["stock"] = in.Stock
	}

	if in.MaxPerUser > 0 {
		updates["max_per_user"] = in.MaxPerUser
	}

	if in.SortOrder >= 0 {
		updates["sort_order"] = in.SortOrder
	}

	if len(updates) == 0 {
		return nil, xerr.NewErrCode(xerr.NO_DATA_TO_UPDATE)
	}

	if err := l.svcCtx.DB.Model(&ticketType).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 同步 Redis 库存
	if _, ok := updates["stock"]; ok {
		l.svcCtx.InventoryRpc.UpdateStock(l.ctx, &inventoryclient.UpdateStockReq{
			TicketTypeId: in.TicketTypeId,
			Stock:        in.Stock,
		})
	}

	return &event.CommonResp{}, nil
}
