package config

import (
	"github.com/spf13/viper"
)

func InitConfig[T any](path string, conf *T) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(conf)
	if err != nil {
		return err
	}
	return nil
}
