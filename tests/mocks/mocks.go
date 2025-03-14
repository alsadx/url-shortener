package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockCache struct {
    mock.Mock
}

func (m *MockCache) Set(url, alias string) {
    m.Called(url, alias)
}

func (m *MockCache) GetUrl(alias string) string {
    args := m.Called(alias)
    return args.String(0)
}

func (m *MockCache) GetAlias(url string) string {
    args := m.Called(url)
    return args.String(0)
}

type MockDB struct {
    mock.Mock
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
    args := m.Called(value)
    return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
    m.Called(query, args)
    return &gorm.DB{}
}

func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
    args := m.Called(out, where)
    return args.Get(0).(*gorm.DB)
}
