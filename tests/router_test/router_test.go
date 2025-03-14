package router_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/router"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mem := storage.NewMemoryStorage()
	r.POST("/save", func(c *gin.Context) {
		router.PostHandler(c, mem, "")
	})

	reqBody := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/save", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mem := storage.NewMemoryStorage()
	mem.SaveUrl("https://example.com", "abc123")

	r.GET("/tinyurl/:alias", func(c *gin.Context) {
		router.GetHandler(c, mem, "abc123")
	})

	req := httptest.NewRequest(http.MethodGet, "/tinyurl/abc123", nil)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusMovedPermanently, resp.Code)
}
