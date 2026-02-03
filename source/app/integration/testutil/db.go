package testutil

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

// RunMigrations executes the database schema
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// TruncateTables clears all data from test tables for scenario isolation.
// Uses TRUNCATE with CASCADE to handle foreign key constraints.
func TruncateTables(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE person_attributes, person_images, request_log, person, key_value RESTART IDENTITY CASCADE
	`)
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}
	return nil
}

// CreatePerson inserts a test person and returns the UUID
func CreatePerson(ctx context.Context, pool *pgxpool.Pool, name, clientID string) (string, error) {
	var id string
	err := pool.QueryRow(ctx, `
		INSERT INTO person (client_id)
		VALUES ($1)
		RETURNING id::text
	`, clientID).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create person: %w", err)
	}
	return id, nil
}

// GetRawAttributeValue returns the raw encrypted value from the database
func GetRawAttributeValue(ctx context.Context, pool *pgxpool.Pool, personID, key string) ([]byte, error) {
	var rawValue []byte
	err := pool.QueryRow(ctx, `
		SELECT encrypted_value
		FROM person_attributes
		WHERE person_id = $1::uuid AND attribute_key = $2
	`, personID, key).Scan(&rawValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw attribute value: %w", err)
	}
	return rawValue, nil
}

// GetAttributeKeyVersion returns the key_version for an attribute
func GetAttributeKeyVersion(ctx context.Context, pool *pgxpool.Pool, personID, key string) (int32, error) {
	var keyVersion int32
	err := pool.QueryRow(ctx, `
		SELECT key_version
		FROM person_attributes
		WHERE person_id = $1::uuid AND attribute_key = $2
	`, personID, key).Scan(&keyVersion)
	if err != nil {
		return 0, fmt.Errorf("failed to get key version: %w", err)
	}
	return keyVersion, nil
}

// CountAttributes returns the number of attributes for a person
func CountAttributes(ctx context.Context, pool *pgxpool.Pool, personID string) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM person_attributes
		WHERE person_id = $1::uuid
	`, personID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count attributes: %w", err)
	}
	return count, nil
}

// InsertKeyValueDirect inserts a key-value pair directly into the database
func InsertKeyValueDirect(ctx context.Context, pool *pgxpool.Pool, key, value string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO key_value (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = CURRENT_TIMESTAMP
	`, key, value)
	if err != nil {
		return fmt.Errorf("failed to insert key-value: %w", err)
	}
	return nil
}

// GetKeyValueDirect retrieves a key-value pair directly from the database
func GetKeyValueDirect(ctx context.Context, pool *pgxpool.Pool, key string) (string, error) {
	var value string
	err := pool.QueryRow(ctx, `
		SELECT value FROM key_value WHERE key = $1
	`, key).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("failed to get key-value: %w", err)
	}
	return value, nil
}

// GetRequestLog returns an audit log entry by trace ID
func GetRequestLog(ctx context.Context, pool *pgxpool.Pool, traceID, encKey string) (caller, reason string, err error) {
	err = pool.QueryRow(ctx, `
		SELECT caller, reason
		FROM request_log
		WHERE trace_id = $1
	`, traceID).Scan(&caller, &reason)
	if err != nil {
		return "", "", fmt.Errorf("failed to get request log: %w", err)
	}
	return caller, reason, nil
}

// CountRequestLogs returns the number of request log entries for a trace ID
func CountRequestLogs(ctx context.Context, pool *pgxpool.Pool, traceID string) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM request_log
		WHERE trace_id = $1
	`, traceID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count request logs: %w", err)
	}
	return count, nil
}
