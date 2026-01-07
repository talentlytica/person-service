-- name: HealthCheck :exec
-- Perform a simple health check query  
SELECT 1;

-- name: GetValue :one
-- Retrieve a value by key
SELECT value FROM key_value WHERE key = sqlc.arg(key) LIMIT 1;

-- name: GetKeyValue :one
-- Retrieve the full key-value record by key
SELECT key, value, created_at, updated_at FROM key_value WHERE key = sqlc.arg(key) LIMIT 1;

-- name: SetValue :exec
-- Set a value by key
INSERT INTO key_value (key, value) VALUES (sqlc.arg(key), sqlc.arg(value))
ON CONFLICT (key) DO UPDATE SET value = sqlc.arg(value);

-- name: DeleteValue :exec
-- Delete a value by key
DELETE FROM key_value WHERE key = sqlc.arg(key);

-- ============================================================================
-- REQUEST LOG OPERATIONS
-- ============================================================================

-- name: InsertRequestLog :one
-- Insert a new request log entry with encrypted data
INSERT INTO request_log (
    trace_id, 
    caller, 
    reason, 
    encrypted_request_body, 
    encrypted_response_body, 
    key_version
) VALUES (
    sqlc.arg(trace_id), 
    sqlc.arg(caller), 
    sqlc.arg(reason), 
    pgp_sym_encrypt(sqlc.arg(encrypted_request_body), sqlc.arg(enc_key)), 
    pgp_sym_encrypt(sqlc.arg(encrypted_response_body), sqlc.arg(enc_key)), 
    sqlc.arg(key_version)
) RETURNING id, trace_id, created_at;

-- name: GetRequestLogByTraceId :one
-- Retrieve request log by trace_id with decrypted data
SELECT 
    id,
    trace_id,
    caller,
    reason,
    pgp_sym_decrypt(encrypted_request_body, sqlc.arg(enc_key)) AS request_body,
    pgp_sym_decrypt(encrypted_response_body, sqlc.arg(enc_key)) AS response_body,
    key_version,
    created_at
FROM request_log
WHERE trace_id = sqlc.arg(trace_id)
LIMIT 1;

-- name: CheckTraceIdExists :one
-- Check if a trace_id already exists (for idempotency)
SELECT EXISTS(SELECT 1 FROM request_log WHERE trace_id = sqlc.arg(trace_id));

-- ============================================================================
-- PERSON OPERATIONS
-- ============================================================================

-- name: CreatePerson :one
-- Create a new person
INSERT INTO person (client_id)
VALUES (sqlc.arg(client_id))
RETURNING id, client_id, created_at, updated_at, deleted_at;

-- name: GetPersonById :one
-- Get person by internal UUID
SELECT id, client_id, created_at, updated_at, deleted_at
FROM person
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
LIMIT 1;

-- name: GetPersonByClientId :one
-- Get person by client_id
SELECT id, client_id, created_at, updated_at, deleted_at
FROM person
WHERE client_id = sqlc.arg(client_id) AND deleted_at IS NULL
LIMIT 1;

-- name: UpdatePersonClientId :exec
-- Update person's client_id
UPDATE person
SET client_id = sqlc.arg(new_client_id), updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: SoftDeletePerson :exec
-- Soft delete a person
UPDATE person
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: HardDeletePerson :exec
-- Hard delete a person (use with caution)
DELETE FROM person WHERE id = sqlc.arg(id);

-- name: RestorePerson :exec
-- Restore a soft-deleted person
UPDATE person
SET deleted_at = NULL, updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id);

-- name: ListPersons :many
-- List all active persons with pagination
SELECT id, client_id, created_at, updated_at, deleted_at
FROM person
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- ============================================================================
-- PERSON ATTRIBUTES OPERATIONS
-- ============================================================================

