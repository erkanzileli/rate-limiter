package in_memory_cache_repository

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/erkanzileli/rate-limiter/repository"
)

type repo struct {
	cache *ristretto.Cache
}

func New() *repo {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10 * 1024,
		MaxCost:     1 << 20,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}

	return &repo{cache}
}

func (s *repo) Increment(ctx context.Context, key interface{}) (int64, error) {
	if value, ok := s.cache.Get(key); ok {
		newValue := value.(int64) + 1
		s.cache.SetWithTTL(key, newValue, 1, repository.IncrementKeyTTL)
		return newValue, nil
	}
	s.cache.SetWithTTL(key, int64(1), 1, repository.IncrementKeyTTL)
	return 1, nil
}
