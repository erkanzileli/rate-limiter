package main

import (
	"flag"
	"fmt"
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/handler"
	"github.com/erkanzileli/rate-limiter/repository"
	"github.com/erkanzileli/rate-limiter/repository/in-memory-cache-repository"
	"github.com/erkanzileli/rate-limiter/repository/rate-limit-rule-repository"
	"github.com/erkanzileli/rate-limiter/repository/redis-cache-repository"
	"github.com/erkanzileli/rate-limiter/service/rate-limit-service"
	new_relic "github.com/erkanzileli/rate-limiter/tracing/new-relic"
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFilePath = flag.String("config-file", "config.yaml", "Config file path")
)

func init() {
	flag.Parse()
	configs.InitConfigs(*configFilePath)

	err := new_relic.CreateAgent()
	if err != nil {
		panic(err)
	}
}

func main() {
	cacheRepository := createCacheRepository()
	ruleRepository := rate_limit_rule_repository.New()
	rateLimitService := rate_limit_service.New(cacheRepository, ruleRepository)
	h := handler.New(rateLimitService)
	server := createServer(h.Handle)

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
