package main

import (
	"github.com/common-nighthawk/go-figure"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

func main() {

	figure.NewFigure("URL-Shortener", "", true).Print()

	cfg := config.MustLoad()

	// TODO: сделать свой логгер
	log := setupLogger(cfg.Env)
	log.Info("Starting App", slog.String("env", cfg.Env))
	log.Debug("Debug level is enable")

	// TODO: заменить sqlite на mongodb
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to connect to storage", sl.Error(err))
		os.Exit(1) // можно заменить на return
	}

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
