package database

import (
	"context"
	"time"
)

type Database interface {
	DeleteOlderThan(ctx context.Context, table string, cutoff time.Time) error
	Close() error
}

type Config struct {
	URL    string
	Schema string
}

func New(cfg Config) (Database, error) {
	return newDatabase(cfg)
}
