package configs

import (
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

var (
	AppConfig      appConfig
	configFilePath string
)

func InitConfigs(configPath string) {
	configFilePath = configPath
	AppConfig.readWithViper(true)
	AppConfig.v.WatchConfig()
	AppConfig.v.OnConfigChange(func(in fsnotify.Event) {
		err := AppConfig.readWithViper(false)
		if err != nil {
			zap.L().Error("Failed to refreshing application configs due to file change.", zap.Error(err))
			return
		}
		zap.L().Info("Application configs are changed")
	})
}
