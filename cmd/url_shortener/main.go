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
	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("Starting server")

	mem, err := storage.NewStorage(cfg)
	if err != nil {
		log.Error("failed to create storage", "error", err)
		os.Exit(1)
	}

	r := gin.New()

	server := &http.Server{
		Addr:         ":8080",
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

	err = server.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
