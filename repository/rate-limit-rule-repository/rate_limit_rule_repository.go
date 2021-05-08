package rate_limit_rule_repository

import (
	model2 "github.com/erkanzileli/rate-limiter/model"
)

var Rules []model2.Rule

type RateLimitRuleRepository interface {
	GetRules() []model2.Rule
}

type repo struct {
}

func New() RateLimitRuleRepository {
	return &repo{}
}

func (r *repo) GetRules() []model2.Rule {
	return Rules
}
