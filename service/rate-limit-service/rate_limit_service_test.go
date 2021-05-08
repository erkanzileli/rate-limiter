package rate_limit_service_test

import (
	"context"
	"fmt"
	"github.com/erkanzileli/rate-limiter/model"
	testifyAssert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"regexp"
	"strings"
)

func (suite *Suite) Test_it_should_return_true_when_minimum_limit_rule_was_not_exceeded() {
	t := suite.T()
	assert := testifyAssert.New(t)

	// Given
	method, requestPath := "PUT", "/users/123/deactivate"

	rulePattern1 := "PUT /users.*"
	regex1, _ := regexp.Compile(rulePattern1)

	rulePattern2 := "PUT /users/.*/deactivate"
	regex2, _ := regexp.Compile(rulePattern2)

	rules := []model.Rule{{rulePattern1, 10, regex1}, {rulePattern2, 5, regex2}}

	suite.ruleRepository.On("GetRules").Return(rules)
	suite.cacheRepository.On("Increment", mock.Anything, mock.Anything).Return(5, nil)

	// When
	canProceed, err := suite.service.CanProceed(context.Background(), method, requestPath)

	// Then
	assert.Nil(err)
	assert.True(canProceed)

	suite.ruleRepository.AssertExpectations(t)
	suite.cacheRepository.AssertExpectations(t)

	capturedIncrementRequest := suite.cacheRepository.Calls[0].Arguments.String(1)
	assert.True(strings.HasPrefix(capturedIncrementRequest, "PUT /users/.*/deactivate"))
}

func (suite *Suite) Test_it_should_return_false_when_minimum_limit_rule_was_exceeded() {
	t := suite.T()
	assert := testifyAssert.New(t)

	// Given
	method, requestPath := "PUT", "/users/123/deactivate"

	rulePattern := "PUT /users.*"
	regex, _ := regexp.Compile(rulePattern)
	rules := []model.Rule{{rulePattern, 5, regex}}

	suite.ruleRepository.On("GetRules").Return(rules)
	suite.cacheRepository.On("Increment", mock.Anything, mock.Anything).Return(6, nil)

	// When
	canProceed, err := suite.service.CanProceed(context.Background(), method, requestPath)

	// Then
	assert.Nil(err)
	assert.False(canProceed)

	suite.ruleRepository.AssertExpectations(t)
	suite.cacheRepository.AssertExpectations(t)
}

func (suite *Suite) Test_it_should_return_true_when_cache_repository_fails() {
	t := suite.T()
	assert := testifyAssert.New(t)

	// Given
	method, requestPath := "PUT", "/users/123"

	rulePattern := "PUT /users.*"
	regex, _ := regexp.Compile(rulePattern)
	rules := []model.Rule{{rulePattern, 5, regex}}

	suite.ruleRepository.On("GetRules").Return(rules)
	suite.cacheRepository.On("Increment", mock.Anything, mock.Anything).Return(0, fmt.Errorf("cache-error"))

	// When
	canProceed, err := suite.service.CanProceed(context.Background(), method, requestPath)

	// Then
	assert.NotNil(err)
	assert.EqualValues("cache-error", err.Error())
	assert.True(canProceed)

	suite.ruleRepository.AssertExpectations(t)
	suite.cacheRepository.AssertExpectations(t)
}

func (suite *Suite) Test_it_should_return_true_when_there_are_not_any_related_rules() {
	t := suite.T()
	assert := testifyAssert.New(t)

	// Given
	method, requestPath := "PUT", "/users/123"

	rulePattern := "GET /users.*"
	regex, _ := regexp.Compile(rulePattern)
	rules := []model.Rule{{rulePattern, 5, regex}}

	suite.ruleRepository.On("GetRules").Return(rules)

	// When
	canProceed, err := suite.service.CanProceed(context.Background(), method, requestPath)

	// Then
	assert.Nil(err)
	assert.True(canProceed)

	suite.ruleRepository.AssertExpectations(t)
	suite.cacheRepository.AssertNotCalled(t, "Increment")
}
