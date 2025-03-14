package main

import (
	"log/slog"
	"net/http"
	"os"

	"url-shortener/internal/config"
	"url-shortener/internal/router"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {

	// logger
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Debug("Debug info")
	log.Info("Starting server")

	// config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("", "error", err)
		os.Exit(1)
	}

	// init storage
	var mem storage.Storage
	switch cfg.Storage {
	case "inmemory":
		mem = storage.NewMemoryStorage()

	case "postgres":
		db, err := storage.DbConnect(&cfg.DB)
		if err != nil {
			log.Error("", "error", err)
			os.Exit(1)
		}
		cache := storage.GetCache(db)
		mem = storage.NewPostgreStorage(db, cache)

	default:
		log.Error("unknown storage type", "type", cfg.Storage)
		os.Exit(1)
	}

	// router
	r := gin.New()

	server := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	r.Use(router.MiddlewareLogging(log))

	r.POST("/save", router.ValidateUrl(), func(c *gin.Context) {
		router.PostHandler(c, mem, c.GetString("url"))
	})

	r.GET("/tinyurl/:alias", router.ValidateAlias(), func(c *gin.Context) {
		router.GetHandler(c, mem, c.GetString("alias"))
	})

	// err = r.Run(cfg.HttpServer.Address)
	err = server.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
