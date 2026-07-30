package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConf 统一 Redis 配置
type RedisConf struct {
	MasterName    string
	SentinelAddrs []string
	Password      string
	DB            int
}

type RedisClient struct {
	client redis.UniversalClient
}

// NewRedisClient 通过 Sentinel 连接 Redis，go-redis 自动处理主从发现和故障转移
func NewRedisClient(c RedisConf) (*RedisClient, error) {
	if c.MasterName == "" || len(c.SentinelAddrs) == 0 {
		return nil, fmt.Errorf("MasterName and SentinelAddrs are required")
	}
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    c.MasterName,
		SentinelAddrs: c.SentinelAddrs,
		Password:      c.Password,
		DB:            c.DB,
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  5 * time.Second,
		PoolSize:      100,
		MinIdleConns:  10,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

// 限定某个key在windows时间内最多访问limit次
func (r *RedisClient) RateLimit(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	now := time.Now().Unix()
	pipeline := r.client.Pipeline()

	pipeline.Incr(ctx, key)
	pipeline.ExpireAt(ctx, key, time.Unix(now+int64(window.Seconds()), 0))

	result, err := pipeline.Exec(ctx)
	if err != nil {
		return false, err
	}

	count, _ := result[0].(*redis.IntCmd).Result()
	return count <= limit, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) GetInt(ctx context.Context, key string) (int64, error) {
	return r.client.Get(ctx, key).Int64()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Client() redis.Cmdable {
	return r.client
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// 往redis的消息队列里面写一条新消息
func (r *RedisClient) XAdd(ctx context.Context, stream string, values map[string]interface{}) error {
	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Err()
}

// XReadGroup 从 Stream 读取消息
func (r *RedisClient) XReadGroup(ctx context.Context, group, consumer string, streams []string, count int64, block time.Duration) ([]redis.XStream, error) {
	return r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  streams,
		Count:    count,
		Block:    block,
	}).Result()
}

// XAck 确认消息
func (r *RedisClient) XAck(ctx context.Context, stream, group string, IDs ...string) error {
	return r.client.XAck(ctx, stream, group, IDs...).Err()
}

// XGroupCreate 创建消费者组
func (r *RedisClient) XGroupCreate(ctx context.Context, stream, group string) error {
	return r.client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
}

// HIncrBy Hash 字段自增
func (r *RedisClient) HIncrBy(ctx context.Context, key, field string, incr int64) error {
	return r.client.HIncrBy(ctx, key, field, incr).Err()
}

// SAdd 添加元素到 Set
func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SRem 移除 Set 元素
func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// HGet 获取 Hash 字段
func (r *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HSet 设置 Hash 字段
func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

func (r *RedisClient) HMSet(ctx context.Context, key string, values map[string]interface{}) error {
	return r.client.HSet(ctx, key, values).Err()
}

func (r *RedisClient) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return r.client.HMGet(ctx, key, fields...).Result()
}

// Pipeline 获取 pipeline
func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// Incr 自增
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

func (r *RedisClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.DecrBy(ctx, key, value).Result()
}

// Expire 设置过期时间
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// lua
func (r *RedisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}

// 分布式锁
func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}
func (r *RedisClient) Unlock(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
