// 令牌桶限流
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBucket struct {
	client redis.Cmdable
}

func NewTokenBucket(client redis.Cmdable) *TokenBucket {
	return &TokenBucket{
		client: client,
	}
}

func (t *TokenBucket) Allow(ctx context.Context, key string, capacity int, rate int) (bool, error) {

	now := time.Now().Unix()

	result, err := t.client.Eval(
		ctx,
		TokenBucketLua,
		[]string{key},
		capacity,
		rate,
		now,
	).Int()

	if err != nil {
		return false, err
	}

	return result == 1, nil
}
