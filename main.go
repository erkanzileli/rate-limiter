package main

import (
	"fmt"
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/pkg/handler/reverse_proxy_handler"
	"github.com/erkanzileli/rate-limiter/pkg/repository"
	"github.com/erkanzileli/rate-limiter/pkg/repository/in-memory-cache-repository"
	"github.com/erkanzileli/rate-limiter/pkg/repository/rate-limit-rule-repository"
	"github.com/erkanzileli/rate-limiter/pkg/repository/redis-cache-repository"
	"github.com/erkanzileli/rate-limiter/pkg/service/rate-limit-service"
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
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
		err := fasthttp.ListenAndServe(configs.AppConfig.Server.Addr, server.Handler)
		if err != nil {
			panic(fmt.Errorf("server error: %+v", err))
		}
	}()

	log.Println("Running on", configs.AppConfig.Server.Addr)

	handleGracefulShutdown(server)
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
		ReadTimeout:  time.Duration(configs.AppConfig.Server.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(configs.AppConfig.Server.WriteTimeout) * time.Millisecond,
	}
}

func createCacheRepository() repository.CacheRepository {
	if configs.AppConfig.Cache.InMemory {
		return in_memory_cache_repository.New()
	}
	redisClientOptions := redis.Options{
		Addr:     configs.AppConfig.Cache.Redis.Addr,
		Username: configs.AppConfig.Cache.Redis.Username,
		Password: configs.AppConfig.Cache.Redis.Password,
		DB:       configs.AppConfig.Cache.Redis.DB,
	}
	client := redis.NewClient(&redisClientOptions)
	return redis_cache_repository.New(client)
}
