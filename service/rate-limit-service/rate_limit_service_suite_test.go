package rate_limit_service_test

import (
	"context"
	model2 "github.com/erkanzileli/rate-limiter/model"
	repository2 "github.com/erkanzileli/rate-limiter/repository"
	rate_limit_rule_repository2 "github.com/erkanzileli/rate-limiter/repository/rate-limit-rule-repository"
	rate_limit_service2 "github.com/erkanzileli/rate-limiter/service/rate-limit-service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	cacheRepository *cacheRepositoryMock
	ruleRepository  *ruleRepositoryMock
	service         rate_limit_service2.RateLimitService
}

func Test(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) SetupTest() {
	suite.cacheRepository = &cacheRepositoryMock{}
	suite.ruleRepository = &ruleRepositoryMock{}
	suite.service = rate_limit_service2.New(suite.cacheRepository, suite.ruleRepository)
}

type cacheRepositoryMock struct {
	mock.Mock
	repository2.CacheRepository
}

func (c *cacheRepositoryMock) Increment(ctx context.Context, key interface{}) (int64, error) {
	args := c.Called(ctx, key)
	return int64(args.Int(0)), args.Error(1)
}

type ruleRepositoryMock struct {
	mock.Mock
	rate_limit_rule_repository2.RateLimitRuleRepository
}

func (r *ruleRepositoryMock) GetRules() []model2.Rule {
	args := r.Called()
	rules := args.Get(0)
	if rules != nil {
		return rules.([]model2.Rule)
	}
	return []model2.Rule{}
}
