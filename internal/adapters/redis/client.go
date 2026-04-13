package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	LPush(ctx context.Context, key string, value interface{}) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LLen(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type redisClient struct {
	rdb *redis.Client
}

func NewRedisClient(addr, password string, db int) Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &redisClient{rdb: rdb}
}

func (c *redisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *redisClient) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *redisClient) LPush(ctx context.Context, key string, value interface{}) error {
	return c.rdb.LPush(ctx, key, value).Err()
}

func (c *redisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.rdb.LRange(ctx, key, start, stop).Result()
}

func (c *redisClient) LLen(ctx context.Context, key string) (int64, error) {
	return c.rdb.LLen(ctx, key).Result()
}

func (c *redisClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.rdb.Expire(ctx, key, ttl).Err()
}

func (c *redisClient) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
