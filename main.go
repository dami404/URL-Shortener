package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {

	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: logger: slog

	// TODO: storage: mongodb

	// TODO: router: chi, chi render

	// TODO: server:
}
