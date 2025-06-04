package main

import (
	"context"
	"log"
	"os"

	pgstorage "github.com/silentnova42/job_vacancy_poster/service/vacancy/db/pg-storage"
	router "github.com/silentnova42/job_vacancy_poster/service/vacancy/pkg/api"
	"github.com/silentnova42/job_vacancy_poster/service/vacancy/pkg/server"
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

	confForDb, err := pgstorage.NewPgConf(dbConfig.GetUrlConn())
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgstorage.Connact(ctx, confForDb, 5)
	if err != nil {
		log.Fatal(err)
	}

	err = db.RunMigration(dbConfig.GetUrlConn())
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer()
	handler := router.NewHandler(db)

	if err = server.Run(":"+os.Getenv("PORT"), handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
