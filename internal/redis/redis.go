package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type IntCmd = redis.IntCmd
type BoolCmd = redis.BoolCmd
type StringCliceCmd = redis.StringSliceCmd
type Z = redis.Z
type ZRangeBy = redis.ZRangeBy

type RedisClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd
	ZAdd(ctx context.Context, key string, members ...Z) *IntCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd
	ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringCliceCmd
	ZCard(ctx context.Context, key string) *IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
}

func add() {
	var _ RedisClient = (*redis.Client)(nil)
	redis.NewClient(&redis.Options{})
}
