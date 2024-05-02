package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"` // struct tags
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  struct {
		Address     string        `yaml:"address" env-default:"localhost:8082"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	} `yaml:"http_server"`
}

func MustLoad() *Config {

	if err := godotenv.Load(); err != nil && os.Getenv("CONFIG_PATH") == "" {
		log.Fatal("no .env file found")
	}
	configPath := os.Getenv("CONFIG_PATH")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// еще не инициализирован основной логгер, поэтому используем этот
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}
	return &cfg
}
