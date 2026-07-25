package logic

import (
	"context"

	"Ticket/app/event/event-rpc/event"
	"Ticket/app/event/event-rpc/internal/svc"
	db "Ticket/app/event/model"
	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateTicketTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTicketTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTicketTypeLogic {
	return &CreateTicketTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建票种
func (l *CreateTicketTypeLogic) CreateTicketType(in *event.CreateTicketTypeReq) (*event.CommonResp, error) {
	var even db.Event
	//先检查活动是否存在
	if err := l.svcCtx.DB.First(&even, in.EventId).Error; err != nil {
		return nil, err
	}

	//参数校验
	if in.Price < 0 {
		return nil, xerr.NewErrMsg("票价不能小于0")
	}
	if in.Stock < 0 {
		return nil, xerr.NewErrMsg("票数不能小于0")
	}
	if in.Name == "" {
		return nil, xerr.NewErrMsg("票名不能为空字符串")
	}
	if in.MaxPerUser <= 0 {
		return nil, xerr.NewErrMsg("限购数必须大于0")
	}
	//检查场次是否存在
	var show db.Show

	if err := l.svcCtx.DB.First(&show, in.ShowId).Error; err != nil {
		return nil, xerr.NewErrCode(xerr.SHOW_NOT_FOUND)
	}
	// 场次是否属于这个活动
	if show.EventID != in.EventId {
		return nil, xerr.NewErrCode(xerr.SHOW_EVENT_NOT_MATCH)
	}
	ticketType := &db.TicketType{
		EventID:    in.EventId,
		ShowID:     in.ShowId,
		Name:       in.Name,
		Price:      in.Price,
		Stock:      in.Stock,
		MaxPerUser: in.MaxPerUser,
		SortOrder:  in.SortOrder,
	}

	//用grom事务实现，保证创建票种和更新总库存操作一起成功或失败
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(ticketType).Error; err != nil {
			return err
		}

		if err := tx.Model(&even).Update("total_stock", gorm.Expr("total_stock + ?", in.Stock)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 初始化 Redis 库存
	l.svcCtx.InventoryRpc.InitStock(l.ctx, &inventoryclient.InitStockReq{
		TicketTypeId: uint64(ticketType.ID),
		Stock:        in.Stock,
	})

	return &event.CommonResp{}, nil
}
