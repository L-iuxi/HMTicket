package redis

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

var (
	ErrLockFailed = errors.New("获取分布式锁失败")
)

type DistributedLock struct {
	client *goredis.Client
}

// 创建锁对象
func NewDistributedLock(client *goredis.Client) *DistributedLock {
	return &DistributedLock{
		client: client,
	}
}

// 获取锁
func (l *DistributedLock) Lock(ctx context.Context, key string, expire time.Duration) (string, error) {
	value := uuid.New().String()

	ok, err := l.client.SetNX(ctx, key, value, expire).Result()
	if err != nil {
		return "", err
	}

	if !ok {
		return "", ErrLockFailed
	}

	return value, nil
}

// 解锁
func (l *DistributedLock) Unlock(ctx context.Context, key string, value string) error {

	result, err := l.client.Eval(ctx, UnlockLua, []string{key}, value).Int()
	if err != nil {
		return err
	}

	if result == 0 {
		return errors.New("锁不存在或者锁不是当前持有者")
	}
	return nil
}
