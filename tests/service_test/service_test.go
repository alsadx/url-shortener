package service_test

import (
	"regexp"
	"testing"
	"url-shortener/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestShortener(t *testing.T) {
	longUrl := "https://example.com"
	alias := service.Shortener(longUrl)

	assert.NotEmpty(t, alias)

	assert.Len(t, alias, 10)

	aliasRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{10}$`)
	assert.True(t, aliasRegex.MatchString(alias))
}
