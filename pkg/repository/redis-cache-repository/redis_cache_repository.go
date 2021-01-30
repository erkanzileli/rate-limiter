package redis_cache_repository

import (
	"github.com/go-redis/redis/v8"
	"rate-limiter/pkg/repository"
)

type repo struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) repository.CacheRepository {
	return &repo{
		redisClient: redisClient,
	}
}

func (s *repo) Set(key, value interface{}) {
	panic("implement me")
}

func (s *repo) Increment(key interface{}) int {
	panic("implement me")
}

func (s *repo) Get(key interface{}) interface{} {
	panic("implement me")
}
