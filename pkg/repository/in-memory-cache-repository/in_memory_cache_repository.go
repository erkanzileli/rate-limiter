package in_memory_cache_repository

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"rate-limiter/pkg/repository"
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

func (s *repo) Increment(ctx context.Context, key interface{}) (int64, error) {
	if value, ok := s.cache.Get(key); ok {
		newValue := value.(int64) + 1
		s.cache.Set(key, newValue, 1)
		return newValue, nil
	}
	s.cache.Set(key, int64(1), 1)
	return 1, nil
}
