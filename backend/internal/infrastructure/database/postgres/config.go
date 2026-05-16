package postgres

import (
	"fmt"
	"net/url"
	"os"
	"strings"
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

	dsn, err := ensurePostgresSSL(dsn)
	if err != nil {
		return nil, err
	}

	return &PGConfig{driver: driver, dsn: dsn}, nil
}

// ensurePostgresSSL adds sslmode=require for hosted Postgres (e.g. Railway proxy URLs).
func ensurePostgresSSL(dsn string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn, nil
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return dsn, nil
	}

	host := strings.ToLower(u.Hostname())
	if host == "" || host == "localhost" || host == "127.0.0.1" || strings.HasSuffix(host, ".internal") {
		return dsn, nil
	}

	q := u.Query()
	if q.Get("sslmode") != "" {
		return dsn, nil
	}

	q.Set("sslmode", "require")
	u.RawQuery = q.Encode()
	return u.String(), nil
}
