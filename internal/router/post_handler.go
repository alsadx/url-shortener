package router

import (
	"log"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Url string `json:"url"`
}

func PostHandler(c *gin.Context, mem storage.Storage, url string) {
	// проверка url в памяти
	alias, exists := mem.CheckUrl(url)
	if exists {
		c.JSON(200, gin.H{
			"alias": alias,
		})
		return
	}

	alias = service.Shortener(url)
	// убрать
	log.Println(alias)

	err := mem.SaveUrl(url, alias)
	for err == storage.ErrAlreadyExists {
		log.Println(err)
		alias = service.Shortener(alias)
		err = mem.SaveUrl(url, alias)
	}

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "failed to save URL",
		})
		return
	}

	c.JSON(200, gin.H{
		"shortUrl": alias,
	})

}
