package redis_cache_repository

import (
	"context"
	"github.com/erkanzileli/rate-limiter/pkg/repository"
	"github.com/go-redis/redis/v8"
	"time"
)

type repo struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) repository.CacheRepository {
	return &repo{
		redisClient: redisClient,
	}
}

func (r *repo) Increment(ctx context.Context, key interface{}) (int64, error) {
	var incr *redis.IntCmd

	_, err := r.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		key := key.(string)
		incr = pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute*2)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}
