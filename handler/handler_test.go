package handler_test

import (
	"bytes"
	"fmt"
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/handler"
	testifyAssert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
)

func (suite *Suite) Test_it_should_redirect_when_rate_limit_was_not_exceeded() {
	t := suite.T()
	assert := testifyAssert.New(t)

	// Given
	var capturedReqUrl, capturedReqMethod, capturedReqBody, capturedHeaderValue string

	suite.appServer.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.Response.SetBody([]byte("OK"))
		ctx.Response.SetStatusCode(http.StatusOK)
		ctx.Response.Header.Set("x-key", "value")

		capturedReqUrl = string(ctx.Request.URI().Path())
		capturedReqMethod = string(ctx.Request.Header.Method())
		capturedReqBody = string(ctx.Request.Body())
		capturedHeaderValue = string(ctx.Request.Header.Peek("x-key"))
	}

	suite.rateLimitService.On("CanProceed", mock.Anything, "POST", "/users").Return(true, nil)

	reqBody := bytes.NewBufferString(`{"key":"value"}`)

	// When
	request, _ := http.NewRequest("POST", mockRateLimiterServerHttpAddr+"/users", reqBody)
	request.Header.Set("x-key", "value")
	resp, err := suite.httpClient.Do(request)

	// Then
	assert.Nil(err)

	assert.Equal("/users", capturedReqUrl)
	assert.Equal("POST", capturedReqMethod)
	assert.Equal(`{"key":"value"}`, capturedReqBody)
	assert.Equal("value", capturedHeaderValue)

	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("value", resp.Header.Get("x-key"))

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal("OK", string(body))

	suite.rateLimitService.AssertExpectations(t)
}

func (suite *Suite) Test_it_should_redirect_when_rate_limiting_is_fails() {
	t := suite.T()
	assert := testifyAssert.New(t)

	// Given
	suite.appServer.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.Response.SetStatusCode(http.StatusOK)
	}

	suite.rateLimitService.On("CanProceed", mock.Anything, "GET", "/users?name=gümüşhacıköy").
		Return(false, fmt.Errorf("cache-error"))

	// When
	resp, err := http.Get(mockRateLimiterServerHttpAddr + "/users?name=gümüşhacıköy")

	// Then
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	suite.rateLimitService.AssertExpectations(t)
}

func (suite *Suite) Test_it_should_not_redirect_when_rate_limit_was_exceeded() {
	t := suite.T()
	assert := testifyAssert.New(t)

	// Given
	suite.rateLimitService.On("CanProceed", mock.Anything, "GET", "/users").Return(false, nil)

	// When
	resp, err := http.Get(mockRateLimiterServerHttpAddr + "/users")

	// Then
	assert.Nil(err)
	assert.Equal(http.StatusTooManyRequests, resp.StatusCode)

	suite.rateLimitService.AssertExpectations(t)
}

func (suite *Suite) Test_it_should_return_error_when_redirection_fails() {
	t := suite.T()
	assert := testifyAssert.New(t)

	// Given
	configs.Config.AppConfig.Hosts = []string{"asdf"}
	suite.handler = handler.New(suite.rateLimitService)
	suite.rateLimiterServer.Handler = suite.handler.Handle

	suite.rateLimitService.On("CanProceed", mock.Anything, "GET", "/users").Return(true, nil)

	// When
	resp, err := http.Get(mockRateLimiterServerHttpAddr + "/users")

	// Then
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)

	suite.rateLimitService.AssertExpectations(t)
}
