package admin

import (
	"context"
	"fmt"

	"Ticket/app/admin/internal/svc"
	"Ticket/app/admin/internal/types"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteCompensationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExecuteCompensationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteCompensationLogic {
	return &ExecuteCompensationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Redis INCRBY 回滚库存
func (l *ExecuteCompensationLogic) ExecuteCompensation(req *types.ExecuteCompensationReq) (*types.ExecuteCompensationResp, error) {
	stockKey := redis.StockKey(req.TicketTypeID)

	// 先确认库存 key 存在
	exists, err := l.svcCtx.Redis.Exists(l.ctx, stockKey)
	if err != nil {
		return failResp(fmt.Sprintf("查询库存失败: %v", err)), nil
	}
	if !exists {
		return failResp("库存 key 不存在，请确认 ticketTypeId 正确"), nil
	}

	// 从补偿记录中读取数量
	compKey := fmt.Sprintf("%s%d", redis.CompensatePrefix, req.TicketTypeID)
	vals, err := l.svcCtx.Redis.HMGet(l.ctx, compKey, "quantity")
	if err != nil || vals[0] == nil {
		return failResp("未找到补偿记录"), nil
	}

	var quantity int64
	fmt.Sscanf(fmt.Sprintf("%v", vals[0]), "%d", &quantity)
	if quantity <= 0 {
		return failResp("补偿数量无效"), nil
	}

	// 直接 Redis INCRBY 回滚
	newStock, err := l.svcCtx.Redis.IncrBy(l.ctx, stockKey, quantity)
	if err != nil {
		return failResp(fmt.Sprintf("库存回滚失败: %v", err)), nil
	}

	// 删除补偿记录
	if err := l.svcCtx.Redis.RemoveCompensation(l.ctx, req.TicketTypeID); err != nil {
		l.Errorf("删除补偿记录失败(库存已回滚): ticketTypeId=%d, err=%v", req.TicketTypeID, err)
	}

	l.Infof("[COMPENSATE] 手动回滚库存成功: ticketTypeId=%d, quantity=%d, newStock=%d",
		req.TicketTypeID, quantity, newStock)

	return &types.ExecuteCompensationResp{
		Success:  true,
		Message:  fmt.Sprintf("回滚成功，新库存: %d", newStock),
		NewStock: newStock,
	}, nil
}

func failResp(msg string) *types.ExecuteCompensationResp {
	return &types.ExecuteCompensationResp{
		Success: false,
		Message: msg,
	}
}
