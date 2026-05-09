package postgres

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/hackathon/*.sql
var migrations embed.FS

func MigrateDB(db *Postgres) error {
	goose.SetBaseFS(migrations)

	if err := goose.Up(db.SqlDB(), "migrations/hackathon"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
