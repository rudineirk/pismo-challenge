package config

import "github.com/caarlos0/env/v10"

type Config struct {
	IsProduction bool
	GoEnv        string `env:"GO_ENV"       envDefault:"development"`
	HTTPPort     int    `env:"HTTP_PORT"    envDefault:"3000"`
	DatabaseURL  string `env:"DATABASE_URL" envDefault:"postgresql://dev:development@127.0.0.1/pismo_challenge?sslmode=disable"` //nolint:lll // default value
	LogLevel     string `env:"LOG_LEVEL"    envDefault:"info"`
	LogFormat    string `env:"LOG_FORMAT"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return &Config{}, err
	}

	if cfg.GoEnv == "production" {
		cfg.IsProduction = true
	}

	return &cfg, nil
}
