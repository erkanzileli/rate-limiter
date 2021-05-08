package rate_limit_service_test

import (
	"context"
	"github.com/erkanzileli/rate-limiter/pkg/model"
	"github.com/erkanzileli/rate-limiter/pkg/repository"
	"github.com/erkanzileli/rate-limiter/pkg/repository/rate-limit-rule-repository"
	rate_limit_service "github.com/erkanzileli/rate-limiter/pkg/service/rate-limit-service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	cacheRepository *cacheRepositoryMock
	ruleRepository  *ruleRepositoryMock
	service         rate_limit_service.RateLimitService
}

func Test(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) SetupTest() {
	suite.cacheRepository = &cacheRepositoryMock{}
	suite.ruleRepository = &ruleRepositoryMock{}
	suite.service = rate_limit_service.New(suite.cacheRepository, suite.ruleRepository)
}

type cacheRepositoryMock struct {
	mock.Mock
	repository.CacheRepository
}

func (c *cacheRepositoryMock) Increment(ctx context.Context, key interface{}) (int64, error) {
	args := c.Called(ctx, key)
	return int64(args.Int(0)), args.Error(1)
}

type ruleRepositoryMock struct {
	mock.Mock
	rate_limit_rule_repository.RateLimitRuleRepository
}

func (r *ruleRepositoryMock) GetRules() []model.Rule {
	args := r.Called()
	rules := args.Get(0)
	if rules != nil {
		return rules.([]model.Rule)
	}
	return []model.Rule{}
}
