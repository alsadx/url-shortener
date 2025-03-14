package storage_test

import (
    "testing"
    "url-shortener/internal/storage"
    "github.com/stretchr/testify/assert"
)

func TestMemoryStorage_SaveUrl(t *testing.T) {
    mem := storage.NewMemoryStorage()

    err := mem.SaveUrl("https://example.com", "abc123")
    assert.NoError(t, err)

    url, err := mem.GetUrl("abc123")
    assert.NoError(t, err)
    assert.Equal(t, "https://example.com", url)
}

func TestMemoryStorage_GetUrl_NotFound(t *testing.T) {
    mem := storage.NewMemoryStorage()

    _, err := mem.GetUrl("nonexist")
    assert.Error(t, err)
}

func TestMemoryStorage_CheckUrl(t *testing.T) {
    mem := storage.NewMemoryStorage()
    mem.SaveUrl("https://example.com", "abc123")

    alias, exists := mem.CheckUrl("https://example.com")
    assert.True(t, exists)
    assert.Equal(t, "abc123", alias)
}

func TestMemoryStorage_CheckInvalidUrl(t *testing.T) {
	mem := storage.NewMemoryStorage()

    alias, exists := mem.CheckUrl("nonexist.none")
    assert.False(t, exists)
    assert.Equal(t, "", alias)
}