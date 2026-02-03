package testutil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pgContainer *postgres.PostgresContainer
	pool        *pgxpool.Pool
	once        sync.Once
	initErr     error
)

const (
	dbName     = "testdb"
	dbUser     = "testuser"
	dbPassword = "testpass"
)

// GetPool returns a singleton PostgreSQL connection pool.
// It starts the container on first call and reuses it for all tests.
func GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	once.Do(func() {
		initErr = startContainer(ctx)
	})

	if initErr != nil {
		return nil, initErr
	}
	return pool, nil
}

func startContainer(ctx context.Context) error {
	var err error

	pgContainer, err = postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to get connection string: %w", err)
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 5 * time.Minute

	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// GetConnectionString returns the database connection string
func GetConnectionString(ctx context.Context) (string, error) {
	if pgContainer == nil {
		return "", fmt.Errorf("container not started")
	}
	return pgContainer.ConnectionString(ctx, "sslmode=disable")
}

// CleanupContainer terminates the PostgreSQL container.
// Should be called in AfterSuite.
func CleanupContainer(ctx context.Context) error {
	if pool != nil {
		pool.Close()
	}
	if pgContainer != nil {
		return pgContainer.Terminate(ctx)
	}
	return nil
}
