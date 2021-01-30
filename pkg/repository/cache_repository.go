package repository

import "time"

type CacheRepository interface {
	Get(key interface{}) interface{}
	Set(key, value interface{})
	SetWithTTL(key, value interface{}, ttl time.Duration)
	Increment(key interface{}) int
}
