package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	errs "person-service/errors"
	health "person-service/healthcheck"
	key_value "person-service/key_value"
	"person-service/middleware"
	person_attributes "person-service/person_attributes"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	db "person-service/internal/db/generated"
)

// ============================================================================
// MAIN - Application Entry Point
// ============================================================================

func setupDb(port string) *db.Queries {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatalf("ERROR: DATABASE_URL environment variable is not set (error_code: %s)\n", errs.ErrDatabaseURLNotSet)
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("ERROR: Invalid PORT value: %s (error_code: %s)\n", err, errs.ErrInvalidPort)
	}

	fmt.Fprintf(os.Stdout, "INFO: Connecting to database...\n")

	// Initialize database connection with SQLC queries using pgxpool
	fmt.Fprintf(os.Stdout, "DEBUG: Opening database connection...\n")

	// Configure connection pool
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("ERROR: Failed to parse database URL: %v (error_code: %s)\n", err, errs.ErrFailedParseDBURL)
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
		log.Fatalf("ERROR: Failed to create connection pool: %v (error_code: %s)\n", err, errs.ErrFailedCreateConnPool)
	}

	fmt.Fprintf(os.Stdout, "DEBUG: Pinging database...\n")
	// Ping database
	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Database ping failed: %v (error_code: %s)\n", err, errs.ErrDatabasePingFailed)
		log.Fatalf("ERROR: Failed to ping database: %v (error_code: %s)\n", err, errs.ErrDatabasePingFailed)
	}

	fmt.Fprintf(os.Stdout, "DEBUG: Database ping successful!\n")

	queries := db.New(pool)
	return queries
}

func main() {
	// Ensure logs are written immediately (unbuffered)
	log.SetFlags(log.LstdFlags)

	_, err := fmt.Fprintf(os.Stdout, "INFO: Application starting...\n")
	if err != nil {
		log.Fatalf("ERROR: Failed to start application: %v\n", err)
	}

	// Load configuration from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	queries := setupDb(port)

	log.Println("INFO: Database connection successful")
	fmt.Fprintf(os.Stdout, "INFO: Database connection successful\n")

	// Create and setup Echo server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	healthHandler := health.NewHealthCheckHandler(queries)
	keyValueHandler := key_value.NewKeyValueHandler(queries)
	personAttributesHandler := person_attributes.NewPersonAttributesHandler(queries)

	// Setup routes
	e.GET("/health", healthHandler.Check)

	// Key-value API routes
	e.POST("/api/key-value", keyValueHandler.SetValue)
	e.GET("/api/key-value/:key", keyValueHandler.GetValue)
	e.DELETE("/api/key-value/:key", keyValueHandler.DeleteValue)

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

	// Log before starting server
	log.Printf("INFO: Server starting on port %s\n", port)
	fmt.Fprintf(os.Stdout, "INFO: Server starting on port %s\n", port)

	// Start server in goroutine
	go func() {
		if err := e.Start(e.Server.Addr); err != nil && err != http.ErrServerClosed {
			log.Printf("ERROR: Server error: %v (error_code: %s)\n", err, errs.ErrFailedStartServer)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	log.Printf("INFO: Server ready and listening on port %s\n", port)
	fmt.Fprintf(os.Stdout, "INFO: Server ready on port %s\n", port)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Fprintf(os.Stdout, "INFO: Shutting down server...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("ERROR: Server shutdown failed: %v (error_code: %s)\n", err, errs.ErrFailedShutdownServer)
	}
	fmt.Fprintf(os.Stdout, "INFO: Server gracefully stopped\n")
}
