package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"log"
	"rate-limiter/configs"
	"rate-limiter/pkg/repository"
	in_memory_cache_repository "rate-limiter/pkg/repository/in-memory-cache-repository"
	"rate-limiter/pkg/repository/rate-limit-rule-repository"
	redis_cache_repository "rate-limiter/pkg/repository/redis-cache-repository"
	"rate-limiter/pkg/reverse_proxy_handler"
	"rate-limiter/pkg/service/rate-limit-service"
)

func init() {
	configs.InitConfigs()
}

func main() {
	var (
		cacheRepository  = initializeCacheRepository()
		ruleRepository   = rate_limit_rule_repository.New()
		rateLimitService = rate_limit_service.New(cacheRepository, ruleRepository)
		handler          = reverse_proxy_handler.New(rateLimitService)
	)

	server := &fasthttp.Server{}
	server.Handler = handler.Handle

	log.Println("Running on", configs.AppConfig.ServerAddr)
	log.Fatalln(fasthttp.ListenAndServe(configs.AppConfig.ServerAddr, server.Handler))
}

func initializeCacheRepository() repository.CacheRepository {
	if configs.AppConfig.Redis != nil {
		redisClientOptions := redis.Options{
			Addr:     configs.AppConfig.Redis.Addr,
			Username: configs.AppConfig.Redis.Username,
			Password: configs.AppConfig.Redis.Password,
			DB:       configs.AppConfig.Redis.DB,
		}
		client := redis.NewClient(&redisClientOptions)
		return redis_cache_repository.New(client)
	}
	return in_memory_cache_repository.New()

}
