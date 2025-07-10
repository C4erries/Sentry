package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Adapter struct {
	client *redis.Client
}

func NewAdapter(client *redis.Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return a.client.HGet(ctx, key, field)
}

func (a *Adapter) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return a.client.HMGet(ctx, key, fields...)
}

func (a *Adapter) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return a.client.HSet(ctx, key, values...)
}

func (a *Adapter) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return a.client.SetNX(ctx, key, value, expiration)
}

func (a *Adapter) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return a.client.ZAdd(ctx, key, members...)
}

func (a *Adapter) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return a.client.ZRemRangeByScore(ctx, key, min, max)
}

func (a *Adapter) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return a.client.ZRangeByScore(ctx, key, opt)
}

func (a *Adapter) ZCard(ctx context.Context, key string) *redis.IntCmd {
	return a.client.ZCard(ctx, key)
}

func (a *Adapter) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return a.client.Expire(ctx, key, expiration)
}

func (a *Adapter) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return a.client.Exists(ctx, keys...)
}
