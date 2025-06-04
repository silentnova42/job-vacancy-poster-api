package main

import (
	"log"
	"os"

	config "github.com/silentnova42/job_vacancy_poster/service/gatway/configs"
	router "github.com/silentnova42/job_vacancy_poster/service/gatway/pkg/api"
	"github.com/silentnova42/job_vacancy_poster/service/gatway/pkg/proxy"
	"github.com/silentnova42/job_vacancy_poster/service/gatway/pkg/server"
)

func main() {
	server := server.NewServer()

	proxyConfig, err := config.InitProxyConfig()
	if err != nil {
		log.Fatal(err)
	}

	proxies, err := proxy.NewProxyManager().InitProxy(proxyConfig)
	if err != nil {
		log.Fatal(err)
	}

	handler, err := router.NewHandler(proxies)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(":"+os.Getenv("PORT"), handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}
