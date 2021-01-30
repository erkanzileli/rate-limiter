package redis_cache_repository

import (
	"rate-limiter/pkg/infra/repository"
	"time"
)

type repo struct {
}

func New() repository.CacheRepository {
	return &repo{}
}

func (s *repo) SetWithTTL(key, value interface{}, ttl time.Duration) {
	panic("implement me")
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
