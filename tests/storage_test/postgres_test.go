package storage_test

import (
    "errors"
    "testing"
    "url-shortener/tests/mocks"
    "url-shortener/internal/storage"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"
)

func TestPostgreStorage_SaveUrl(t *testing.T) {
    mockDB := new(mocks.MockDB)

    cache := &storage.Cache{}

    ps := storage.NewPostgreStorage(mockDB, cache)

    mockDB.On("Create", mock.Anything).Return(&gorm.DB{})

    err := ps.SaveUrl("https://example.com", "abc123")

    assert.NoError(t, err)

    mockDB.AssertCalled(t, "Create", mock.Anything)
}

func TestPostgreStorage_SaveUrl_Error(t *testing.T) {
    mockDB := new(mocks.MockDB)

    cache := &storage.Cache{}

    ps := storage.NewPostgreStorage(mockDB, cache)

    mockDB.On("Create", mock.Anything).Return(&gorm.DB{Error: errors.New("database error")})

    err := ps.SaveUrl("https://example.com", "abc123")

    assert.Error(t, err)
    assert.Equal(t, "failed to save URL: database error", err.Error())

    mockDB.AssertCalled(t, "Create", mock.Anything)
}

func TestPostgreStorage_GetUrl(t *testing.T) {
    mockDB := new(mocks.MockDB)

    cache := &storage.Cache{}
    cache.Set("https://example.com", "abc123")

    ps := storage.NewPostgreStorage(mockDB, cache)

    var result storage.UrlData
    result.Url = "https://example.com"

    mockDB.On("Where", "alias = ?", []interface{}{"abc123"}).Return(mockDB)
    mockDB.On("First", &storage.UrlData{}, []interface{}{mock.Anything}).Run(func(args mock.Arguments) {
        arg := args.Get(0).(*storage.UrlData)
        *arg = result
    }).Return(&gorm.DB{})

    url, err := ps.GetUrl("abc123")

    assert.NoError(t, err)
    assert.Equal(t, "https://example.com", url)

}

func TestPostgreStorage_CheckUrl(t *testing.T) {
    mockDB := new(mocks.MockDB)

    cache := &storage.Cache{}
    cache.Set("https://example.com", "abc123")
    
    ps := storage.NewPostgreStorage(mockDB, cache)

    var result storage.UrlData
    result.Alias = "abc123"

    mockDB.On("Where", "url = ?", []interface{}{"https://example.com"}).Return(mockDB)
    mockDB.On("First", &storage.UrlData{}, []interface{}{mock.Anything}).Run(func(args mock.Arguments) {
        arg := args.Get(0).(*storage.UrlData)
        *arg = result
    }).Return(&gorm.DB{})

    alias, exists := ps.CheckUrl("https://example.com")

    assert.True(t, exists)
    assert.Equal(t, "abc123", alias)
}