package health

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	db "person-service/internal/db/generated"
	"person-service/internal/testdb"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	pool, err = testdb.GetPool(ctx)
	if err != nil {
		log.Fatalf("Failed to get pool: %v", err)
	}
	if err := testdb.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	os.Exit(m.Run())
}

func TestNewHealthCheckHandler(t *testing.T) {
	queries := db.New(pool)
	handler := NewHealthCheckHandler(queries)
	assert.NotNil(t, handler)
	assert.Equal(t, queries, handler.queries)
}

func TestCheck_Success(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewHealthCheckHandler(queries)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = handler.Check(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "healthy")
}

func TestCheck_DatabaseError(t *testing.T) {
	// Create a handler with a closed pool to simulate database error
	closedPool, err := createClosedPool()
	assert.NoError(t, err)

	queries := db.New(closedPool)
	handler := NewHealthCheckHandler(queries)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = handler.Check(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "HC_001_HEALTH_CHECK_FAILED")
}

// createClosedPool creates a pool and immediately closes it to simulate database errors
func createClosedPool() (*pgxpool.Pool, error) {
	ctx := context.Background()
	connStr, err := testdb.GetConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	closedPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Close the pool immediately
	closedPool.Close()

	return closedPool, nil
}
