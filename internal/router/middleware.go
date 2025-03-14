package router

import (
	"log/slog"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
)

func MiddlewareLogging(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("Incoming request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_addr", c.ClientIP(),
		)
		c.Next()
	}
}

func ValidateUrl() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid JSON"})
			c.Abort()
			return
		}

		// проверка валидности URL
		_, err := url.ParseRequestURI(req.Url)
		if err != nil || req.Url == "" {
			c.JSON(400, gin.H{"error": "invalid URL"})
			c.Abort()
			return
		}

		c.Set("url", req.Url)
		c.Next()
	}
}

func ValidateAlias() gin.HandlerFunc {
	aliasRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{10}$`)
	return func(c *gin.Context) {
		alias := c.Param("alias")

		if !aliasRegex.MatchString(alias) {
			c.JSON(400, gin.H{"error": "invalid alias format"})
			c.Abort()
			return
		}

		c.Set("alias", alias) 
		c.Next()
	}
}
