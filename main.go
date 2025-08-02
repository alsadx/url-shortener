package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("%+v\n", cfg)

	// TODO: init logger: slog

	// TODO: init storage: postgresql

	// TODO: init router: chi

	// TODO: run server
}
