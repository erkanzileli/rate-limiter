package configs

import (
	"github.com/spf13/viper"
	"rate-limiter/pkg/model"
)

type appConfig struct {
	v             *viper.Viper
	ServerAddr    string        `yaml:"serverAddr"`
	AppServerAddr string        `yaml:"appServerAddr"`
	Redis         *RedisConfig   `yaml:"redis"`
	Rules         []*model.Rule `yaml:"rules"`
}

func (a *appConfig) readWithViper(shouldPanic bool) error {
	if a.v == nil {
		v := viper.New()
		v.SetConfigFile("config.yaml")
		a.v = v
	}

	a.v.BindEnv("VERSION")

	err := a.v.ReadInConfig()
	if err != nil {
		if shouldPanic {
			panic(err)
		}
		return err
	}

	err = a.v.Unmarshal(&AppConfig)
	if err != nil {
		if shouldPanic {
			panic(err)
		}
		return err
	}

	return nil
}
