package db_cleaner

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mssql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCleanerRemovesExpiredRowsFromPostgreSQL(t *testing.T) {
	ctx := context.Background()
	net := setupNetwork(t, ctx)

	pg := setupPostgres(t, ctx, net, "fixtures/init-postgresql.sql")
	setupWiremock(t, ctx, net, "fixtures/wiremock")
	setupCleaner(t, ctx, net,
		map[string]string{
			"DATA_CATALOG_URL":             "http://wiremock.testcontainer.docker:8080",
			"DATA_CATALOG_TOKEN":           "eyJraWQiOiJHYj",
			"DATA_CATALOG_DATABASE_SCHEMA": "service.database.schema",
			"DATABASE_URL":                 "postgres://postgres:postgres@postgresql.testcontainer.docker:5432/database?sslmode=disable",
			"DATABASE_SCHEMA":              "schema",
		})

	connStr, err := pg.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)
	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)

	assertMaxID(t, db, `"schema".table_a`, 3)
	assertMaxID(t, db, `"schema".table_b`, 1)
	assertMaxID(t, db, `"schema".table_c`, 1)
}

//go:embed fixtures/init-mssql.sql
var sqlServerInitScript []byte

func TestCleanerRemovesExpiredRowsFromSQLServer(t *testing.T) {
	ctx := context.Background()
	net := setupNetwork(t, ctx)

	sqlServer := setupSQLServer(t, ctx, net)
	setupWiremock(t, ctx, net, "fixtures/wiremock")
	setupCleaner(t, ctx, net,
		map[string]string{
			"DATA_CATALOG_URL":             "http://wiremock.testcontainer.docker:8080",
			"DATA_CATALOG_TOKEN":           "eyJraWQiOiJHYj",
			"DATA_CATALOG_DATABASE_SCHEMA": "service.database.schema",
			"DATABASE_URL":                 "sqlserver://sa:SuperStrong@Passw0rd@mssql.testcontainer.docker:1433?database=database&encrypt=disable",
			"DATABASE_SCHEMA":              "schema",
		})

	connStr, err := sqlServer.ConnectionString(ctx, "database=database", "encrypt=disable")
	assert.NoError(t, err)
	db, err := sql.Open("sqlserver", connStr)
	assert.NoError(t, err)

	assertMaxID(t, db, "[schema].table_a", 3)
	assertMaxID(t, db, "[schema].table_b", 1)
	assertMaxID(t, db, "[schema].table_c", 1)
}

func setupNetwork(t *testing.T, ctx context.Context) *testcontainers.DockerNetwork {
	t.Helper()
	net, err := network.New(ctx)
	assert.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)
	return net
}

func setupPostgres(t *testing.T, ctx context.Context, net *testcontainers.DockerNetwork, initScript string) *postgres.PostgresContainer {
	t.Helper()
	pg, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithInitScripts(filepath.Join(initScript)),
		postgres.WithDatabase("database"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
		network.WithNetwork([]string{"postgresql.testcontainer.docker"}, net),
	)
	assert.NoError(t, err)
	testcontainers.CleanupContainer(t, pg)
	return pg
}

func setupSQLServer(t *testing.T, ctx context.Context, net *testcontainers.DockerNetwork) *mssql.MSSQLServerContainer {
	t.Helper()
	sqlc, err := mssql.Run(ctx,
		"mcr.microsoft.com/mssql/server:2022-CU14-ubuntu-22.04",
		mssql.WithAcceptEULA(),
		mssql.WithInitSQL(bytes.NewReader(sqlServerInitScript)),
		mssql.WithPassword("SuperStrong@Passw0rd"),
		network.WithNetwork([]string{"mssql.testcontainer.docker"}, net),
	)
	assert.NoError(t, err)
	testcontainers.CleanupContainer(t, sqlc)
	return sqlc
}

func setupWiremock(t *testing.T, ctx context.Context, net *testcontainers.DockerNetwork, mappings string) testcontainers.Container {
	t.Helper()
	wm, err := testcontainers.Run(ctx,
		"wiremock/wiremock:3.13.2",
		testcontainers.WithExposedPorts("8080/tcp"),
		testcontainers.WithFiles(testcontainers.ContainerFile{
			HostFilePath:      filepath.Join(mappings),
			ContainerFilePath: "/home/wiremock",
			FileMode:          0o755,
		}),
		testcontainers.WithWaitStrategy(wait.ForHTTP("/__admin").
			WithPort("8080/tcp").
			WithStartupTimeout(30*time.Second),
		),
		network.WithNetwork([]string{"wiremock.testcontainer.docker"}, net),
	)
	assert.NoError(t, err)
	testcontainers.CleanupContainer(t, wm)
	return wm
}

func setupCleaner(t *testing.T, ctx context.Context, net *testcontainers.DockerNetwork, env map[string]string) testcontainers.Container {
	t.Helper()
	cleaner, err := testcontainers.Run(ctx, "",
		testcontainers.WithDockerfile(testcontainers.FromDockerfile{
			Context:    "../..",
			Dockerfile: "cmd/db-cleaner/Dockerfile",
			KeepImage:  false,
		}),
		testcontainers.WithEnv(env),
		testcontainers.WithWaitStrategy(wait.ForExit()),
		network.WithNetwork(nil, net),
	)
	assert.NoError(t, err)
	testcontainers.CleanupContainer(t, cleaner)
	return cleaner
}

func assertMaxID(t *testing.T, db *sql.DB, table string, expected int) {
	t.Helper()
	var maxID int
	err := db.QueryRow(fmt.Sprintf("SELECT MAX(id) FROM %s", table)).Scan(&maxID)
	assert.NoError(t, err)
	assert.Equal(t, expected, maxID)
}
