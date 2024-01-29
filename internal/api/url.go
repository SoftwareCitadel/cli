package api

import (
	"citadel/internal/util"

	"github.com/spf13/viper"
)

var ApiBaseUrl string = "https://console.softwarecitadel.com"

func RetrieveApiBaseUrl() string {
	configDir, err := util.InitConfigDir()
	if err != nil {
		return ApiBaseUrl
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return ApiBaseUrl
	}

	consoleUrl := viper.GetString("console_url")
	if consoleUrl == "" {
		return ApiBaseUrl
	}

	return consoleUrl
}
