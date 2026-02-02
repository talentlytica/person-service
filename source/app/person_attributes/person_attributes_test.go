package person_attributes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
