package configs

import (
	"github.com/fsnotify/fsnotify"
	"log"
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
			log.Println("Error on refreshing application configs due to file change, error: ", err)
			return
		}
		log.Println("Application configs are changed")
	})
}
