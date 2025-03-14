package router

import (
	"fmt"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context, mem storage.Storage, alias string) {
	// alias := c.Param("alias")

	longUrl, err := mem.GetUrl(alias)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "not found",
		})
		return
	}

	c.Redirect(301, longUrl)
}
