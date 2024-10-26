package env

import (
	"errors"
	"os"
)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func (p *pgConfig) DSN() string {
	return p.dsn
}

// NewPGConfig creates a new PostgreSQL configuration based on environment variables.
func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}
	return &pgConfig{
		dsn: dsn,
	}, nil
}