-- name: CreateOrUpdatePersonAttribute :one
-- Create or update a person attribute with encryption
INSERT INTO person_attributes (
    person_id, 
    attribute_key, 
    encrypted_value, 
    key_version
) VALUES (
    sqlc.arg(person_id), 
    sqlc.arg(attribute_key), 
    pgp_sym_encrypt(sqlc.arg(attribute_value), sqlc.arg(enc_key)), 
    sqlc.arg(key_version)
)
ON CONFLICT (person_id, attribute_key) 
DO UPDATE SET 
    encrypted_value = pgp_sym_encrypt(sqlc.arg(attribute_value), sqlc.arg(enc_key)),
    key_version = sqlc.arg(key_version),
    updated_at = CURRENT_TIMESTAMP
RETURNING id, person_id, attribute_key, key_version, created_at, updated_at;

-- name: GetPersonAttribute :one
-- Get a single decrypted attribute for a person
SELECT 
    id,
    person_id,
    attribute_key,
    pgp_sym_decrypt(encrypted_value, sqlc.arg(enc_key)) AS attribute_value,
    key_version,
    created_at,
    updated_at
FROM person_attributes
WHERE person_id = sqlc.arg(person_id) AND attribute_key = sqlc.arg(attribute_key)
LIMIT 1;

-- name: GetAllPersonAttributes :many
-- Get all decrypted attributes for a person
SELECT 
    id,
    person_id,
    attribute_key,
    pgp_sym_decrypt(encrypted_value, sqlc.arg(enc_key)) AS attribute_value,
    key_version,
    created_at,
    updated_at
FROM person_attributes
WHERE person_id = sqlc.arg(person_id)
ORDER BY attribute_key;

-- name: GetMultiplePersonAttributes :many
-- Get multiple specific attributes for a person (pass array of keys)
SELECT 
    id,
    person_id,
    attribute_key,
    pgp_sym_decrypt(encrypted_value, sqlc.arg(enc_key)) AS attribute_value,
    key_version,
    created_at,
    updated_at
FROM person_attributes
WHERE person_id = sqlc.arg(person_id) AND attribute_key = ANY(sqlc.arg(attribute_keys)::citext[])
ORDER BY attribute_key;

-- name: DeletePersonAttribute :exec
-- Delete a specific attribute for a person
DELETE FROM person_attributes
WHERE person_id = sqlc.arg(person_id) AND attribute_key = sqlc.arg(attribute_key);

-- name: DeleteAllPersonAttributes :exec
-- Delete all attributes for a person
DELETE FROM person_attributes
WHERE person_id = sqlc.arg(person_id);

-- name: ListAttributeKeys :many
-- List all unique attribute keys used across all persons
SELECT DISTINCT attribute_key
FROM person_attributes
ORDER BY attribute_key;

-- name: CountPersonAttributes :one
-- Count attributes for a person
SELECT COUNT(*) FROM person_attributes WHERE person_id = sqlc.arg(person_id);

-- ============================================================================
-- PERSON IMAGES OPERATIONS
-- ============================================================================

-- name: CreateOrUpdatePersonImage :one
-- Create or update a person image with encryption
INSERT INTO person_images (
    person_id,
    attribute_key,
    image_type,
    encrypted_image_data,
    key_version,
    mime_type,
    file_size,
    width,
    height
) VALUES (
    sqlc.arg(person_id), 
    sqlc.arg(attribute_key), 
    sqlc.arg(image_type), 
    pgp_sym_encrypt(sqlc.arg(image_data), sqlc.arg(enc_key)), 
    sqlc.arg(key_version), 
    sqlc.arg(mime_type), 
    sqlc.arg(file_size), 
    sqlc.arg(width), 
    sqlc.arg(height)
)
ON CONFLICT (person_id, attribute_key)
DO UPDATE SET
    image_type = sqlc.arg(image_type),
    encrypted_image_data = pgp_sym_encrypt(sqlc.arg(image_data), sqlc.arg(enc_key)),
    key_version = sqlc.arg(key_version),
    mime_type = sqlc.arg(mime_type),
    file_size = sqlc.arg(file_size),
    width = sqlc.arg(width),
    height = sqlc.arg(height),
    updated_at = CURRENT_TIMESTAMP
