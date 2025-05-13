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

func Connect(ctx context.Context, conf *pgxpool.Config, attempts int) (*Db, error) {
	var (
		client  *pgxpool.Pool
		backoff = time.Second
		err     error
	)

	for attempts > 0 {
		client, err = pgxpool.NewWithConfig(ctx, conf)
		if err != nil {
			attempts--
			time.Sleep(backoff)
			continue
		}

		if err = client.Ping(ctx); err != nil {
			attempts--
			time.Sleep(backoff)
			backoff *= 2
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

func (db *Db) RunMigrate(url string) error {
	m, err := migrate.New("file://migrate", url)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
