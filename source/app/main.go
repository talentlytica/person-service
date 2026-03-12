package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	errs "person-service/errors"
	health "person-service/healthcheck"
	dbpkg "person-service/internal/db"
	db "person-service/internal/db/generated"
	key_value "person-service/key_value"
	"person-service/logging"
	"person-service/middleware"
	person "person-service/person"
	person_attributes "person-service/person_attributes"
)

// Version is set at build time via ldflags
var Version = "dev"

// ============================================================================
// MAIN - Application Entry Point
// ============================================================================

func setupDb(port string) (*db.Queries, *pgxpool.Pool) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logging.Error("DATABASE_URL environment variable is not set",
			"error_code", errs.ErrDatabaseURLNotSet)
		os.Exit(1)
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		logging.Error("Invalid PORT value",
			"error", err,
			"error_code", errs.ErrInvalidPort)
		os.Exit(1)
	}

	logging.Info("Connecting to database")

	// Configure connection pool
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		logging.Error("Failed to parse database URL",
			"error", err,
			"error_code", errs.ErrFailedParseDBURL)
		os.Exit(1)
	}

	// Set connection pool parameters
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 1 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	// Create connection pool with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logging.Error("Failed to create connection pool",
			"error", err,
			"error_code", errs.ErrFailedCreateConnPool)
		os.Exit(1)
	}

	logging.Debug("Pinging database")

	// Ping database
	if err := pool.Ping(ctx); err != nil {
		logging.Error("Database ping failed",
			"error", err,
			"error_code", errs.ErrDatabasePingFailed)
		os.Exit(1)
	}

	logging.Debug("Database ping successful")

	// Run migrations before creating queries
	// Convert postgres:// to pgx5:// for golang-migrate compatibility
	migrateURL := databaseURL
	if len(migrateURL) >= 11 && migrateURL[:11] == "postgres://" {
		migrateURL = "pgx5://" + migrateURL[11:]
	} else if len(migrateURL) >= 13 && migrateURL[:13] == "postgresql://" {
		migrateURL = "pgx5://" + migrateURL[13:]
	}

	if err := dbpkg.RunMigrations(ctx, pool, migrateURL); err != nil {
		logging.Error("Database migration failed",
			"error", err)
		os.Exit(1)
	}

	queries := db.New(pool)
	return queries, pool
}

func main() {
	// Initialize structured logging
	logging.Init()

	logging.Info("Application starting")

	// Load configuration from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	queries, _ := setupDb(port)

	logging.Info("Database connection successful")

	// Create and setup Echo server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Apply trace middleware globally (must be first to capture all requests)
	e.Use(middleware.TraceMiddleware())

	healthHandler := health.NewHealthCheckHandler(queries)
	keyValueHandler := key_value.NewKeyValueHandler(queries)
	personAttributesHandler := person_attributes.NewPersonAttributesHandler(queries)

	// Setup routes
	e.GET("/health", healthHandler.Check)
	e.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"version": Version})
	})

	// Key-value API routes
	e.POST("/api/key-value", keyValueHandler.SetValue)
	e.GET("/api/key-value/:key", keyValueHandler.GetValue)
	e.DELETE("/api/key-value/:key", keyValueHandler.DeleteValue)

	// Person CRUD API routes - protected with Bearer token middleware
	personHandler := person.NewPersonHandler(queries)
	personGroup := e.Group("/api/person", middleware.BearerMiddleware())
	personGroup.POST("", personHandler.CreatePerson)
	personGroup.GET("/:id", personHandler.GetPerson)
	personGroup.PATCH("/:id", personHandler.UpdatePerson)
	personGroup.DELETE("/:id", personHandler.DeletePerson)

	// Person attributes API routes - protected with API key middleware
	personAttributesGroup := e.Group("/persons", middleware.APIKeyMiddleware())
	personAttributesGroup.POST("/:personId/attributes", personAttributesHandler.CreateAttribute)
	personAttributesGroup.PUT("/:personId/attributes", personAttributesHandler.CreateAttribute)
	personAttributesGroup.GET("/:personId/attributes", personAttributesHandler.GetAllAttributes)
	personAttributesGroup.GET("/:personId/attributes/:attributeId", personAttributesHandler.GetAttribute)
	personAttributesGroup.PUT("/:personId/attributes/:attributeId", personAttributesHandler.UpdateAttribute)
	personAttributesGroup.DELETE("/:personId/attributes/:attributeId", personAttributesHandler.DeleteAttribute)

	// Configure server
	e.Server = &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logging.Info("Server starting", "port", port)

	// Start server in goroutine
	go func() {
		if err := e.Start(e.Server.Addr); err != nil && err != http.ErrServerClosed {
			logging.Error("Server error",
				"error", err,
				"error_code", errs.ErrFailedStartServer)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	logging.Info("Server ready and listening", "port", port)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logging.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logging.Error("Server shutdown failed",
			"error", err,
			"error_code", errs.ErrFailedShutdownServer)
		os.Exit(1)
	}
	logging.Info("Server gracefully stopped")
}
