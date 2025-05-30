package main

import (
	"log"
	"os"

	config "github.com/silentnova42/job_vacancy_poster/service/gatway/configs"
	router "github.com/silentnova42/job_vacancy_poster/service/gatway/pkg/api"
	"github.com/silentnova42/job_vacancy_poster/service/gatway/pkg/server"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatal(err)
	}

	server := server.NewServer()

	proxy, err := config.InitProxy()
	if err != nil {
		log.Fatal(err)
	}

	handler, err := router.NewHandler(proxy...)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(":"+os.Getenv("PORT"), handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}
