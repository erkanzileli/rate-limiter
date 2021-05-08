package configs

type ruleScope string

const (
	PathScope ruleScope = "path"
	RuleScope ruleScope = "rule"
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

// IsValid is for validating given rule scope.
func (a ruleScope) IsValid() bool {
	switch a {
	case PathScope, RuleScope:
		return true
	}
	return false
}
