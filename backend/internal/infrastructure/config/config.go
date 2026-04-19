package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string
	Port             string
	MongoURI         string
	MongoDBName      string
	JWTSecret        string
	JWTIssuer        string
	JWTExpiresInHour int
	AllowedOrigins   []string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	expires, err := strconv.Atoi(getEnv("JWT_EXPIRES_IN_HOURS", "24"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_EXPIRES_IN_HOURS: %w", err)
	}

	allowedOriginsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:5174,http://127.0.0.1:5173,http://127.0.0.1:5174")
	allowedOrigins := []string{}
	for _, origin := range strings.Split(allowedOriginsStr, ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			allowedOrigins = append(allowedOrigins, trimmed)
		}
	}

	cfg := Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		Port:             getEnv("PORT", "8080"),
		MongoURI:         getEnv("MONGO_URI", ""),
		MongoDBName:      getEnv("MONGO_DB_NAME", "inventory_management"),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTIssuer:        getEnv("JWT_ISSUER", "inventory-api"),
		JWTExpiresInHour: expires,
		AllowedOrigins:   allowedOrigins,
	}

	if cfg.MongoURI == "" {
		return Config{}, fmt.Errorf("MONGO_URI is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
