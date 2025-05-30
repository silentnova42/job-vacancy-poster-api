package config

import (
	"github.com/spf13/viper"
)

var (
	VacancyProxy = viper.GetString("service.vacancy")
	ProfileProxy = viper.GetString("service.profile")
	AuthProxy    = viper.GetString("service.auth")
)

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
