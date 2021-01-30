package rate_limit_rule_repository

import (
	"rate-limiter/configs"
	"rate-limiter/pkg/model"
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
