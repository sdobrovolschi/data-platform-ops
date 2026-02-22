package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type httpClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func (c *httpClient) GetTables(ctx context.Context, databaseSchema string) ([]Table, error) {
	req, err := c.newGetTablesRequest(ctx, databaseSchema)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch tables: %s", string(body))
	}

	return c.parseTables(resp)
}

func (c *httpClient) newGetTablesRequest(ctx context.Context, databaseSchema string) (*http.Request, error) {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse base URL: %w", err)
	}

	endpoint := base.JoinPath("/v1/tables")
	query := endpoint.Query()
	query.Set("databaseSchema", databaseSchema)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	return req, nil
}

func (c *httpClient) parseTables(resp *http.Response) ([]Table, error) {
	var respPayload struct {
		Data []struct {
			Name            string
			RetentionPeriod *string
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	tables := make([]Table, len(respPayload.Data))
	for i, t := range respPayload.Data {
		var dur *time.Duration
		if t.RetentionPeriod != nil {
			parsed, err := time.ParseDuration(*t.RetentionPeriod)
			if err != nil {
				return nil, fmt.Errorf("unable to parse retention period for table %s: %w", t.Name, err)
			}
			dur = &parsed
		}
		tables[i] = Table{
			Name:            t.Name,
			RetentionPeriod: dur,
		}
	}
	return tables, nil
}

func (c *httpClient) Close() {
	c.client.CloseIdleConnections()
}
