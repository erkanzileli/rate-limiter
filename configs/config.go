package configs

import (
	rate_limit_rule_repository "github.com/erkanzileli/rate-limiter/repository/rate-limit-rule-repository"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Config         config
	configFilePath string
)

type config struct {
	v *viper.Viper

	// AppConfig is target application's settings.
	AppConfig appConfig

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

func (a *config) readWithViper(shouldPanic bool) error {
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

	err = a.v.Unmarshal(&Config)
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

func InitConfigs(configPath string) {
	configFilePath = configPath
	Config.readWithViper(true)
	Config.v.WatchConfig()
	Config.v.OnConfigChange(func(in fsnotify.Event) {
		err := Config.readWithViper(false)
		if err != nil {
			zap.L().Error("Failed to refreshing application configs due to file change.", zap.Error(err))
			return
		}
		zap.L().Info("Application configs are changed")
	})
}