RETURNING id, person_id, attribute_key, image_type, key_version, mime_type, file_size, width, height, created_at, updated_at;

-- name: GetPersonImage :one
-- Get a specific decrypted image for a person
SELECT 
    id,
    person_id,
    attribute_key,
    image_type,
    pgp_sym_decrypt(encrypted_image_data, sqlc.arg(enc_key)) AS image_data,
    key_version,
    mime_type,
    file_size,
    width,
    height,
    created_at,
    updated_at
FROM person_images
WHERE person_id = sqlc.arg(person_id) AND attribute_key = sqlc.arg(attribute_key)
LIMIT 1;

-- name: GetPersonImageMetadata :one
-- Get image metadata without decrypting the image data (for performance)
SELECT 
    id,
    person_id,
    attribute_key,
    image_type,
    key_version,
    mime_type,
    file_size,
    width,
    height,
    created_at,
    updated_at
FROM person_images
WHERE person_id = sqlc.arg(person_id) AND attribute_key = sqlc.arg(attribute_key)
LIMIT 1;

-- name: ListPersonImages :many
-- List all image metadata for a person (without decrypting)
SELECT 
    id,
    person_id,
    attribute_key,
    image_type,
    key_version,
    mime_type,
    file_size,
    width,
    height,
    created_at,
    updated_at
FROM person_images
WHERE person_id = sqlc.arg(person_id)
ORDER BY created_at DESC;

-- name: ListPersonImagesByType :many
-- List images of a specific type for a person (without decrypting)
SELECT 
    id,
    person_id,
    attribute_key,
    image_type,
    key_version,
    mime_type,
    file_size,
    width,
    height,
    created_at,
    updated_at
FROM person_images
WHERE person_id = sqlc.arg(person_id) AND image_type = sqlc.arg(image_type)
ORDER BY created_at DESC;

-- name: DeletePersonImage :exec
-- Delete a specific image for a person
DELETE FROM person_images
WHERE person_id = sqlc.arg(person_id) AND attribute_key = sqlc.arg(attribute_key);

-- name: DeleteAllPersonImages :exec
-- Delete all images for a person
DELETE FROM person_images
WHERE person_id = sqlc.arg(person_id);

-- name: CountPersonImages :one
-- Count images for a person
SELECT COUNT(*) FROM person_images WHERE person_id = sqlc.arg(person_id);

-- ============================================================================
-- COMBINED OPERATIONS
-- ============================================================================

-- name: GetPersonWithAttributes :one
-- Get person basic info (to be combined with attributes in application layer)
SELECT 
    p.id,
    p.client_id,
    p.created_at,
    p.updated_at,
    p.deleted_at
FROM person p
WHERE p.id = sqlc.arg(id) AND p.deleted_at IS NULL
LIMIT 1;

-- name: SearchPersonsByAttribute :many
-- Search persons by a specific decrypted attribute value (note: performance intensive)
SELECT DISTINCT
    p.id,
    p.client_id,
    p.created_at,
    p.updated_at
FROM person p
JOIN person_attributes pa ON p.id = pa.person_id
WHERE pa.attribute_key = sqlc.arg(attribute_key)
    AND pgp_sym_decrypt(pa.encrypted_value, sqlc.arg(enc_key)) = sqlc.arg(attribute_value)
    AND p.deleted_at IS NULL;

-- name: BulkCreatePersonAttributes :copyfrom
-- Bulk insert person attributes (use with COPY FROM)
INSERT INTO person_attributes (
    person_id,
    attribute_key,
    encrypted_value,
    key_version
) VALUES (
    sqlc.arg(person_id), 
    sqlc.arg(attribute_key), 
    sqlc.arg(encrypted_value), 
    sqlc.arg(key_version)
);

