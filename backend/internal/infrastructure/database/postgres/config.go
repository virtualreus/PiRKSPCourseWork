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

	if driver == "" || dsn == "" {
		return nil, fmt.Errorf("PG_DRIVER and PG_DSN environment variables are required")
	}

	return &PGConfig{driver: driver, dsn: dsn}, nil
}
