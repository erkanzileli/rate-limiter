package model

import "regexp"

type Rule struct {
	Regex   *regexp.Regexp
	Pattern string
	Limit   int64
}
