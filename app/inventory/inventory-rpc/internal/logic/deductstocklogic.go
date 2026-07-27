
package logic

import (
	"context"
	"fmt"

	"Ticket/app/inventory/inventory-rpc/internal/svc"
	"Ticket/app/inventory/inventory-rpc/inventory"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockLogic {
	return &DeductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 扣库存（Redis Lua 原子扣减）
func (l *DeductStockLogic) DeductStock(in *inventory.DeductStockReq) (*inventory.DeductStockResp, error) {
	key := redis.StockKey(in.TicketTypeId)

	result, err := l.svcCtx.Redis.Eval(l.ctx, redis.DeductStockLua, []string{key}, in.Quantity)
	if err != nil {
		return &inventory.DeductStockResp{
			Success: false,
			Message: fmt.Sprintf("Redis 执行失败: %v", err),
		}, err
	}

	code, ok := result.(int64)
	if !ok {
		return &inventory.DeductStockResp{
			Success: false,
			Message: "Redis 返回值异常",
		}, nil
	}

	switch code {
	case -1:
		return &inventory.DeductStockResp{
			Success: false,
			Message: "库存未初始化",
		}, nil
	case 0:
		return &inventory.DeductStockResp{
			Success: false,
			Message: "库存不足",
		}, nil
	case 1:
		remain, _ := l.svcCtx.Redis.GetInt(l.ctx, key)
		return &inventory.DeductStockResp{
			Success:     true,
			Message:     "扣减成功",
			RemainStock: int32(remain),
		}, nil
	default:
		return &inventory.DeductStockResp{
			Success: false,
			Message: fmt.Sprintf("未知 Lua 返回值: %d", code),
		}, nil
	}
}
