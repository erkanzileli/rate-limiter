package rate_limit_rule_repository

import (
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/erkanzileli/rate-limiter/pkg/model"
)

type RateLimitRuleRepository interface {
	GetRules() []*model.Rule
}

type repo struct {
	rules []*model.Rule
}

func New() RateLimitRuleRepository {
	return &repo{}
}

func (r *repo) GetRules() []*model.Rule {
	return configs.AppConfig.Rules
}
