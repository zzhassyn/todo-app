package tasks_postgres_repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	
	core_pgx_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool/pgx"
)

func setupTestDB(t *testing.T) *TasksRepository {
	ctx := context.Background()

	// Start Postgres container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("test-user"),
		postgres.WithPassword("test-password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Run migrations
	migrationsPath, err := filepath.Abs("../../../../../migrations")
	require.NoError(t, err)
	migrationsPath = filepath.ToSlash(migrationsPath)
	m, err := migrate.New("file://"+migrationsPath, connStr)
	require.NoError(t, err)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations: %v", err)
	}

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Initialize pool
	cfg := core_pgx_pool.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "test-user",
		Password: "test-password",
		Database: "test-db",
		Timeout:  5 * time.Second,
	}

	pool, err := core_pgx_pool.NewPool(ctx, cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return NewTasksRepository(pool)
}
