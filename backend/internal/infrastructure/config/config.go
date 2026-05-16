package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	AuthSecret  string
}

func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/autopark?sslmode=disable"),
		AuthSecret:  env("AUTH_SECRET", "dev-secret"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
