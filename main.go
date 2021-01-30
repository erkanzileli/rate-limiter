package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	rate_limit_service "rate-limiter/pkg/domain/service/rate-limit-service"
	in_memory_cache_repository "rate-limiter/pkg/infra/repository/in-memory-cache-repository"
	rate_limit_rule_repository "rate-limiter/pkg/infra/repository/rate-limit-rule-repository"
)

const (
	serverAddr      = ":8080"
	appServerAddr   = "http://localhost:8081"
	cacheMode       = "in_memory"
	keyFormat       = "%s_%s"
	tooManyRequests = "Too Many Requests!"
)

var (
	cacheRepository  = in_memory_cache_repository.New()
	ruleRepository   = rate_limit_rule_repository.New()
	rateLimitService = rate_limit_service.New(cacheRepository, ruleRepository)
)

func main() {
	server := &fasthttp.Server{}

	server.Handler = reverseProxyHandler

	fmt.Println("Running on", serverAddr)
	log.Fatalln(fasthttp.ListenAndServe(serverAddr, server.Handler))
}

func reverseProxyHandler(ctx *fasthttp.RequestCtx) {
	method := string(ctx.Method())
	uri := ctx.Request.URI().RequestURI()
	routingUrl := fmt.Sprintf("%s%s", appServerAddr, uri)
	key := fmt.Sprintf(keyFormat, method, uri)

	fmt.Printf("Received %s %s\n", method, uri)

	if ok := rateLimitService.CanProceed(key); !ok {
		ctx.Response.SetBody([]byte(tooManyRequests))
		ctx.Response.SetStatusCode(429)
		fmt.Println("Too many requests!")
		return
	}

	fmt.Printf("Routing to -> %s\n", routingUrl)

	clientReq := fasthttp.AcquireRequest()
	clientResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(clientReq)
	defer fasthttp.ReleaseResponse(clientResp)

	clientReq.Header.SetMethod(method)
	clientReq.SetRequestURI(routingUrl)
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		clientReq.Header.Set(string(key), string(value))
	})

	err := fasthttp.Do(clientReq, clientResp)
	if err != nil {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.Response.SetStatusCode(500)
	}

	ctx.Response.SetBody(clientResp.Body())
	ctx.Response.SetStatusCode(clientResp.StatusCode())
	clientResp.Header.VisitAll(func(key, value []byte) {
		ctx.Response.Header.Set(string(key), string(value))
	})
}
