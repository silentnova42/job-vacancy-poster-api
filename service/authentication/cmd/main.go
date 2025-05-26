package main

import (
	"log"
	"os"

	router "github.com/silentnova42/job_vacancy_poster/pkg/api"
	auth "github.com/silentnova42/job_vacancy_poster/pkg/jwt"
	"github.com/silentnova42/job_vacancy_poster/pkg/server"
)

func main() {
	server := server.NewServer()

	authService, err := auth.NewAuthService()
	if err != nil {
		log.Fatal(err)
	}

	handler, err := router.NewHandler(authService, auth.ExpForRefresh)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(":"+os.Getenv("PORT"), handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}
