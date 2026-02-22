package config

import (
	"fmt"
	"os"
)

func Load() (*Config, error) {
	cfg := DefaultConfig()

	requiredEnv := map[string]*string{
		"DATABASE_URL":                 &cfg.DatabaseURL,
		"DATABASE_SCHEMA":              &cfg.DatabaseSchema,
		"DATA_CATALOG_URL":             &cfg.DataCatalogURL,
		"DATA_CATALOG_TOKEN":           &cfg.DataCatalogToken,
		"DATA_CATALOG_DATABASE_SCHEMA": &cfg.DataCatalogDatabaseSchema,
	}

	for key, dest := range requiredEnv {
		value := os.Getenv(key)
		if value == "" {
			return nil, fmt.Errorf("%s is required", key)
		}
		*dest = value
	}

	return &cfg, nil
}
