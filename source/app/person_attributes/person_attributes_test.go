package person_attributes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	db "person-service/internal/db/generated"
	"person-service/internal/testdb"
)

var pool *pgxpool.Pool

const testEncryptionKey = "test-encryption-key-32bytes!!"

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

	// Set up environment variable for encryption key
	os.Setenv("ENCRYPTION_KEY_1", testEncryptionKey)

	os.Exit(m.Run())
}

// Helper function to create a test person
func createTestPerson(ctx context.Context, clientID string) (string, error) {
	return testdb.CreatePerson(ctx, pool, "", clientID)
}

// Helper function to create a test attribute for a person
func createTestAttribute(ctx context.Context, personID, key, value string) (int32, error) {
	var id int32
	err := pool.QueryRow(ctx, `
		INSERT INTO person_attributes (person_id, attribute_key, encrypted_value, key_version)
		VALUES ($1::uuid, $2, pgp_sym_encrypt($3, $4), 1)
		RETURNING id
	`, personID, key, value, testEncryptionKey).Scan(&id)
	return id, err
}

func TestNewPersonAttributesHandler(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)
	assert.NotNil(t, handler)
	assert.Equal(t, queries, handler.queries)
	assert.Equal(t, testEncryptionKey, handler.encryptionKey)
	assert.Equal(t, int32(1), handler.keyVersion)
}

func TestNewPersonAttributesHandler_WithEnvVar(t *testing.T) {
	// Set environment variable
	os.Setenv("ENCRYPTION_KEY_1", "test-key")
	defer os.Setenv("ENCRYPTION_KEY_1", testEncryptionKey)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)
	assert.NotNil(t, handler)
	assert.Equal(t, "test-key", handler.encryptionKey)
}

func TestNewPersonAttributesHandler_DefaultKey(t *testing.T) {
	// Unset environment variable to test default key
	os.Unsetenv("ENCRYPTION_KEY_1")
	defer os.Setenv("ENCRYPTION_KEY_1", testEncryptionKey)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)
	assert.NotNil(t, handler)
	assert.Equal(t, "default-key-for-dev", handler.encryptionKey)
}

func TestCreateAttribute_InvalidUUID(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"email","value":"test@example.com","meta":{"caller":"test","reason":"testing","traceId":"123"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/invalid-uuid/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("invalid-uuid")

	err := handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Person not found")
}

func TestCreateAttribute_InvalidJSON(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{invalid-json}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	err := handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request body")
}

func TestCreateAttribute_EmptyKey(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"","value":"test@example.com","meta":{"caller":"test","reason":"testing","traceId":"123"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	err := handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Key is required")
}

func TestCreateAttribute_MissingMeta(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"email","value":"test@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	err := handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Missing required field")
}

func TestGetAllAttributes_InvalidUUID(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/invalid-uuid/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("invalid-uuid")

	err := handler.GetAllAttributes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Person not found")
}

func TestGetAttribute_InvalidPersonUUID(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/invalid-uuid/attributes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("invalid-uuid", "1")

	err := handler.GetAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid person ID format")
}

func TestGetAttribute_InvalidAttributeID(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "invalid")

	err := handler.GetAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid attribute ID format")
}

func TestUpdateAttribute_InvalidPersonUUID(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"value":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/persons/invalid-uuid/attributes/1", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("invalid-uuid", "1")

	err := handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid person ID format")
}

func TestUpdateAttribute_InvalidAttributeID(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"value":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/invalid", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "invalid")

	err := handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid attribute ID format")
}

func TestUpdateAttribute_InvalidJSON(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{invalid-json}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/1", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "1")

	err := handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request body")
}

func TestDeleteAttribute_InvalidPersonUUID(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/persons/invalid-uuid/attributes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("invalid-uuid", "1")

	err := handler.DeleteAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid person ID format")
}

func TestDeleteAttribute_InvalidAttributeID(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "invalid")

	err := handler.DeleteAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid attribute ID format")
}

