package configs

import (
	"github.com/erkanzileli/rate-limiter/pkg/model"
	"github.com/spf13/viper"
	"log"
	"regexp"
)

type appConfig struct {
	v *viper.Viper

	// AppServerAddr is a url with http scheme which will used to be redirect the requests from rate-limiter.
	AppServerAddr string

	// Server includes server configurations.
	Server serverConfig

	// Cache includes cache configurations.
	Cache cacheConfig

	// Rules is regexes and its limits to limit requests for 60 second periods.
	Rules []*model.Rule
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

	a.Rules = compileRules(a.Rules)

	return nil
}

// compileRules compiles given rule's patterns and filters non-valid patterns
func compileRules(rules []*model.Rule) []*model.Rule {
	tempRules := make([]*model.Rule, 0)

	for _, rule := range rules {
		regex, err := regexp.Compile(rule.Pattern)
		if err != nil {
			log.Printf("error compiling rule pattern into a regexp: %+v\n", err)
			continue
		}

		rule.Regex = regex
		tempRules = append(tempRules, rule)
	}

	return tempRules
}
