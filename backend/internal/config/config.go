package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL       string
	ServerPort        string
	JWTSecret         string
	CookieSecure      bool
	CORSAllowedOrigin string
	LogFormat         string // "text" (default, local dev) or "json" (log aggregators)
	OTPTTL            time.Duration
	// OTPDevLogCodes prints generated OTP codes to the server log so login is
	// testable without a real SMS provider — never set this in production (see
	// internal/sms.ConsoleSender).
	OTPDevLogCodes bool
	// MSG91AuthKey/MSG91TemplateID configure the real SMS sender (internal/sms.MSG91Sender).
	// Takes priority over OTPDevLogCodes when both happen to be set.
	MSG91AuthKey    string
	MSG91TemplateID string
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

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "text"
	}

	otpTTL := 5 * time.Minute
	if raw := os.Getenv("OTP_TTL_MINUTES"); raw != "" {
		minutes, err := strconv.Atoi(raw)
		if err != nil || minutes <= 0 {
			return nil, errors.New("OTP_TTL_MINUTES must be a positive integer")
		}
		otpTTL = time.Duration(minutes) * time.Minute
	}

	return &Config{
		DatabaseURL:       databaseURL,
		ServerPort:        port,
		JWTSecret:         jwtSecret,
		CookieSecure:      os.Getenv("COOKIE_SECURE") == "true",
		CORSAllowedOrigin: corsAllowedOrigin,
		LogFormat:         logFormat,
		OTPTTL:            otpTTL,
		OTPDevLogCodes:    os.Getenv("OTP_DEV_LOG_CODES") == "true",
		MSG91AuthKey:      os.Getenv("MSG91_AUTH_KEY"),
		MSG91TemplateID:   os.Getenv("MSG91_TEMPLATE_ID"),
	}, nil
}
