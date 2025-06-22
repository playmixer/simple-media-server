package config

import (
	"errors"
	"fmt"
	"os"
	"simple-media-server/internal/adapters/api/rest"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Rest      rest.Config
	BaseURL   string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	SecretKey string `env:"SECRET_KEY" envDefault:""`
	LogLevel  string `env:"LOG_LEVEL" envDefault:"debug"`
}

func Init() (*Config, error) {
	cfg := &Config{
		Rest: rest.Config{},
	}

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return cfg, fmt.Errorf("failed load enviorements from file: %w", err)
	}

	if err := env.Parse(cfg); err != nil {
		return cfg, fmt.Errorf("failed parse env: %w", err)
	}

	return cfg, nil
}
