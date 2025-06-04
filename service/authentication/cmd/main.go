package main

import (
	"log"
	"os"

	router "github.com/silentnova42/job_vacancy_poster/service/auth/pkg/api"
	"github.com/silentnova42/job_vacancy_poster/service/auth/pkg/auth"
	"github.com/silentnova42/job_vacancy_poster/service/auth/pkg/server"
)

func main() {
	server := server.NewServer()

	authService, err := auth.NewAuthService()
	if err != nil {
		log.Fatal(err)
	}

	urlProfile := os.Getenv("PROFILE_URL")

	handler, err := router.NewHandler(authService, auth.ExpForRefresh, urlProfile)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(":"+os.Getenv("PORT"), handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}
