package configs

import (
	"github.com/erkanzileli/rate-limiter/repository/rate-limit-rule-repository"
	"github.com/spf13/viper"
	"log"
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

	DefaultRuleScope string
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
