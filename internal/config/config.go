package config

import (
	"log"
	"os"
)

type Config struct {
	AppPort string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	JWTPublicKey string

	APIBaseURL string
}

func Load() *Config {
	cfg := &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPass:       getEnv("DB_PASS", "postgres"),
		DBName:       getEnv("DB_NAME", "ticketing_booking"),
		JWTPublicKey: getEnv("JWT_PUBLIC_KEY", ""),
		APIBaseURL:   getEnv("API_BASE_URL", ""),
	}

	if cfg.JWTPublicKey == "" {
		log.Fatal("JWT_PUBLIC_KEY is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
