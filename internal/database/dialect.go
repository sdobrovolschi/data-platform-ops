package database

import (
	"fmt"
	"net/url"
)

const (
	Postgres  string = "postgres"
	SQLServer string = "sqlserver"
)

type dialect interface {
	Driver() string
	Quote(identifier string) string
	Placeholder(n int) string
}

func newDialect(rawURL string) (dialect, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid database URL: %w", err)
	}

	switch u.Scheme {
	case Postgres:
		return postgresDialect{}, nil
	case SQLServer:
		return sqlServerDialect{}, nil
	default:
		return nil, fmt.Errorf("unsupported database: %s", u.Scheme)
	}
}

type postgresDialect struct{}

func (postgresDialect) Driver() string {
	return Postgres
}

func (postgresDialect) Quote(id string) string {
	return fmt.Sprintf(`"%s"`, id)
}

func (postgresDialect) Placeholder(n int) string {
	return fmt.Sprintf("$%d", n)
}

type sqlServerDialect struct{}

func (sqlServerDialect) Driver() string {
	return SQLServer
}

func (sqlServerDialect) Quote(id string) string {
	return fmt.Sprintf(`[%s]`, id)
}

func (sqlServerDialect) Placeholder(n int) string {
	return fmt.Sprintf("@p%d", n)
}
