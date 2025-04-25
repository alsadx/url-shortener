package storage

import (
	"errors"
	"fmt"
	"url-shortener/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBInterface interface {
    Create(value interface{}) *gorm.DB
    Where(query interface{}, args ...interface{}) *gorm.DB
    First(out interface{}, where ...interface{}) *gorm.DB
}

type PostgreStorage struct {
	db    DBInterface
	cache *Cache
}

type UrlData struct {
	ID    uint   `gorm:"primaryKey"`
	Url   string `gorm:"type:varchar(255);not null"`
	Alias string `gorm:"type:varchar(255);not null;unique"`
}

// TODO: закрытие коннекта к бд

func DbConnect(info *config.DbInfo) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s port=%s dbname=%s sslmode=disable password='%s'", info.Host, info.User, info.Port, info.Name, info.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, ErrDBConnection
	}

	db.AutoMigrate(&UrlData{})

	return db, nil
}

func NewPostgreStorage(db DBInterface, cache *Cache) *PostgreStorage {
	return &PostgreStorage{
		db:    db,
		cache: cache,
	}
}

func (p *PostgreStorage) CheckUrl(url string) (string, bool) {
	if val := p.cache.GetAlias(url); val != "" {
		return val, true
	}

	return "", false
}

func (p *PostgreStorage) SaveUrl(url, alias string) error {
	newUrl := UrlData{
		Url:   url,
		Alias: alias,
	}

	// проверка на уникальность alias
	if p.cache.GetUrl(alias) != "" {
		return ErrAlreadyExists
	}

	result := p.db.Create(&newUrl)
	if result.Error != nil {
		// доп проверка на уникальность
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to save URL: %v", result.Error)
	}

	// запись в кэш
	p.cache.Set(url, alias)

	return nil
}

func (p *PostgreStorage) GetUrl(alias string) (string, error) {
	// проверка в кэше
	if val := p.cache.GetUrl(alias); val != "" {
		return val, nil
	}

	// проверка в бд
	urlData := UrlData{}
	result := p.db.Where("alias = ?", alias).First(&urlData)
	if result.Error != nil {
		return "", fmt.Errorf("not found")
	}


	return urlData.Url, nil
}
