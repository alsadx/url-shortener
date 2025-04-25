package storage

import (
	"errors"
	"fmt"
	"url-shortener/internal/config"
)

type Storage interface {
	SaveUrl(longUrl, shortUrl string) error
	GetUrl(shortUrl string) (string, error)
	CheckUrl(longUrl string) (string, bool)
}

var (
	ErrAlreadyExists = errors.New("URL already exists")
	ErrDBConnection  = errors.New("database connection error")
)

func NewStorage(cfg *config.Config) (Storage, error) {
	var mem Storage
	switch cfg.Storage {
	case "inmemory":
		mem = NewMemoryStorage()

	case "postgres":
		db, err := DbConnect(&cfg.DB)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %s", err)
		}
		cache := GetCache(db)
		mem = NewPostgreStorage(db, cache)

	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage)
	}

	return mem, nil
}
