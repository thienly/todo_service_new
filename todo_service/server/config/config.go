package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Db Db `mapstructure:"db"`
}
type Db struct {
	Conn string `mapstructure:"conn"`
}
var (
	PanicErr = "Can not load the application config"
)
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	path, err := os.Getwd()
	if err != nil {
		panic(PanicErr)
	}
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(PanicErr)
	}
	cf := &Config{}
	err = viper.Unmarshal(cf)
	if err != nil {
		return nil, err
	}
	return cf, nil
}
