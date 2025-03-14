package storage

import "errors"

type Storage interface {
	SaveUrl(longUrl, shortUrl string) error
	GetUrl(shortUrl string) (string, error)
	CheckUrl(longUrl string) (string, bool)
}

var (
	ErrAlreadyExists = errors.New("URL already exists")
	ErrDBConnection  = errors.New("database connection error")
)
