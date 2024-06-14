package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	WebServerPort string `mapstructure:"APP_PORT"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	var err error
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err = viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg, err
}
