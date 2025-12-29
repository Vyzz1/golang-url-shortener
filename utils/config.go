package utils

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	DBSOURCE          string `mapstructure:"DB_SOURCE"`
	BaseURL           string `mapstructure:"BASE_URL"`
	FrontendURL       string `mapstructure:"FRONTEND_URL"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
