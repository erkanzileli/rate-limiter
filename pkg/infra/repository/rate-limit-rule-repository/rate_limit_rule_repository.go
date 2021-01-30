package rate_limit_rule_repository

type RateLimitRuleRepository interface {
	GetRules() *map[string]int
}

type repo struct {
	rules *map[string]int
}

func New() RateLimitRuleRepository {
	// hardcoded for now, it will update somehow these rules
	return &repo{
		rules: &map[string]int{
			"POST_/suppliers/12/packages": 3,
			"POST_/suppliers/13/packages": 4,
		},
	}
}

func (r *repo) GetRules() *map[string]int {
	return r.rules
}
