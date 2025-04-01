package storage

import (
	"context"
	"fmt"
	"url-shortener/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// type DBInterface interface {
// 	Create(value interface{}) *gorm.DB
// 	Where(query interface{}, args ...interface{}) *gorm.DB
// 	First(out interface{}, where ...interface{}) *gorm.DB
// }

type PostgreStorage struct {
	dbPool *pgxpool.Pool
	cache  *Cache
}

func InitDbPool(info *config.DbInfo) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", info.User, info.Password, info.Host, info.Port, info.Name)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pool config: %s", err)
	}
	config.MaxConns = info.MaxConn

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %s", err)
	}

	return pool, nil
}

// func DbConnect(info *config.DbInfo) (*gorm.DB, error) {

// 	dsn := fmt.Sprintf("host=%s user=%s port=%s dbname=%s sslmode=disable password='%s'", info.Host, info.User, info.Port, info.Name, info.Password)
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, ErrDBConnection
// 	}

// 	db.AutoMigrate(&UrlData{})

// 	return db, nil
// }

func NewPostgreStorage(pool *pgxpool.Pool, cache *Cache) *PostgreStorage {
	return &PostgreStorage{
		dbPool: pool,
		cache:  cache,
	}
}

func (p *PostgreStorage) CheckUrl(url string) (string, bool) {
	if val := p.cache.GetAlias(url); val != "" {
		return val, true
	}

	return "", false
}

func (p *PostgreStorage) SaveUrl(url, alias string) error {
	// newUrl := UrlData{
	// 	Url:   url,
	// 	Alias: alias,
	// }

	// проверка на уникальность alias
	if p.cache.GetUrl(alias) != "" {
		return ErrAlreadyExists
	}

	_, err := p.dbPool.Exec(context.Background(), "INSERT INTO urls (url, alias) VALUES ($1, $2)", url, alias)
	if err != nil {
		return fmt.Errorf("failed to save URL: %v", err)
	}
	// result := p.db.Create(&newUrl)
	// if result.Error != nil {
	// 	// доп проверка на уникальность
	// 	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
	// 		return ErrAlreadyExists
	// 	}
	// 	return fmt.Errorf("failed to save URL: %v", result.Error)
	// }

	// запись в кэш
	p.cache.Set(url, alias)

	return nil
}

func (p *PostgreStorage) GetUrl(alias string) (string, error) {
	// проверка в кэше
	if val := p.cache.GetUrl(alias); val != "" {
		return val, nil
	}

	var url string

	err := p.dbPool.QueryRow(context.Background(), "SELECT url FROM urls WHERE alias = $1", alias).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("url not found")
	}

	// проверка в бд
	// urlData := UrlData{}
	// result := p.db.Where("alias = ?", alias).First(&urlData)
	// if result.Error != nil {
	// 	return "", fmt.Errorf("not found")
	// }

	return url, nil
}

func (p *PostgreStorage) Close() {
	p.dbPool.Close()
}

func CheckTable(pool *pgxpool.Pool) error{
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM pg_catalog.pg_tables
            WHERE tablename = $1
        )
    `

	var exists bool
	err := pool.QueryRow(context.Background(), query, "urls").Scan(&exists)
    if err != nil {
        return fmt.Errorf("failed to check table existence: %s", err)
    }

	if !exists {
        createQuery := `
            CREATE TABLE urls (
                url TEXT PRIMARY KEY,
                alias TEXT NOT NULL
            )
        `
        _, err := pool.Exec(context.Background(), createQuery)
        if err != nil {
            return fmt.Errorf("failed to create table: %s", err)
        }
    }

	return nil
}