package proxy

import router "github.com/silentnova42/job_vacancy_poster/service/gatway/pkg/api"

type ProxyConfig interface {
	GetUrlVacancyProxy() string
	GetUrlProfileProxy() string
	GetUrlAuthProxy() string
}

type ProxyManager struct{}

func NewProxyManager() *ProxyManager {
	return &ProxyManager{}
}

func (pm *ProxyManager) InitProxy(config ProxyConfig) ([]router.ReverseProxy, error) {
	arr := make([]router.ReverseProxy, 0)

	vacancy, err := NewVacancyProxy(config.GetUrlVacancyProxy())
	if err != nil {
		return nil, err
	}
	arr = append(arr, vacancy)

	profile, err := NewProfileProxy(config.GetUrlProfileProxy())
	if err != nil {
		return nil, err
	}
	arr = append(arr, profile)

	auth, err := NewAuthProxy(config.GetUrlAuthProxy())
	if err != nil {
		return nil, err
	}
	arr = append(arr, auth)

	return arr, nil
}
