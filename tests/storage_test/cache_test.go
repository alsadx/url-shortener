package storage_test

import (
    "testing"
    // "url-shortener/internal/storage"
	"url-shortener/tests/mocks"
    "github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
    mockCache := new(mocks.MockCache)

    mockCache.On("Set", "https://example.com", "abc123").Return()
    mockCache.On("GetUrl", "abc123").Return("https://example.com")
    mockCache.On("GetAlias", "https://example.com").Return("abc123")

    mockCache.Set("https://example.com", "abc123")

    url := mockCache.GetUrl("abc123")
    assert.Equal(t, "https://example.com", url)

    alias := mockCache.GetAlias("https://example.com")
    assert.Equal(t, "abc123", alias)

    mockCache.AssertCalled(t, "Set", "https://example.com", "abc123")
    mockCache.AssertCalled(t, "GetUrl", "abc123")
    mockCache.AssertCalled(t, "GetAlias", "https://example.com")
}

func TestCache_GetUrl_NotFound(t *testing.T) {
    mockCache := new(mocks.MockCache)

    mockCache.On("GetUrl", "nonexistent_alias").Return("")

    url := mockCache.GetUrl("nonexistent_alias")
    assert.Equal(t, "", url)

    mockCache.AssertCalled(t, "GetUrl", "nonexistent_alias")
}

// func TestCache_SetAndGet(t *testing.T) {
//     cache := storage.GetCache(nil)

//     cache.Set("https://example1.com", "abc111")
// 	cache.Set("https://example2.com", "abc222")

// 	// поиск по alias
//     url := cache.GetUrl("abc222")
//     assert.Equal(t, "https://example2.com", url)

// 	// поиск по url
//     alias := cache.GetAlias("https://example1.com")
//     assert.Equal(t, "abc111", alias)

// 	// поиск по несуществующему url
// 	alias = cache.GetAlias("https://example3.com")
// 	assert.Equal(t, "", alias)

// 	// поиск по несуществующему alias
// 	url = cache.GetUrl("abc333")
// 	assert.Equal(t, "", url)
// }