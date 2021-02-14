package configs

type cacheConfig struct {
	// InMemory is a flag that represents the cache mode. If true then in-memory cache will be used.
	// If not true then a Redis connection will be open with specified Redis configs.
	InMemory bool

	// Redis includes configurations to connect a Redis as cache service.
	Redis redisConfig
}

type redisConfig struct {
	// Addr is host address of the Redis
	Addr string

	// Username is credential for connecting to Redis
	Username string

	// Password is credential for connecting to Redis
	Password string

	// DB represents Redis DB namespace
	DB int
}
