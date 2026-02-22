package metadata

import (
	"context"
	"net/http"
	"time"
)

type Client interface {
	GetTables(ctx context.Context, databaseSchema string) ([]Table, error)
	Close()
}

type Config struct {
	BaseURL string
	Token   string
	Timeout time.Duration
}

func New(cfg Config) Client {
	return &httpClient{
		baseURL: cfg.BaseURL,
		token:   cfg.Token,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}
