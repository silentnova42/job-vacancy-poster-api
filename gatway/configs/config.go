package config

import (
	"github.com/spf13/viper"
)

type ProxyConfig struct {
	vacancyProxy string
	profileProxy string
	authProxy    string
}

func InitProxyConfig() (*ProxyConfig, error) {
	if err := initViperConfig(); err != nil {
		return nil, err
	}

	return &ProxyConfig{
		vacancyProxy: viper.GetString("service.vacancy"),
		profileProxy: viper.GetString("service.profile"),
		authProxy:    viper.GetString("service.auth"),
	}, nil
}

func initViperConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func (pc *ProxyConfig) GetUrlVacancyProxy() string {
	return pc.vacancyProxy
}

func (pc *ProxyConfig) GetUrlProfileProxy() string {
	return pc.profileProxy
}

func (pc *ProxyConfig) GetUrlAuthProxy() string {
	return pc.authProxy
}
