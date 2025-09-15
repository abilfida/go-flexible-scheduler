package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN  string
	Port string
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

	return &Config{
		DSN:  dsn,
		Port: port,
	}, nil
}
