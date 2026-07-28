package admin

import (
	"context"
	"time"

	"Ticket/app/admin/internal/svc"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

// StartAutoCompensation 后台定时自动执行补偿记录。
// 这个函数设计为在 admin 启动时以 goroutine 运行。
func StartAutoCompensation(svcCtx *svc.ServiceContext) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		// 启动后先等 30 秒再跑第一次
		time.Sleep(30 * time.Second)

		for range ticker.C {
			ctx := context.Background()
			runCompensation(ctx, svcCtx.Redis)
		}
	}()
}

func runCompensation(ctx context.Context, r *redis.RedisClient) {
	records, err := r.ListCompensations(ctx)
	if err != nil {
		logx.Errorf("[AUTO-COMPENSATE] 读取补偿列表失败: %v", err)
		return
	}

	if len(records) == 0 {
		return
	}

	logx.Infof("[AUTO-COMPENSATE] 发现 %d 条待补偿记录", len(records))

	for _, rec := range records {
		stockKey := redis.StockKey(rec.TicketTypeID)

		// 先确认 key 存在
		exists, err := r.Exists(ctx, stockKey)
		if err != nil {
			logx.Errorf("[AUTO-COMPENSATE] 查询库存失败: ticketTypeId=%d, err=%v", rec.TicketTypeID, err)
			continue
		}
		if !exists {
			logx.Errorf("[AUTO-COMPENSATE] 库存 key 不存在，跳过: ticketTypeId=%d", rec.TicketTypeID)
			continue
		}

		newStock, err := r.IncrBy(ctx, stockKey, int64(rec.Quantity))
		if err != nil {
			logx.Errorf("[AUTO-COMPENSATE] 回滚失败保留记录: ticketTypeId=%d, quantity=%d, err=%v",
				rec.TicketTypeID, rec.Quantity, err)
			continue
		}

		// 回滚成功，删除记录
		if err := r.RemoveCompensation(ctx, rec.TicketTypeID); err != nil {
			logx.Errorf("[AUTO-COMPENSATE] 删除记录失败(库存已回滚): ticketTypeId=%d, err=%v",
				rec.TicketTypeID, err)
		}

		logx.Infof("[AUTO-COMPENSATE] 自动补偿成功: ticketTypeId=%d, quantity=%d, newStock=%d, orderNo=%s",
			rec.TicketTypeID, rec.Quantity, newStock, rec.LastOrderNo)
	}
}
