package pgstorage

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	client *pgxpool.Pool
}

func NewPgConf(url string) (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 15
	config.MinConns = 5
	config.MaxConnIdleTime = 30 * time.Minute

	return config, nil
}

func Connact(cxt context.Context, config *pgxpool.Config, attempts int) (*Db, error) {
	var (
		client  *pgxpool.Pool
		backoff = time.Second
		err     error
	)

	for attempts > 0 {
		if client, err = pgxpool.NewWithConfig(cxt, config); err != nil {
			attempts--
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if err = client.Ping(cxt); err != nil {
			time.Sleep(backoff)
			backoff *= 2
			attempts--
			continue
		}
		break
	}

	if attempts == 0 {
		return nil, err
	}

	return &Db{
		client: client,
	}, nil
}

func (db *Db) RunMigration(url string) error {
	m, err := migrate.New("file://migrate", url)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
