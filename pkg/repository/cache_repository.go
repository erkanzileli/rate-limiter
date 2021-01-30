package repository

type CacheRepository interface {
	Get(key interface{}) interface{}
	Set(key, value interface{})
	Increment(key interface{}) int
}
