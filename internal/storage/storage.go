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
	Close()
}

var (
	ErrAlreadyExists = errors.New("URL already exists")
)

func InitStorage(cfg *config.Config) (Storage, error) {
	var mem Storage
	switch cfg.Storage {
	case "inmemory":
		mem = NewMemoryStorage()

	case "postgres":
		pool, err := InitDbPool(&cfg.DB)
		if err != nil {
			return nil, err
		}
		// defer pool.Close()
		cache := GetCache(pool)
		mem = NewPostgreStorage(pool, cache)
		err = CheckTable(pool)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown storage type %s", cfg.Storage)
	}

	return mem, nil
}
