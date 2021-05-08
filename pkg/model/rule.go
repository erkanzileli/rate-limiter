package model

import "regexp"

// Rule model is represents a valid rule that can be used for rate limiting.
type Rule struct {
	//Scope   string
	Pattern string
	Limit   int64
	Regex   *regexp.Regexp
}
