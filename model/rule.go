package model

import "regexp"

// Rule model is represents a valid rule that can be used for rate limiting.
type Rule struct {
	Scope   RuleScope
	Pattern string
	Limit   int64
	Regex   *regexp.Regexp
}

func (r Rule) IsPathScope() bool {
	return r.Scope == PathScope
}

func (r Rule) IsPatternScope() bool {
	return r.Scope == PatternScope
}
