package config

import (
	"time"
)

type Config struct {
	DataCatalogURL            string
	DataCatalogToken          string
	DataCatalogDatabaseSchema string
	DatabaseURL               string
	DatabaseSchema            string
	HTTPTimeout               time.Duration
}

func DefaultConfig() Config {
	return Config{
		HTTPTimeout: 10 * time.Second,
	}
}
