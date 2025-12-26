package utils

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HttpServerAddress string `mapstructure:"HttpServerAddress"`
	DBSOURCE          string `mapstructure:"DB_SOURCE"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return config, nil
}
