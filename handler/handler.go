package handler

import (
	"context"
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/service/rate-limit-service"
	"github.com/erkanzileli/rate-limiter/tracing/new-relic"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
)

type ReverseProxyHandler interface {
	Handle(reqCtx *fasthttp.RequestCtx)
}

type handler struct {
	client           *fasthttp.HostClient
	rateLimitService rate_limit_service.RateLimitService
}

func New(rateLimitService rate_limit_service.RateLimitService) *handler {
	return &handler{
		client: &fasthttp.HostClient{
			Addr:     configs.Config.AppConfig.GetAddresses(),
			MaxConns: 1024,
		},
		rateLimitService: rateLimitService,
	}
}

func (h *handler) Handle(reqCtx *fasthttp.RequestCtx) {
	method, path := string(reqCtx.Method()), string(reqCtx.Request.URI().RequestURI())
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

	redirectResp, err := h.redirect(ctx, reqCtx)
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

func (h *handler) redirect(ctx context.Context, reqCtx *fasthttp.RequestCtx) (*fasthttp.Response, error) {
	defer new_relic.StartSegment(ctx)

	method, requestUri := string(reqCtx.Method()), string(reqCtx.Request.RequestURI())

	zap.L().Debug("Redirecting.", zap.String("requestUri", requestUri))

	redirectReq := fasthttp.AcquireRequest()
	redirectResp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(redirectReq)

	redirectReq.Header.SetMethod(method)
	redirectReq.SetRequestURI(string(reqCtx.Request.RequestURI()))
	redirectReq.SetBody(reqCtx.Request.Body())

	reqCtx.Request.Header.Del("Connection")
	reqCtx.Request.Header.VisitAll(func(key, value []byte) {
		redirectReq.Header.Set(string(key), string(value))
	})

	err := h.client.DoTimeout(redirectReq, redirectResp, configs.Config.AppConfig.Timeout)
	if err != nil {
		return nil, err
	}

	return redirectResp, nil
}

func passResponse(ctx context.Context, reqCtx *fasthttp.RequestCtx, resp *fasthttp.Response) {
	defer new_relic.StartSegment(ctx)

	reqCtx.Response.SetBody(resp.Body())
	reqCtx.Response.SetStatusCode(resp.StatusCode())

	resp.Header.Del("Connection")
	resp.Header.VisitAll(func(key, value []byte) {
		reqCtx.Response.Header.Set(string(key), string(value))
	})
}
