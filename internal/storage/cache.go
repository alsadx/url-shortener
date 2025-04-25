package storage

import (
	"sync"

	"gorm.io/gorm"
)

type Cache struct {
	UrlAlias sync.Map
	AliasUrl sync.Map
	once      sync.Once
	db        *gorm.DB
}

var (
	cacheInstance *Cache
	cacheMutex sync.Mutex
)

func GetCache(db *gorm.DB) *Cache {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if cacheInstance == nil {
		cacheInstance = &Cache{db: db}
		cacheInstance.loadFromDB()
	}
	return cacheInstance
}

func (c *Cache) loadFromDB() {
	urls := []UrlData{}
	c.db.Find(&urls)

	for _, url := range urls {
		c.AliasUrl.Store(url.Alias, url.Url)
		c.UrlAlias.Store(url.Url, url.Alias)
	}
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
