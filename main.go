package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"rate-limiter/configs"
	"rate-limiter/pkg/repository"
	"rate-limiter/pkg/repository/in-memory-cache-repository"
	"rate-limiter/pkg/repository/rate-limit-rule-repository"
	"rate-limiter/pkg/repository/redis-cache-repository"
	"rate-limiter/pkg/reverse_proxy_handler"
	"rate-limiter/pkg/service/rate-limit-service"
	"syscall"
	"time"
)

func init() {
	configs.InitConfigs()
}

func main() {
	cacheRepository := createCacheRepository()
	ruleRepository := rate_limit_rule_repository.New()
	rateLimitService := rate_limit_service.New(cacheRepository, ruleRepository)
	handler := reverse_proxy_handler.New(rateLimitService)
	server := createServer(handler.Handle)

	go func() {
		err := fasthttp.ListenAndServe(configs.AppConfig.ServerConfig.Addr, server.Handler)
		if err != nil {
			panic(fmt.Errorf("server error: %+v", err))
		}
	}()

	handleGracefulShutdown(server)

	log.Println("Running on", configs.AppConfig.ServerConfig.Addr)
}

func handleGracefulShutdown(server *fasthttp.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server...")

	if err := server.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server shut downed.")
}

func createServer(handler func(ctx *fasthttp.RequestCtx)) *fasthttp.Server {
	return &fasthttp.Server{
		Handler: handler,
		ErrorHandler: func(ctx *fasthttp.RequestCtx, err error) {
			log.Printf("Server error occurred %+v", err)
		},
		ReadTimeout:  time.Duration(configs.AppConfig.ServerConfig.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(configs.AppConfig.ServerConfig.WriteTimeout) * time.Millisecond,
	}
}

func createCacheRepository() repository.CacheRepository {
	if configs.AppConfig.CacheConfig.InMemory {
		return in_memory_cache_repository.New()
	}
	redisClientOptions := redis.Options{
		Addr:     configs.AppConfig.CacheConfig.Redis.Addr,
		Username: configs.AppConfig.CacheConfig.Redis.Username,
		Password: configs.AppConfig.CacheConfig.Redis.Password,
		DB:       configs.AppConfig.CacheConfig.Redis.DB,
	}
	client := redis.NewClient(&redisClientOptions)
	return redis_cache_repository.New(client)
}
