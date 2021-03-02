package rate_limit_rule_repository

import (
	"github.com/erkanzileli/rate-limiter/pkg/model"
)

var Rules []*model.Rule

type RateLimitRuleRepository interface {
	GetRules() []*model.Rule
}

type repo struct {
}

func New() RateLimitRuleRepository {
	return &repo{}
}

func (r *repo) GetRules() []*model.Rule {
	return Rules
}
