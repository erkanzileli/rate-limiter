package redis_cache_repository

import (
	"context"
	"github.com/erkanzileli/rate-limiter/pkg/repository"
	"github.com/go-redis/redis/v8"
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
	return r.redisClient.Incr(ctx, key.(string)).Result()
}
