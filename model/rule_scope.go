package model

type RuleScope string

const (
	PathScope    RuleScope = "path"
	PatternScope RuleScope = "pattern"

	defaultScope = PathScope
)

func NewRuleScope(scope string) RuleScope {
	rScope := RuleScope(scope)

	switch rScope {
	case PathScope, PatternScope:
		return rScope
	}

	return defaultScope
}