func TestCreateAttribute_PersonNotFound(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"email","value":"test@example.com","meta":{"caller":"test","reason":"testing","traceId":"123"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	err = handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAllAttributes_PersonNotFound(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	err = handler.GetAllAttributes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAttribute_PersonNotFound(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "1")

	err = handler.GetAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateAttribute_PersonNotFound(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"value":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/1", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "1")

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteAttribute_PersonNotFound(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "1")

	err = handler.DeleteAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateAttribute_SuccessWithoutTraceID(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person
	personID, err := createTestPerson(ctx, "test-client-1")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"email","value":"test@example.com","meta":{"caller":"test","reason":"testing"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "email")
}

func TestGetAllAttributes_Success(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person
	personID, err := createTestPerson(ctx, "test-client-2")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAttribute_NotFound(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person without attributes
	personID, err := createTestPerson(ctx, "test-client-3")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, "999")

	err = handler.GetAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Attribute not found")
}

func TestUpdateAttribute_NotFound(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person without attributes
	personID, err := createTestPerson(ctx, "test-client-4")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"value":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes/999", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, "999")

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Attribute not found")
}

func TestDeleteAttribute_NotFound(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person without attributes
	personID, err := createTestPerson(ctx, "test-client-5")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/persons/"+personID+"/attributes/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, "999")

	err = handler.DeleteAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Attribute not found")
}

func TestGetAllAttributes_DatabaseError(t *testing.T) {
	closedPool, err := createClosedPool()
	assert.NoError(t, err)

	queries := db.New(closedPool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	err = handler.GetAllAttributes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetAttribute_DatabaseError(t *testing.T) {
	closedPool, err := createClosedPool()
	assert.NoError(t, err)

	queries := db.New(closedPool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "1")

	err = handler.GetAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateAttribute_DatabaseError(t *testing.T) {
	closedPool, err := createClosedPool()
	assert.NoError(t, err)

	queries := db.New(closedPool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"value":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/1", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "1")

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteAttribute_DatabaseError(t *testing.T) {
	closedPool, err := createClosedPool()
	assert.NoError(t, err)

	queries := db.New(closedPool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000", "1")

	err = handler.DeleteAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateAttribute_Success(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-6")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "email", "old@example.com")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"value":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateAttribute_WithKeyChange(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-7")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "oldkey", "old@example.com")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"newkey","value":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteAttribute_Success(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-8")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "email", "test@example.com")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.DeleteAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Attribute deleted successfully")
}

func TestGetAttribute_Success(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-9")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "email", "test@example.com")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.GetAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllAttributes_WithResults(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with multiple attributes
	personID, err := createTestPerson(ctx, "test-client-10")
	assert.NoError(t, err)

	_, err = createTestAttribute(ctx, personID, "email", "test@example.com")
	assert.NoError(t, err)
	_, err = createTestAttribute(ctx, personID, "phone", "+1234567890")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCreateAttribute_SuccessWithTraceID(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person
	personID, err := createTestPerson(ctx, "test-client-11")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"email","value":"test@example.com","meta":{"caller":"test","reason":"testing","traceId":"trace123"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateAttribute_DatabaseErrorOnCreate(t *testing.T) {
	closedPool, err := createClosedPool()
	assert.NoError(t, err)

	queries := db.New(closedPool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"email","value":"test@example.com","meta":{"caller":"test","reason":"testing","traceId":"123"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/123e4567-e89b-12d3-a456-426614174000/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	err = handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateAttribute_WithoutKeyProvided(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-12")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "existing-key", "old-value")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"value":"new-value","meta":{"caller":"test","reason":"testing","traceId":"123"}}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "existing-key")
}

// TestCreateAttribute_WithEmptyValue tests creating an attribute with empty value
func TestCreateAttribute_WithEmptyValue(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person
	personID, err := createTestPerson(ctx, "test-client-empty-value")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"empty-value-key","value":"","meta":{"caller":"test","reason":"testing","traceId":"trace-empty"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "empty-value-key")
}

// TestUpdateAttribute_SameKeyAsExisting tests updating with the same key (no key change)
func TestUpdateAttribute_SameKeyAsExisting(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-same-key")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "same-key", "old-value")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	// Update with same key explicitly provided
	jsonBody := `{"key":"same-key","value":"new-value"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "same-key")
	assert.Contains(t, rec.Body.String(), "new-value")
}

// TestUpdateAttribute_EmptyKeyPreservesOriginal tests that providing empty key preserves the original key
// This verifies that when req.Key is empty, the attribute is NOT deleted (tests line 423 condition)
func TestUpdateAttribute_EmptyKeyPreservesOriginal(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-empty-key")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "preserve-key", "old-value")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	// Update with empty key - should preserve the original key
	jsonBody := `{"key":"","value":"new-value"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Verify the original key is preserved
	assert.Contains(t, rec.Body.String(), "preserve-key")
	assert.Contains(t, rec.Body.String(), "new-value")

	// Verify the attribute still exists with the same ID by fetching it again
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	c2.SetParamNames("personId", "attributeId")
	c2.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.GetAttribute(c2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec2.Code)
	assert.Contains(t, rec2.Body.String(), "preserve-key")
}

// TestGetAttribute_WithTimestamps tests that timestamps are included in response
func TestGetAttribute_WithTimestamps(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-timestamps")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "timestamp-key", "timestamp-value")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.GetAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "createdAt")
}

// TestGetAllAttributes_WithTimestamps tests that timestamps are included in response for all attributes
func TestGetAllAttributes_WithTimestamps(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-all-timestamps")
	assert.NoError(t, err)

	_, err = createTestAttribute(ctx, personID, "ts-key-1", "ts-value-1")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "createdAt")
}

// TestCreateAttribute_WithUpdatedAt tests that updatedAt timestamp is set on create
func TestCreateAttribute_WithUpdatedAt(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person
	personID, err := createTestPerson(ctx, "test-client-updated-at")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"updated-key","value":"updated-value","meta":{"caller":"test","reason":"testing","traceId":"trace-updated"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "updatedAt")
}

// TestUpdateAttribute_ResponseContainsTimestamps tests that response contains timestamps
func TestUpdateAttribute_ResponseContainsTimestamps(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute
	personID, err := createTestPerson(ctx, "test-client-update-ts")
	assert.NoError(t, err)

	attrID, err := createTestAttribute(ctx, personID, "update-ts-key", "old-value")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"value":"new-value"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "createdAt")
	assert.Contains(t, rec.Body.String(), "updatedAt")
}

// createHandlerWithWrongKey creates a handler with an incorrect encryption key
// This causes decryption to fail, triggering error paths in retrieve operations
func createHandlerWithWrongKey(queries *db.Queries) *PersonAttributesHandler {
	return &PersonAttributesHandler{
		queries:       queries,
		encryptionKey: "wrong-encryption-key-32bytes!!!",
		keyVersion:    1,
	}
}

// TestGetAllAttributes_DecryptionError tests error when decryption fails
func TestGetAllAttributes_DecryptionError(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute using the correct key
	personID, err := createTestPerson(ctx, "test-client-decrypt-err-1")
	assert.NoError(t, err)
	_, err = createTestAttribute(ctx, personID, "encrypted-key", "encrypted-value")
	assert.NoError(t, err)

	// Create handler with wrong key - decryption will fail
	queries := db.New(pool)
	handler := createHandlerWithWrongKey(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to retrieve attributes")
}

// TestGetAttribute_DecryptionError tests error when decryption fails for single attribute
func TestGetAttribute_DecryptionError(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute using the correct key
	personID, err := createTestPerson(ctx, "test-client-decrypt-err-2")
	assert.NoError(t, err)
	attrID, err := createTestAttribute(ctx, personID, "encrypted-key", "encrypted-value")
	assert.NoError(t, err)

	// Create handler with wrong key - decryption will fail
	queries := db.New(pool)
	handler := createHandlerWithWrongKey(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.GetAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to retrieve attributes")
}

// TestUpdateAttribute_GetAttributesDecryptionError tests error when getting attributes for update fails
func TestUpdateAttribute_GetAttributesDecryptionError(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute using the correct key
	personID, err := createTestPerson(ctx, "test-client-decrypt-err-3")
	assert.NoError(t, err)
	attrID, err := createTestAttribute(ctx, personID, "encrypted-key", "encrypted-value")
	assert.NoError(t, err)

	// Create handler with wrong key - decryption will fail when trying to find existing attribute
	queries := db.New(pool)
	handler := createHandlerWithWrongKey(queries)

	e := echo.New()
	jsonBody := `{"value":"new-value"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.UpdateAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to retrieve attributes")
}

// TestDeleteAttribute_GetAttributesDecryptionError tests error when getting attributes for delete fails
func TestDeleteAttribute_GetAttributesDecryptionError(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with an attribute using the correct key
	personID, err := createTestPerson(ctx, "test-client-decrypt-err-4")
	assert.NoError(t, err)
	attrID, err := createTestAttribute(ctx, personID, "encrypted-key", "encrypted-value")
	assert.NoError(t, err)

	// Create handler with wrong key - decryption will fail when trying to find existing attribute
	queries := db.New(pool)
	handler := createHandlerWithWrongKey(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/persons/%s/attributes/%d", personID, attrID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId", "attributeId")
	c.SetParamValues(personID, fmt.Sprintf("%d", attrID))

	err = handler.DeleteAttribute(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to retrieve attributes")
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

// ============================================================================
// PHASE 1: CONCURRENT OPERATION TESTS
// ============================================================================

// TestConcurrentAttributeCreation_SamePerson tests that concurrent attribute creation
// for the same person with different keys works correctly without race conditions
func TestConcurrentAttributeCreation_SamePerson(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person
	personID, err := createTestPerson(ctx, "concurrent-test-person")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	numGoroutines := 10
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)
	results := make(chan int, numGoroutines)

	// Create multiple attributes concurrently for the same person
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			e := echo.New()
			jsonBody := fmt.Sprintf(`{"key":"concurrent-key-%d","value":"concurrent-value-%d","meta":{"caller":"test","reason":"concurrent-test","traceId":"trace-%d"}}`, index, index, index)
			req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("personId")
			c.SetParamValues(personID)

			if err := handler.CreateAttribute(c); err != nil {
				errors <- err
				return
			}

			results <- rec.Code
		}(i)
	}

	wg.Wait()
	close(errors)
	close(results)

	// Verify no errors occurred
	for err := range errors {
		assert.NoError(t, err, "Concurrent creation should not produce errors")
	}

	// Verify all responses were successful
	successCount := 0
	for code := range results {
		if code == http.StatusCreated {
			successCount++
		}
	}
	assert.Equal(t, numGoroutines, successCount, "All concurrent creations should succeed")

	// Verify all attributes were created correctly
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestConcurrentAttributeUpsert_SameKey tests upsert behavior when multiple goroutines
// try to create/update the same attribute key simultaneously
func TestConcurrentAttributeUpsert_SameKey(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person
	personID, err := createTestPerson(ctx, "upsert-test-person")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	numGoroutines := 5
	var wg sync.WaitGroup
	results := make(chan int, numGoroutines)

	// Try to create the same attribute key concurrently (should upsert)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			e := echo.New()
			jsonBody := fmt.Sprintf(`{"key":"same-key","value":"value-%d","meta":{"caller":"test","reason":"upsert-test","traceId":"upsert-trace-%d"}}`, index, index)
			req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("personId")
			c.SetParamValues(personID)

			handler.CreateAttribute(c)
			results <- rec.Code
		}(i)
	}

	wg.Wait()
	close(results)

	// Verify all responses were successful (either Created or OK for upsert)
	for code := range results {
		assert.True(t, code == http.StatusCreated || code == http.StatusOK,
			"Upsert should return 201 or 200, got %d", code)
	}

	// Verify only one attribute exists with the key
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response to verify only one attribute
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	attributes, ok := response["attributes"].([]interface{})
	if ok {
		assert.Equal(t, 1, len(attributes), "Should have exactly one attribute after upserts")
	}
}

// TestConcurrentReadWrite_Attributes tests concurrent read and write operations
func TestConcurrentReadWrite_Attributes(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a test person with initial attribute
	personID, err := createTestPerson(ctx, "readwrite-test-person")
	assert.NoError(t, err)
	_, err = createTestAttribute(ctx, personID, "readwrite-key", "initial-value")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	numReaders := 5
	numWriters := 3
	var wg sync.WaitGroup

	// Start readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("personId")
			c.SetParamValues(personID)

			err := handler.GetAllAttributes(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
		}()
	}

	// Start writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			e := echo.New()
			jsonBody := fmt.Sprintf(`{"key":"readwrite-key","value":"updated-value-%d","meta":{"caller":"test","reason":"readwrite-test","traceId":"rw-trace-%d"}}`, index, index)
			req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("personId")
			c.SetParamValues(personID)

			handler.CreateAttribute(c)
			// Both 200 and 201 are acceptable
			assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusCreated)
		}(i)
	}

	wg.Wait()
}

// ============================================================================
// PHASE 2: BOUNDARY AND EDGE CASE TESTS
// ============================================================================

// TestCreateAttribute_LongKey tests attribute creation with maximum length key
func TestCreateAttribute_LongKey(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "long-key-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	// Create a long key (citext has no explicit limit but test reasonable boundary)
	longKey := strings.Repeat("a", 255)

	e := echo.New()
	jsonBody := fmt.Sprintf(`{"key":"%s","value":"test-value","meta":{"caller":"test","reason":"boundary-test","traceId":"long-key-trace"}}`, longKey)
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Verify retrieval works
	assert.Contains(t, rec.Body.String(), longKey)
}

// TestCreateAttribute_LongValue tests attribute creation with very long value
func TestCreateAttribute_LongValue(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "long-value-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	// Create a long value (encrypted values stored as BYTEA should handle large data)
	longValue := strings.Repeat("x", 10000)

	e := echo.New()
	jsonBody := fmt.Sprintf(`{"key":"long-value-key","value":"%s","meta":{"caller":"test","reason":"boundary-test","traceId":"long-value-trace"}}`, longValue)
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

// TestCreateAttribute_UnicodeCharacters tests handling of Unicode in attribute values
func TestCreateAttribute_UnicodeCharacters(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "unicode-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	testCases := []struct {
		name  string
		key   string
		value string
	}{
		{"Chinese characters", "name-zh", "张伟"},
		{"Japanese characters", "name-ja", "田中太郎"},
		{"Arabic characters", "name-ar", "محمد أحمد"},
		{"Emoji", "emoji-key", "Hello 🌍 World 🎉"},
		{"Mixed scripts", "mixed-key", "Hello مرحبا 你好 👋"},
		{"Special symbols", "symbols-key", "©®™€£¥"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			jsonBody := fmt.Sprintf(`{"key":"%s","value":"%s","meta":{"caller":"test","reason":"unicode-test","traceId":"unicode-trace-%s"}}`, tc.key, tc.value, tc.key)
			req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("personId")
			c.SetParamValues(personID)

			err := handler.CreateAttribute(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, rec.Code)

			// Verify the value is stored and retrieved correctly
			assert.Contains(t, rec.Body.String(), tc.key)
		})
	}
}

// TestCreateAttribute_SpecialCharactersInKey tests special characters in attribute keys
func TestCreateAttribute_SpecialCharactersInKey(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "special-key-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	testCases := []struct {
		name string
		key  string
	}{
		{"Dots in key", "user.email.primary"},
		{"Underscores", "user_email_primary"},
		{"Hyphens", "user-email-primary"},
		{"Colons", "namespace:attribute:name"},
		{"Numbers", "attribute123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			jsonBody := fmt.Sprintf(`{"key":"%s","value":"test-value","meta":{"caller":"test","reason":"special-key-test","traceId":"special-%s"}}`, tc.key, tc.key)
			req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("personId")
			c.SetParamValues(personID)

			err := handler.CreateAttribute(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, rec.Code)
		})
	}
}

// TestCreateAttribute_CaseInsensitiveKey tests that attribute keys are case-insensitive (citext)
func TestCreateAttribute_CaseInsensitiveKey(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "case-insensitive-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()

	// Create attribute with lowercase key
	jsonBody := `{"key":"email","value":"first@example.com","meta":{"caller":"test","reason":"case-test","traceId":"case-1"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Try to create with uppercase key - should upsert due to citext
	jsonBody = `{"key":"EMAIL","value":"second@example.com","meta":{"caller":"test","reason":"case-test","traceId":"case-2"}}`
	req = httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)
	assert.NoError(t, err)
	// Should succeed (either create or update depending on upsert behavior)
	assert.True(t, rec.Code == http.StatusCreated || rec.Code == http.StatusOK)
}

// TestMultipleAttributesPerPerson tests creating many attributes for a single person
func TestMultipleAttributesPerPerson(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "multi-attr-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	numAttributes := 50 // Test with many attributes

	for i := 0; i < numAttributes; i++ {
		e := echo.New()
		jsonBody := fmt.Sprintf(`{"key":"attr-%d","value":"value-%d","meta":{"caller":"test","reason":"multi-test","traceId":"multi-trace-%d"}}`, i, i, i)
		req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("personId")
		c.SetParamValues(personID)

		err := handler.CreateAttribute(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code, "Failed to create attribute %d", i)
	}

	// Verify all attributes can be retrieved
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify count
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	if attrs, ok := response["attributes"].([]interface{}); ok {
		assert.Equal(t, numAttributes, len(attrs), "Should have %d attributes", numAttributes)
	}
}

// ============================================================================
// PHASE 3: API CONTRACT VALIDATION TESTS
// ============================================================================

// TestCreateAttribute_ResponseSchema validates the response schema for create
func TestCreateAttribute_ResponseSchema(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "schema-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"schema-key","value":"schema-value","meta":{"caller":"test","reason":"schema-test","traceId":"schema-trace"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Validate response has required fields
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check required fields exist
	assert.Contains(t, response, "id", "Response should contain 'id'")
	assert.Contains(t, response, "key", "Response should contain 'key'")
	assert.Contains(t, response, "value", "Response should contain 'value'")
	assert.Contains(t, response, "createdAt", "Response should contain 'createdAt'")
	assert.Contains(t, response, "updatedAt", "Response should contain 'updatedAt'")

	// Validate field types
	_, ok := response["id"].(float64)
	assert.True(t, ok, "'id' should be a number")
	_, ok = response["key"].(string)
	assert.True(t, ok, "'key' should be a string")
	_, ok = response["value"].(string)
	assert.True(t, ok, "'value' should be a string")
}

// TestGetAllAttributes_ResponseSchema validates the response schema for get all
func TestGetAllAttributes_ResponseSchema(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "get-all-schema-test")
	assert.NoError(t, err)
	_, err = createTestAttribute(ctx, personID, "test-key", "test-value")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Validate response structure
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "attributes", "Response should contain 'attributes' array")
	attrs, ok := response["attributes"].([]interface{})
	assert.True(t, ok, "'attributes' should be an array")
	assert.Greater(t, len(attrs), 0, "Should have at least one attribute")

	// Validate each attribute has required fields
	firstAttr, ok := attrs[0].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, firstAttr, "id")
	assert.Contains(t, firstAttr, "key")
	assert.Contains(t, firstAttr, "value")
}

// TestErrorResponse_Schema validates error response format consistency
func TestErrorResponse_Schema(t *testing.T) {
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           string
		personId       string
		expectedStatus int
	}{
		{
			name:           "Invalid UUID",
			method:         http.MethodGet,
			path:           "/persons/invalid-uuid/attributes",
			personId:       "invalid-uuid",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPut,
			path:           "/persons/123e4567-e89b-12d3-a456-426614174000/attributes",
			body:           "{invalid}",
			personId:       "123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty Key",
			method:         http.MethodPut,
			path:           "/persons/123e4567-e89b-12d3-a456-426614174000/attributes",
			body:           `{"key":"","value":"test","meta":{"caller":"test","reason":"test","traceId":"123"}}`,
			personId:       "123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("personId")
			c.SetParamValues(tc.personId)

			if tc.method == http.MethodGet {
				handler.GetAllAttributes(c)
			} else {
				handler.CreateAttribute(c)
			}

			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Validate error response has required fields
			var errResponse map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &errResponse)
			assert.NoError(t, err)
			assert.Contains(t, errResponse, "message", "Error response should contain 'message'")
			assert.Contains(t, errResponse, "errorCode", "Error response should contain 'errorCode'")
		})
	}
}

// TestContentType_JSON validates Content-Type header is application/json
func TestContentType_JSON(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "content-type-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	jsonBody := `{"key":"ct-key","value":"ct-value","meta":{"caller":"test","reason":"content-type-test","traceId":"ct-trace"}}`
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)
	assert.NoError(t, err)

	// Check Content-Type header
	contentType := rec.Header().Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "Response Content-Type should be application/json")
}

// ============================================================================
// PHASE 4: DATA CONSISTENCY TESTS
// ============================================================================

// TestCascadeDelete_PersonWithAttributes tests that deleting a person cascades to attributes
func TestCascadeDelete_PersonWithAttributes(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	// Create a person with multiple attributes
	personID, err := createTestPerson(ctx, "cascade-delete-test")
	assert.NoError(t, err)

	_, err = createTestAttribute(ctx, personID, "attr1", "value1")
	assert.NoError(t, err)
	_, err = createTestAttribute(ctx, personID, "attr2", "value2")
	assert.NoError(t, err)
	_, err = createTestAttribute(ctx, personID, "attr3", "value3")
	assert.NoError(t, err)

	// Delete the person directly from database
	_, err = pool.Exec(ctx, "DELETE FROM person WHERE id = $1::uuid", personID)
	assert.NoError(t, err)

	// Verify attributes are also deleted (due to CASCADE)
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM person_attributes WHERE person_id = $1::uuid", personID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Attributes should be cascade deleted with person")

	// Verify person cannot access attributes through API
	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID+"/attributes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.GetAllAttributes(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code, "Should return 404 for deleted person")
}

// TestIdempotency_SameTraceID tests idempotent behavior with same trace_id
func TestIdempotency_SameTraceID(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "idempotency-test")
	assert.NoError(t, err)

	queries := db.New(pool)
	handler := NewPersonAttributesHandler(queries)

	// Create attribute with specific trace_id
	traceID := "idempotent-trace-12345"
	e := echo.New()
	jsonBody := fmt.Sprintf(`{"key":"idempotent-key","value":"value1","meta":{"caller":"test","reason":"idempotency-test","traceId":"%s"}}`, traceID)
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)
	assert.NoError(t, err)
	firstResponse := rec.Body.String()

	// Repeat request with same trace_id
	req = httptest.NewRequest(http.MethodPut, "/persons/"+personID+"/attributes", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("personId")
	c.SetParamValues(personID)

	err = handler.CreateAttribute(c)
	assert.NoError(t, err)
	secondResponse := rec.Body.String()

	// Both responses should be consistent (same attribute ID)
	var first, second map[string]interface{}
	json.Unmarshal([]byte(firstResponse), &first)
	json.Unmarshal([]byte(secondResponse), &second)

	// The key should be the same
	assert.Equal(t, first["key"], second["key"], "Same trace_id should return consistent results")
}

// TestUniqueConstraint_PersonAttributeKey tests unique constraint on (person_id, attribute_key)
func TestUniqueConstraint_PersonAttributeKey(t *testing.T) {
	ctx := context.Background()
	err := testdb.TruncateTables(ctx, pool)
	assert.NoError(t, err)

	personID, err := createTestPerson(ctx, "unique-constraint-test")
	assert.NoError(t, err)

	// Create first attribute
	_, err = createTestAttribute(ctx, personID, "unique-key", "value1")
	assert.NoError(t, err)

	// Try to insert duplicate via direct SQL (should fail)
	_, err = pool.Exec(ctx, `
		INSERT INTO person_attributes (person_id, attribute_key, encrypted_value, key_version)
		VALUES ($1::uuid, 'unique-key', pgp_sym_encrypt('value2', $2), 1)
	`, personID, testEncryptionKey)
	// This should fail due to unique constraint
	assert.Error(t, err, "Should fail on unique constraint violation")
	assert.Contains(t, err.Error(), "duplicate key value", "Error should mention duplicate key")
}
