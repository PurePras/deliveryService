package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL       string
	ServerPort        string
	JWTSecret         string
	CookieSecure      bool
	CORSAllowedOrigin string
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsAllowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if corsAllowedOrigin == "" {
		corsAllowedOrigin = "http://localhost:5173"
	}

	return &Config{
		DatabaseURL:       databaseURL,
		ServerPort:        port,
		JWTSecret:         jwtSecret,
		CookieSecure:      os.Getenv("COOKIE_SECURE") == "true",
		CORSAllowedOrigin: corsAllowedOrigin,
	}, nil
}
