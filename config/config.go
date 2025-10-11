package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN                  string
	Port                 string
	Timezone             string
	JWTSecret            string
	CleanupIntervalHours int
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

	// Muat interval cleanup, default ke 24 jam jika tidak ada atau error
	cleanupInterval, err := strconv.Atoi(os.Getenv("CLEANUP_INTERVAL_HOURS"))
	if err != nil {
		cleanupInterval = 24 // Default 24 jam
	}

	return &Config{
		DSN:                  dsn,
		Port:                 port,
		Timezone:             tz,
		JWTSecret:            jwtSecret,
		CleanupIntervalHours: cleanupInterval,
	}, nil
}
