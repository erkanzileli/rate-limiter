package configs

type cacheConfig struct {
	// InMemory is a flag that represents the cache mode. If true then in-memory cache will be used.
	// If not true then a Redis connection will be open with specified Redis configs.
	InMemory bool `yaml:"inMemory"`

	// Redis includes configurations to connect a Redis as cache service.
	Redis redisConfig `yaml:"redis"`
}

type redisConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
