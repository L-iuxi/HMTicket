package redis

import (
	"context"
	"fmt"
	"time"
)

// CompensationRecord 待补偿库存记录
type CompensationRecord struct {
	TicketTypeID uint64 `json:"ticketTypeId"`
	Quantity     int32  `json:"quantity"`
	Reason       string `json:"reason"`
	LastOrderNo  string `json:"lastOrderNo"`
	CreatedAt    string `json:"createdAt"`
}

func compensateKey(ticketTypeID uint64) string {
	return fmt.Sprintf("%s%d", CompensatePrefix, ticketTypeID)
}

// RecordCompensation 记录一条补偿任务
// 写入 Hash + 加入待处理 Set
func (r *RedisClient) RecordCompensation(ctx context.Context, ticketTypeID uint64, quantity int32, reason, orderNo string) error {
	key := compensateKey(ticketTypeID)
	now := time.Now().Format(time.RFC3339)

	if err := r.HMSet(ctx, key, map[string]interface{}{
		"quantity":      quantity,
		"reason":        reason,
		"last_order_no": orderNo,
		"created_at":    now,
	}); err != nil {
		return fmt.Errorf("记录补偿任务失败: %w", err)
	}

	// 不过期，等人工处理完再删
	if err := r.SAdd(ctx, CompensateListKey, ticketTypeID); err != nil {
		return fmt.Errorf("加入补偿列表失败: %w", err)
	}

	return nil
}

// ListCompensations 列出所有待补偿记录
func (r *RedisClient) ListCompensations(ctx context.Context) ([]CompensationRecord, error) {
	ids, err := r.SMembers(ctx, CompensateListKey)
	if err != nil {
		return nil, fmt.Errorf("读取补偿列表失败: %w", err)
	}

	records := make([]CompensationRecord, 0, len(ids))
	for _, idStr := range ids {
		var ticketTypeID uint64
		fmt.Sscanf(idStr, "%d", &ticketTypeID)

		key := compensateKey(ticketTypeID)
		vals, err := r.HMGet(ctx, key, "quantity", "reason", "last_order_no", "created_at")
		if err != nil {
			continue
		}

		var quantity int64
		fmt.Sscanf(coalesceStr(vals[0]), "%d", &quantity)

		records = append(records, CompensationRecord{
			TicketTypeID: ticketTypeID,
			Quantity:     int32(quantity),
			Reason:       coalesceStr(vals[1]),
			LastOrderNo:  coalesceStr(vals[2]),
			CreatedAt:    coalesceStr(vals[3]),
		})
	}

	return records, nil
}

// RemoveCompensation 补偿执行成功后删除记录
func (r *RedisClient) RemoveCompensation(ctx context.Context, ticketTypeID uint64) error {
	if err := r.Del(ctx, compensateKey(ticketTypeID)); err != nil {
		return err
	}
	return r.SRem(ctx, CompensateListKey, ticketTypeID)
}

func coalesceStr(v interface{}) string {
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
