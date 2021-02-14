package rate_limit_rule_repository

import (
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/pkg/model"
)

// todo: future work
// rules can be obtained from another resource
type RateLimitRuleRepository interface {
	GetRules() []*model.Rule
}

type repo struct {
}

func New() RateLimitRuleRepository {
	return &repo{}
}

func (r *repo) GetRules() []*model.Rule {
	return configs.AppConfig.Rules
}
