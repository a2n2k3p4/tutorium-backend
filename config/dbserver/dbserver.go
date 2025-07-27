package dbserver

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func (c *Config) DBUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func ConnectDB(cfg *Config) (*pgxpool.Pool, error) {
	dbUrl := cfg.DBUrl()

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
