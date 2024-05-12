package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host              string
	Port              int
	User              string
	Password          string
	DBName            string
	SSLMode           string
	ConnectionTimeout time.Duration
}

func Connect(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := getConnectionString(cfg)
	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return pool, err
}

func getConnectionString(cfg Config) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.Password, cfg.SSLMode)
}
