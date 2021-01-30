package in_memory_cache_repository

import (
	"github.com/dgraph-io/ristretto"
	"rate-limiter/pkg/infra/repository"
	"time"
)

type repo struct {
	cache *ristretto.Cache
}

func New() repository.CacheRepository {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	return &repo{cache}
}

func (s *repo) Increment(key interface{}) int {
	if value, ok := s.cache.Get(key); ok {
		newValue := value.(int) + 1
		s.Set(key, newValue)
		return newValue
	}
	s.Set(key, 1)
	return 1
}

func (s *repo) Get(key interface{}) interface{} {
	if value, ok := s.cache.Get(key); ok {
		return value
	}
	return nil
}

func (s *repo) Set(key, value interface{}) {
	s.cache.Set(key, value, 1)
}

func (s *repo) SetWithTTL(key, value interface{}, ttl time.Duration) {
	s.cache.SetWithTTL(key, value, 1, ttl)
}
