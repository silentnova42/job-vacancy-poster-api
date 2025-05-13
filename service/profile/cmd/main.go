package main

import (
	"context"
	"log"
	"os"

	pgstorage "github.com/silentnova42/job_vacancy_poster/db/pg-storage"
	ginrouter "github.com/silentnova42/job_vacancy_poster/pkg/gin-router"
	"github.com/silentnova42/job_vacancy_poster/pkg/server"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()

	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	config := pgstorage.Config{
		Username: viper.GetString("db.username"),
		Host:     viper.GetString("db.host"),
		Password: os.Getenv("POSTGRES_PSASWORD"),
		Port:     viper.GetString("db.port"),
		Dbname:   viper.GetString("db.dbname"),
	}

	confDb, err := pgstorage.NewPgConf(config.GetUrlConn())
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgstorage.Connect(ctx, confDb, 5)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.RunMigrate(config.GetUrlConn()); err != nil {
		log.Fatal(err)
	}

	handler := ginrouter.NewHandler(db)
	if err := server.NewServer().Run(":"+os.Getenv("PORT"), handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
