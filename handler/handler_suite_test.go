package handler_test

import (
	"context"
	"fmt"
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/handler"
	"github.com/erkanzileli/rate-limiter/service/rate-limit-service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
	"time"
)

var (
	mockAppServerAddr = "127.0.0.1:9000"

	mockRateLimiterServerAddr     = "127.0.0.1:9001"
	mockRateLimiterServerHttpAddr = fmt.Sprintf("http://%s", mockRateLimiterServerAddr)
)

type Suite struct {
	suite.Suite
	appServer         *fasthttp.Server
	rateLimiterServer *fasthttp.Server
	httpClient        *http.Client
	rateLimitService  *rateLimitServiceMock
	handler           handler.ReverseProxyHandler
}

func Test(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) SetupSuite() {
	suite.appServer = &fasthttp.Server{}
	go func() {
		if err := suite.appServer.ListenAndServe(mockAppServerAddr); err != nil {
			panic(err)
		}
	}()

	suite.rateLimiterServer = &fasthttp.Server{}
	go func() {
		if err := suite.rateLimiterServer.ListenAndServe(mockRateLimiterServerAddr); err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)
}

func (suite *Suite) SetupTest() {
	time.Sleep(1 * time.Second)
	configs.Config.AppConfig.Port = "8080"
	configs.Config.AppConfig.Hosts = []string{mockAppServerAddr}
	configs.Config.AppConfig.Timeout = time.Second

	suite.httpClient = &http.Client{}
	suite.rateLimitService = &rateLimitServiceMock{}
	suite.handler = handler.New(suite.rateLimitService)
	suite.rateLimiterServer.Handler = suite.handler.Handle
}

type rateLimitServiceMock struct {
	mock.Mock
	rate_limit_service.RateLimitService
}

func (r *rateLimitServiceMock) CanProceed(ctx context.Context, method, path string) (bool, error) {
	args := r.Called(ctx, method, path)
	return args.Bool(0), args.Error(1)
}
