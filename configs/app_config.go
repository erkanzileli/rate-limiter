package configs

import (
	"github.com/erkanzileli/rate-limiter/repository/rate-limit-rule-repository"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type appConfig struct {
	v *viper.Viper

	// AppServerAddr is an url with http scheme which will use to for redirecting the requests from rate-limiter.
	AppServerAddr string

	// Server contains server configurations.
	Server serverConfig

	// Cache contains cache configurations.
	Cache cacheConfig

	// Algorithm contains algorithm options. Not required.
	Algorithm algorithmConfig

	// Rules are basically rate limiting rules.
	Rules []ruleConfig

	DefaultRuleScope string

	// Tracing contains tracing options. Only NewRelic is supported for now.
	Tracing tracingConfig
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
			zap.L().Fatal("config read error", zap.Error(err))
		}
		return err
	}

	err = a.v.Unmarshal(&AppConfig)
	if err != nil {
		if shouldPanic {
			zap.L().Fatal("config unmarshall error", zap.Error(err))
		}
		return err
	}

	rate_limit_rule_repository.Rules = compileRules(a.Rules)

	a.Tracing.validateProvider()

	return nil
}
