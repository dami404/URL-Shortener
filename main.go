package main

import (
	"github.com/common-nighthawk/go-figure"
	"log/slog"
	"os"
	"url-shortener/internal/config"
)

func main() {

	figure.NewFigure("URL-Shortener", "", true).Print()

	cfg := config.MustLoad()

	// ! сделать свой логгер
	log := setupLogger(cfg.Env)
	log.Info("Starting App", slog.String("env", cfg.Env))
	log.Debug("Debug level is enable")

	// TODO: storage: mongodb

	// TODO: router: chi, chi render

	// TODO: server:
}

func setupLogger(env string) (log *slog.Logger) {
	envMap := map[string]*slog.Logger{
		"local":       slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		"development": slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		"production":  slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}
	log = envMap[env]
	return
}
