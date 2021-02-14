package configs

import (
	"github.com/spf13/viper"
	"github.com/erkanzileli/rate-limiter/pkg/model"
)

type appConfig struct {
	v *viper.Viper

	// AppServerAddr is a url with http scheme which will used to be redirect the requests from rate-limiter.
	AppServerAddr string `yaml:"appServerAddr"`

	// ServerConfig includes server configurations.
	ServerConfig serverConfig `yaml:"server"`

	// CacheConfig includes cache configurations.
	CacheConfig cacheConfig `yaml:"cacheConfig"`

	// Rules is regexes and its limits to limit requests for 60 second periods.
	Rules []*model.Rule `yaml:"rules"`
}

func (a *appConfig) readWithViper(shouldPanic bool) error {
	if a.v == nil {
		v := viper.New()
		v.SetConfigFile("config.yaml")
		a.v = v
	}

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
