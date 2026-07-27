package admin

import (
	"context"

	"Ticket/app/admin/internal/svc"
	"Ticket/app/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCompensationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCompensationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCompensationsLogic {
	return &ListCompensationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 查看补偿订单
func (l *ListCompensationsLogic) ListCompensations(req *types.ListCompensationsReq) (*types.ListCompensationsResp, error) {
	records, err := l.svcCtx.Redis.ListCompensations(l.ctx)
	if err != nil {
		return nil, err
	}

	items := make([]types.CompensationItem, 0, len(records))
	for _, r := range records {
		items = append(items, types.CompensationItem{
			TicketTypeID: r.TicketTypeID,
			Quantity:     r.Quantity,
			Reason:       r.Reason,
			LastOrderNo:  r.LastOrderNo,
			CreatedAt:    r.CreatedAt,
		})
	}

	return &types.ListCompensationsResp{
		Records: items,
	}, nil
}
