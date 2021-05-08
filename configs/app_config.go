package configs

import (
	"github.com/erkanzileli/rate-limiter/pkg/model"
	rate_limit_rule_repository "github.com/erkanzileli/rate-limiter/pkg/repository/rate-limit-rule-repository"
	"github.com/spf13/viper"
	"log"
	"regexp"
)

type appConfig struct {
	v *viper.Viper

	// AppServerAddr is a url with http scheme which will used to be redirect the requests from rate-limiter.
	AppServerAddr string

	// Server contains server configurations.
	Server serverConfig

	// Cache contains cache configurations.
	Cache cacheConfig

	// Algorithm contains algorithm options. Not required.
	Algorithm algorithmConfig

	Rules []ruleConfig
}

func (a *appConfig) readWithViper(shouldPanic bool) error {
	if a.v == nil {
		v := viper.New()
		v.SetConfigFile(configFilePath)
		a.v = v
	}

	err := a.v.ReadInConfig()
	if err != nil {
		if shouldPanic {
			log.Fatalf("config read error: %+v", err)
		}
		return err
	}

	err = a.v.Unmarshal(&AppConfig)
	if err != nil {
		if shouldPanic {
			log.Fatalf("config unmarshall error: %+v", err)
		}
		return err
	}

	rate_limit_rule_repository.Rules = compileRules(a.Rules)

	return nil
}

// compileRules compiles given rule's patterns and filters non-valid patterns
func compileRules(rules []ruleConfig) []model.Rule {
	tempRules := make([]model.Rule, 0)

	for _, rule := range rules {
		regex, err := regexp.Compile(rule.Pattern)
		if err != nil {
			log.Printf("error compiling rule pattern into a regexp: %+v\n", err)
			continue
		}

		tempRules = append(tempRules, model.Rule{
			//Scope:   rule.Scope,
			Pattern: rule.Pattern,
			Limit:   rule.Limit,
			Regex:   regex,
		})
	}

	return tempRules
}
