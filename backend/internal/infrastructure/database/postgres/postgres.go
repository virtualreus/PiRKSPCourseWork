package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db      *sqlx.DB
	Builder squirrel.StatementBuilderType
}

func NewPostgres(cfg Config) (*Postgres, error) {
	db, err := sqlx.Open(cfg.GetDriver(), cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Postgres{
		db:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func (p *Postgres) SqlDB() *sql.DB {
	return p.db.DB
}

func (p *Postgres) SqlxDB() *sqlx.DB {
	return p.db
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
