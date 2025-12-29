package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	DBSOURCE          string `mapstructure:"DB_SOURCE"`
	BaseURL           string `mapstructure:"BASE_URL"`
	FrontendURL       string `mapstructure:"FRONTEND_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path + "/.env")
	viper.SetConfigType("env")

	_ = viper.ReadInConfig()

	viper.AutomaticEnv()
	viper.BindEnv("HTTP_SERVER_ADDRESS")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("BASE_URL")
	viper.BindEnv("FRONTEND_URL")

	err = viper.Unmarshal(&config)
	return
}
