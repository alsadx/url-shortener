package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"

	"url-shortener/internal/config"
)

type Storage struct {
	dbPool *pgxpool.Pool
}

func New(storage *config.Storage) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", storage.Host, storage.User, storage.Password, storage.Database)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pool config: %s", err)
	}
	cfg.MaxConns = storage.MaxConn

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %s", err)
	}

	return &Storage{dbPool: pool}, nil
}
