package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type IntCmd = redis.IntCmd
type BoolCmd = redis.BoolCmd
type StringCmd = redis.StringCmd
type StatusCmd = redis.StatusCmd
type SliceCmd = redis.SliceCmd
type StringSliceCmd = redis.StringSliceCmd
type Z = redis.Z
type ZRangeBy = redis.ZRangeBy

type RedisClient interface {
	//Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	//Get(ctx context.Context, key string) *StringCmd
	HGet(ctx context.Context, key string, field string) *StringCmd
	HMGet(ctx context.Context, key string, fields ...string) *SliceCmd
	HSet(ctx context.Context, key string, values ...interface{}) *IntCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd
	ZAdd(ctx context.Context, key string, members ...Z) *IntCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd
	ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZCard(ctx context.Context, key string) *IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	Exists(ctx context.Context, keys ...string) *IntCmd
}

func add() {
	redis.NewClient(&redis.Options{})
}
