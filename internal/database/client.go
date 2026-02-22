package database

import (
	"context"
	"net/url"
	"time"
)

type Database interface {
	DeleteOlderThan(ctx context.Context, table string, cutoff time.Time) error
	Close() error
}

type Config struct {
	URL      string
	Username string
	Password string
	Schema   string
}

func New(cfg Config) (Database, error) {
	return newDatabase(cfg)
}

func (c Config) DSN() (string, error) {
	u, err := url.Parse(c.URL)
	if err != nil {
		return "", err
	}

	u.User = url.UserPassword(c.Username, c.Password)

	//if c.Schema != "" {
	//	q := u.Query()
	//	q.Set("schema", c.Schema)
	//	u.RawQuery = q.Encode()
	//}

	return u.String(), nil
}
