package redis_cache_repository

import (
	"context"
	"github.com/erkanzileli/rate-limiter/repository"
	"github.com/erkanzileli/rate-limiter/tracing/new-relic"
	"github.com/go-redis/redis/v8"
)

type repo struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *repo {
	return &repo{
		redisClient: redisClient,
	}
}

func (r *repo) Increment(ctx context.Context, key interface{}) (int64, error) {
	defer new_relic.StartSegment(ctx)

	incr := new(redis.IntCmd)

	_, err := r.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		key := key.(string)
		incr = pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, repository.IncrementKeyTTL)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}
