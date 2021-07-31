package configs

import (
	"github.com/erkanzileli/rate-limiter/model"
	"go.uber.org/zap"
	"regexp"
)

type ruleConfig struct {
	// Scope specifies limiting indicator for a rule.
	//
	// If scope is PathScope;
	// Each requests that matches with the rule pattern will be limited according to their request paths.
	// That means all different request paths that matches the rule will increment their usage.
	//
	// If scope is RuleScope;
	// All requests that matches with the rule will be limited according to the matched rule.
	// That means all different requests that matches the rule will increment the rule's usage.
	Scope string

	// Pattern is for request-rule matching. When matching a request "METHOD /url/path" format is used.
	Pattern string

	// Limit specifies the limit of the rule.
	Limit int64
}

// compileRules compiles given rule's patterns and filters non-valid patterns
func compileRules(rules []ruleConfig) []model.Rule {
	tempRules := make([]model.Rule, 0, len(rules))

	for _, r := range rules {
		regex, err := regexp.Compile(r.Pattern)
		if err != nil {
			zap.L().Error("error compiling rule pattern into a regexp.", zap.Error(err))
			continue
		}

		tempRules = append(tempRules, model.Rule{
			Scope:   model.NewRuleScope(r.Scope),
			Pattern: r.Pattern,
			Limit:   r.Limit,
			Regex:   regex,
		})
	}

	return tempRules
}
