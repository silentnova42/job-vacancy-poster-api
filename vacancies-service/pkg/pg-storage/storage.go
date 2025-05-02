package pgstorage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	client *pgxpool.Pool
}

func NewPgDb(cxt context.Context, url string, attempts int) (*Db, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 15
	config.MinConns = 5
	config.MaxConnIdleTime = 30 * time.Minute

	var client *pgxpool.Pool
	for attempts >= 0 {
		if client, err = pgxpool.NewWithConfig(cxt, config); err != nil {
			attempts--
			time.Sleep(time.Second)
			continue
		}

		if err = client.Ping(cxt); err != nil {
			attempts--
			time.Sleep(time.Second)
			continue
		}
		break
	}

	return &Db{
		client: client,
	}, nil
}
