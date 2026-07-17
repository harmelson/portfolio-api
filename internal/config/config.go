package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	Port           string
	GoogleClientID string
	DevAuthEnabled bool
	DevAuthToken   string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	cfg := Config{
		DatabaseURL:    GetEnv("DATABASE_URL", ""),
		Port:           GetEnv("PORT", "8080"),
		GoogleClientID: GetEnv("GOOGLE_CLIENT_ID", ""),
		DevAuthEnabled: GetEnv("DEV_AUTH_ENABLED", "false") == "true",
		DevAuthToken:   GetEnv("DEV_AUTH_TOKEN", "dev-token"),
	}

	configValidator(cfg)

	return cfg
}

func GetEnv(key string, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return val
}

func configValidator(cfg Config) {
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.Port == "" {
		log.Fatal("PORT is required")
	}

	if cfg.GoogleClientID == "" && !cfg.DevAuthEnabled {
		log.Fatal("GOOGLE_CLIENT_ID is required when dev auth is disabled")
	}
}
