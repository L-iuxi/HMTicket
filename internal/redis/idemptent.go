package redis

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var (
	ErrDuplicateRequest = errors.New("重复请求")
)

type Idempotent struct {
	client *goredis.Client
}

func NewIdempotent(client *goredis.Client) *Idempotent {
	return &Idempotent{
		client: client,
	}
}

// Start 尝试标记请求为 pending，返回 false 表示已存在
func (i *Idempotent) Start(ctx context.Context, key string, expire time.Duration) error {
	ok, err := i.client.SetNX(ctx, key, "pending", expire).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrDuplicateRequest
	}
	return nil
}

// Get 获取请求状态："" = 不存在，pending / ok:xxx / failed
func (i *Idempotent) Get(ctx context.Context, key string) (string, error) {
	value, err := i.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// Success 标记请求成功，result 是返回给客户端的订单号等信息
func (i *Idempotent) Success(ctx context.Context, key string, result string, expire time.Duration) error {
	return i.client.Set(ctx, key, "ok:"+result, expire).Err()
}

// Failed 标记请求失败，短 TTL 避免永久占位
func (i *Idempotent) Failed(ctx context.Context, key string, expire time.Duration) error {
	return i.client.Set(ctx, key, "failed", expire).Err()
}

// Delete 删除幂等记录
func (i *Idempotent) Delete(ctx context.Context, key string) error {
	return i.client.Del(ctx, key).Err()
}
