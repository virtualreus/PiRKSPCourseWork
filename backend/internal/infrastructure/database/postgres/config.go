package postgres

import (
	"fmt"
	"os"
)

type Config interface {
	GetDSN() string
	GetDriver() string
}

type PGConfig struct {
	driver string
	dsn    string
}

func (p PGConfig) GetDSN() string {
	return p.dsn
}

func (p PGConfig) GetDriver() string {
	return p.driver
}

func NewConfig() (Config, error) {
	driver := os.Getenv("PG_DRIVER")
	dsn := os.Getenv("PG_DSN")

	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if driver == "" && dsn != "" {
		driver = "postgres"
	}

	if driver == "" || dsn == "" {
		return nil, fmt.Errorf("set PG_DSN or DATABASE_URL (and optionally PG_DRIVER)")
	}

	return &PGConfig{driver: driver, dsn: dsn}, nil
}
