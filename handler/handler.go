package handler

import (
	"context"
	"fmt"
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/service/rate-limit-service"
	"github.com/erkanzileli/rate-limiter/tracing/new-relic"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ReverseProxyHandler interface {
	Handle(reqCtx *fasthttp.RequestCtx)
}

type handler struct {
	rateLimitService rate_limit_service.RateLimitService
}

func New(rateLimitService rate_limit_service.RateLimitService) *handler {
	return &handler{
		rateLimitService: rateLimitService,
	}
}

func (h *handler) Handle(reqCtx *fasthttp.RequestCtx) {
	method, path := string(reqCtx.Method()), string(reqCtx.Request.URI().Path())
	ctx := context.Background()
	ctx, endTxn := new_relic.StartTransaction(ctx, method+" "+path)
	defer endTxn()

	zapPath, zapMethod := zap.String("path", path), zap.String("method", method)

	zap.L().Debug("Received request.", zapPath, zapMethod)

	if ok, err := h.rateLimitService.CanProceed(ctx, method, path); err != nil {
		zap.L().Error("Rate limiting skipping due to error.", zap.Error(err), zapPath, zapMethod)
	} else if !ok {
		reqCtx.Response.SetStatusCode(http.StatusTooManyRequests)
		zap.L().Debug("Too many requests!", zap.String("method", method), zap.String("path", path))
		return
	}

	redirectResp, err := redirect(ctx, reqCtx)
	if err != nil {
		reqCtx.Response.SetBody([]byte(err.Error()))
		reqCtx.Response.SetStatusCode(http.StatusInternalServerError)
	}

	if redirectResp == nil {
		return
	}

	defer fasthttp.ReleaseResponse(redirectResp)
	passResponse(ctx, reqCtx, redirectResp)
}

func redirect(ctx context.Context, reqCtx *fasthttp.RequestCtx) (*fasthttp.Response, error) {
	defer new_relic.StartSegment(ctx)

	method, routingUrl := string(reqCtx.Method()), getRoutingUrl(string(reqCtx.Request.URI().Path()))

	zap.L().Debug("Redirecting.", zap.String("routingUrl", routingUrl))

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

func passResponse(ctx context.Context, reqCtx *fasthttp.RequestCtx, resp *fasthttp.Response) {
	defer new_relic.StartSegment(ctx)

	reqCtx.Response.SetBody(resp.Body())
	reqCtx.Response.SetStatusCode(resp.StatusCode())
	resp.Header.VisitAll(func(key, value []byte) {
		reqCtx.Response.Header.Set(string(key), string(value))
	})
}

func getRoutingUrl(path string) string {
	return fmt.Sprintf("%s%s", configs.AppConfig.AppServerAddr, path)
}
