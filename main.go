package main

import (
	"flag"
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/handler"
	"github.com/erkanzileli/rate-limiter/repository"
	"github.com/erkanzileli/rate-limiter/repository/in-memory-cache-repository"
	"github.com/erkanzileli/rate-limiter/repository/rate-limit-rule-repository"
	"github.com/erkanzileli/rate-limiter/repository/redis-cache-repository"
	"github.com/erkanzileli/rate-limiter/service/rate-limit-service"
	"github.com/erkanzileli/rate-limiter/tracing/new-relic"
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFilePath = flag.String("config-file", "config.yaml", "Config file path")
)

func main() {
	flag.Parse()
	configs.InitConfigs(*configFilePath)

	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	zap.ReplaceGlobals(logger)

	defer logger.Sync()

	if err = new_relic.CreateAgent(); err != nil {
		logger.Fatal("failed to create new relic agent", zap.Error(err))
	}

	cacheRepository := newCacheRepository()
	ruleRepository := rate_limit_rule_repository.New()
	rateLimitService := rate_limit_service.New(cacheRepository, ruleRepository)
	h := handler.New(rateLimitService)
	server := createServer(h.Handle)

	go func() {
		if err := fasthttp.ListenAndServe(configs.AppConfig.Server.Addr, server.Handler); err != nil {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	logger.Info("Running on: " + configs.AppConfig.Server.Addr)

	handleGracefulShutdown(server)
}

func newCacheRepository() repository.CacheRepository {
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

func handleGracefulShutdown(server *fasthttp.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutting down the server...")

	if err := server.Shutdown(); err != nil {
		zap.L().Fatal("Server forced to shutdown:", zap.Error(err))
	}

	zap.L().Info("Server shut downed.")
}

func createServer(handler func(ctx *fasthttp.RequestCtx)) *fasthttp.Server {
	return &fasthttp.Server{
		Handler: handler,
		ErrorHandler: func(ctx *fasthttp.RequestCtx, err error) {
			zap.L().Error("Server error occurred.", zap.Error(err))
		},
		ReadTimeout:  time.Duration(configs.AppConfig.Server.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(configs.AppConfig.Server.WriteTimeout) * time.Millisecond,
	}
}
