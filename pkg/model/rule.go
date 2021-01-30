package model

type Rule struct {
	Pattern string `yaml:"pattern"`
	Limit   int    `yaml:"limit"`
}
