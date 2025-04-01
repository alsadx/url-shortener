package storage

import (
	"fmt"
	"sync"
)

type MemoryStorage struct {
	aliasUrl map[string]string
	urlAlias map[string]string
	mux  sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		aliasUrl: make(map[string]string),
		urlAlias: make(map[string]string),
	}
}

func (m *MemoryStorage) CheckUrl(url string) (string, bool) {
	val, ok := m.urlAlias[url]
	
	return val, ok
}

func (m *MemoryStorage) SaveUrl(url, alias string) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, exists := m.aliasUrl[alias]; exists {
		return ErrAlreadyExists
	}

	m.aliasUrl[alias] = url
	m.urlAlias[url] = alias

	return nil
}

func (m *MemoryStorage) GetUrl(alias string) (string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if val, ok := m.aliasUrl[alias]; ok {
		return val, nil
	}

	return "", fmt.Errorf("not found")
}

func (m *MemoryStorage) Close() {}