package main

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/logger/sl"
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

	router := chi.NewRouter()
	router.Use(middleware.RequestID)

	// TODO: добавить в логи айпи адрес пользователя и чето с ним придумать интересное
	//router.Use(middleware.RealIP)

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	log.Info("Starting server", slog.String("address", cfg.HTTPServer.Address))
	// TODO: server
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
