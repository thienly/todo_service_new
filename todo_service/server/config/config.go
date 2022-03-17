package config

import (
	"github.com/spf13/viper"
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
func LoadConfig(path string) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(PanicErr)
	}
	cf := &Config{}
	err := viper.Unmarshal(cf)
	if err != nil {
		return nil, err
	}
	return cf, nil
}
