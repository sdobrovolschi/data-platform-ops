package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"data-platform-ops/internal/config"
	"data-platform-ops/internal/database"
	"data-platform-ops/internal/metadata"
)

type Cleaner struct {
	metadata metadata.Client
	database database.Database
	schema   string
	log      *slog.Logger
}

func NewCleaner(log *slog.Logger) (Cleaner, error) {
	cfg, err := config.Load()
	if err != nil {
		return Cleaner{}, err
	}

	meta := metadata.New(metadata.Config{
		BaseURL: cfg.DataCatalogURL,
		Token:   cfg.DataCatalogToken,
		Timeout: cfg.HTTPTimeout,
	})

	db, err := database.New(database.Config{
		URL:    cfg.DatabaseURL,
		Schema: cfg.DatabaseSchema,
	})
	if err != nil {
		return Cleaner{}, err
	}

	return Cleaner{
		metadata: meta,
		database: db,
		schema:   cfg.DataCatalogDatabaseSchema,
		log:      log,
	}, nil
}

func (c *Cleaner) run(ctx context.Context) error {
	tables, err := c.metadata.GetTables(ctx, c.schema)
	if err != nil {
		return fmt.Errorf("fetch tables: %w", err)
	}

	for _, table := range tables {
		if table.RetentionPeriod != nil {
			cutoff := time.Now().UTC().Add(-*table.RetentionPeriod)
			if err := c.database.DeleteOlderThan(ctx, table.Name, cutoff); err == nil {
				c.log.InfoContext(ctx, fmt.Sprintf("deleted rows from %s older than %s", table.Name, cutoff))
			} else {
				c.log.ErrorContext(ctx, fmt.Sprintf("delete failed for table %s", table.Name), "err", err)
			}
		}
	}

	return nil
}

func (c *Cleaner) Close() error {
	c.metadata.Close()
	return c.database.Close()
}
