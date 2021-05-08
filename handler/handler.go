package handler

import (
	"context"
	"fmt"
	"github.com/erkanzileli/rate-limiter/configs"
	rate_limit_service2 "github.com/erkanzileli/rate-limiter/service/rate-limit-service"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"time"
)

type ReverseProxyHandler interface {
	Handle(reqCtx *fasthttp.RequestCtx)
}

type handler struct {
	rateLimitService rate_limit_service2.RateLimitService
}

func New(rateLimitService rate_limit_service2.RateLimitService) *handler {
	return &handler{
		rateLimitService: rateLimitService,
	}
}

func (h *handler) Handle(reqCtx *fasthttp.RequestCtx) {
	method, path := string(reqCtx.Method()), string(reqCtx.Request.URI().Path())
	ctx := context.Background()

	log.Printf("Received %s %s\n", method, path)

	if ok, err := h.rateLimitService.CanProceed(ctx, method, path); err != nil {
		log.Printf("Rate limit skipping due to error: %+v", err)
	} else if !ok {
		reqCtx.Response.SetStatusCode(http.StatusTooManyRequests)
		log.Println("Too many requests!")
		return
	}

	redirectResp, err := redirect(reqCtx)
	if err != nil {
		reqCtx.Response.SetBody([]byte(err.Error()))
		reqCtx.Response.SetStatusCode(http.StatusInternalServerError)
	}

	if redirectResp == nil {
		return
	}

	defer fasthttp.ReleaseResponse(redirectResp)
	passResponse(reqCtx, redirectResp)
}

func redirect(reqCtx *fasthttp.RequestCtx) (*fasthttp.Response, error) {
	method, routingUrl := string(reqCtx.Method()), getRoutingUrl(string(reqCtx.Request.URI().Path()))

	log.Printf("Redirecting to -> %s\n", routingUrl)

	redirectReq := fasthttp.AcquireRequest()
	redirectResp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(redirectReq)

	redirectReq.Header.SetMethod(method)
	redirectReq.SetRequestURI(routingUrl)
	redirectReq.SetBody(reqCtx.Request.Body())
	reqCtx.Request.Header.VisitAll(func(key, value []byte) {
		redirectReq.Header.Set(string(key), string(value))
	})

	err := fasthttp.DoTimeout(redirectReq, redirectResp, 1*time.Second)
	if err != nil {
		return nil, err
	}

	return redirectResp, nil
}

func passResponse(reqCtx *fasthttp.RequestCtx, resp *fasthttp.Response) {
	reqCtx.Response.SetBody(resp.Body())
	reqCtx.Response.SetStatusCode(resp.StatusCode())
	resp.Header.VisitAll(func(key, value []byte) {
		reqCtx.Response.Header.Set(string(key), string(value))
	})
}

func getRoutingUrl(path string) string {
	return fmt.Sprintf("%s%s", configs.AppConfig.AppServerAddr, path)
}
