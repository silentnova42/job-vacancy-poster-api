package main

import (
	"context"
	"log"
	"os"

	ginrouter "github.com/silentnova42/job_vacancy_poster/pkg/gin-router"
	pgstorage "github.com/silentnova42/job_vacancy_poster/pkg/pg-storage"
	"github.com/silentnova42/job_vacancy_poster/server"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	dbConfig := pgstorage.Config{
		Username: viper.GetString("db.username"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Dbname:   viper.GetString("db.dbname"),
	}

	db, err := pgstorage.NewPgDb(ctx, dbConfig.GetUrlConn(), 5)
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer()
	handler := ginrouter.NewHandler(db)

	if err := server.Run(":"+viper.GetString("port"), handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
