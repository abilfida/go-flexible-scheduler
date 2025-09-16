package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN       string
	Port      string
	Timezone  string
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Abaikan jika .env tidak ada (misal di production, pakai env var asli)
		// return nil, err
	}

	dsn := os.Getenv("DB_DSN")
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Default port
	}

	// Atur timezone, default ke "Asia/Jakarta" jika tidak ada di .env
	tz := os.Getenv("APP_TIMEZONE")
	if tz == "" {
		tz = "Asia/Jakarta" // <-- DEFAULT TIMEZONE
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	return &Config{
		DSN:       dsn,
		Port:      port,
		Timezone:  tz,
		JWTSecret: jwtSecret,
	}, nil
}
