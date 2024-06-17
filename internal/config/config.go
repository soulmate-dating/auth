package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type Postgres struct {
	Host              string        `env:"POSTGRES_HOST,required" example:"localhost"`
	Port              int           `env:"POSTGRES_PORT" envDefault:"5432"`
	User              string        `env:"POSTGRES_USER,required" example:"glimpse"`
	Password          string        `env:"POSTGRES_PASSWORD,required" example:"password"`
	Database          string        `env:"POSTGRES_DB,required" example:"glimpse"`
	SSLMode           string        `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	ConnectionTimeout time.Duration `env:"POSTGRES_CONNECTION_TIMEOUT" envDefault:"60s"`
}

type API struct {
	Network string `env:"API_NETWORK" envDefault:"tcp"`
	Address string `env:"API_ADDRESS,required" example:"localhost:8080"`
}

type JWT struct {
	SecretKey              string        `env:"JWT_ACCESS_SECRET,required"`
	RefreshSecretKey       string        `env:"JWT_REFRESH_SECRET,required"`
	Issuer                 string        `env:"JWT_ISSUER,required"`
	AccessExpirationHours  time.Duration `env:"JWT_ACCESS_EXPIRATION" envDefault:"24h"`
	RefreshExpirationHours time.Duration `env:"JWT_REFRESH_EXPIRATION" envDefault:"720h"`
}

type Metrics struct {
	Address string `env:"METRICS_ADDRESS,required" example:":8080"`
}

type Config struct {
	Postgres Postgres
	API      API
	JWT      JWT
	Metrics  Metrics
}

func Load() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}
	return cfg, nil
}
