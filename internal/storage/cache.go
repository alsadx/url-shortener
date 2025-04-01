package storage

import (
	"fmt"
	"sync"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cache struct {
	UrlAlias sync.Map
	AliasUrl sync.Map
	once     sync.Once
	dbPool   *pgxpool.Pool
}

var (
	cacheInstance *Cache
	cacheMutex    sync.Mutex
)

func GetCache(pool *pgxpool.Pool) *Cache {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if cacheInstance == nil {
		cacheInstance = &Cache{dbPool: pool}
		cacheInstance.loadFromDB()
	}
	return cacheInstance
}

func (c *Cache) loadFromDB() error {
	// urls := []UrlData{}
	// c.db.Find(&urls)

	// for _, url := range urls {
	// 	c.AliasUrl.Store(url.Alias, url.Url)
	// 	c.UrlAlias.Store(url.Url, url.Alias)
	// }

	rows, err := c.dbPool.Query(context.Background(), "SELECT alias, url FROM urls")
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
        var shortURL, longURL string
        if err := rows.Scan(&shortURL, &longURL); err != nil {
            return fmt.Errorf("scan failed: %w", err)
        }
        c.AliasUrl.Store(shortURL, longURL)
		c.UrlAlias.Store(longURL, shortURL)
    }

    if err := rows.Err(); err != nil {
        return fmt.Errorf("rows iteration failed: %w", err)
    }

	return nil
}

func (c *Cache) Set(url, alias string) {
	c.UrlAlias.Store(url, alias)
	c.AliasUrl.Store(alias, url)
}

func (c *Cache) GetUrl(alias string) string {
	if value, ok := c.AliasUrl.Load(alias); ok {
		return value.(string)
	}

	return ""
}

func (c *Cache) GetAlias(url string) string {
	if value, ok := c.UrlAlias.Load(url); ok {
		return value.(string)
	}

	return ""
}
