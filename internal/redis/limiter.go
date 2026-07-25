// 通用限流器(固定窗口)
package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *goredis.Client
}

func NewRateLimiter(client *goredis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
	}
}

// Allow
// true: 允许请求
// false: 超过限制
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {

	result, err := r.client.Eval(ctx, RateLimitLua, []string{key}, int(window.Seconds()), limit).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}
