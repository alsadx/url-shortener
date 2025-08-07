package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"url-shortener/internal/storage"

	"url-shortener/internal/config"
)

type Storage struct {
	dbPool *pgxpool.Pool
}

func New(storage *config.Storage) (*Storage, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		storage.User, storage.Password, storage.Host, storage.Port, storage.Database,
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pool config: %s", err)
	}
	cfg.MaxConns = storage.MaxConn

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %s", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{dbPool: pool}, nil
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave, alias string) (int64, error) {
	op := "storage.postgres.SaveURL"

	var urlID int64

	query := `INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id`

	err := s.dbPool.QueryRow(ctx, query, urlToSave, alias).Scan(&urlID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return urlID, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	op := "storage.postgres.GetURL"

	var url string

	query := `SELECT url FROM url WHERE alias = $1`

	err := s.dbPool.QueryRow(ctx, query, alias).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	op := "storage.postgres.DeleteURL"

	query := `DELETE FROM url WHERE alias = $1`

	commandTag, err := s.dbPool.Exec(ctx, query, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}
