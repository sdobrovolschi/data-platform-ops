package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
)

type sqlDB struct {
	db      *sql.DB
	dialect dialect
	schema  string
}

func newDatabase(cfg Config) (Database, error) {
	dialect, err := newDialect(cfg.URL)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(dialect.Driver(), cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &sqlDB{
		db:      db,
		dialect: dialect,
		schema:  cfg.Schema,
	}, nil
}

func (d *sqlDB) DeleteOlderThan(ctx context.Context, tableName string, cutoff time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := d.deleteQuery(tableName)
	if _, err := d.db.ExecContext(ctx, query, cutoff); err != nil {
		return fmt.Errorf("unable to execute query: %w", err)
	}
	return nil
}

func (d *sqlDB) deleteQuery(table string) string {
	return fmt.Sprintf(
		`DELETE FROM %s.%s WHERE created_at < %s`,
		d.dialect.Quote(d.schema),
		d.dialect.Quote(table),
		d.dialect.Placeholder(1),
	)
}

func (d *sqlDB) Close() error {
	return d.db.Close()
}
