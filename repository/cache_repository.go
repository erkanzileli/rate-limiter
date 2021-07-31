package repository

import (
	"context"
)

type CacheRepository interface {
	Increment(ctx context.Context, key interface{}) (int64, error)
}
