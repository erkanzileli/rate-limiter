package reverse_proxy_handler

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"rate-limiter/configs"
	"rate-limiter/pkg/service/rate-limit-service"
)

const (
	tooManyRequests = "Too Many Requests!"
)

type handler struct {
	rateLimitService rate_limit_service.RateLimitService
}

func New(rateLimitService rate_limit_service.RateLimitService) *handler {
	return &handler{
		rateLimitService: rateLimitService,
	}
}

func (h *handler) Handle(ctx *fasthttp.RequestCtx) {
	method := string(ctx.Method())
	uri := string(ctx.Request.URI().RequestURI())
	routingUrl := fmt.Sprintf("%s%s", configs.AppConfig.AppServerAddr, uri)

	fmt.Printf("Received %s %s\n", method, uri)

	if ok, err := h.rateLimitService.CanProceed(ctx, method, uri); err != nil {
		fmt.Printf("Rate limit skipping due to error: %+v", err)
	} else if !ok {
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
