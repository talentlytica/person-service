package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	db "person-service/internal/db/generated"

	health "person-service/healthcheck"
	key_value "person-service/key_value"
	"person-service/middleware"
	person_attributes "person-service/person_attributes"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	// TestEncryptionKey is the encryption key used for tests
	TestEncryptionKey = "test-encryption-key-32bytes!!"
	// TestAPIKeyBlue is a valid API key for tests (blue)
	TestAPIKeyBlue = "person-service-key-11111111-2222-3333-4444-555555555555"
	// TestAPIKeyGreen is a valid API key for tests (green)
	TestAPIKeyGreen = "person-service-key-aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
)

// TestServer wraps an Echo instance configured for testing
type TestServer struct {
	Echo    *echo.Echo
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

// NewTestServer creates a new test server with all handlers configured
func NewTestServer(pool *pgxpool.Pool) *TestServer {
	// Set up environment variables for the handlers
	os.Setenv("ENCRYPTION_KEY_1", TestEncryptionKey)
	os.Setenv("PERSON_API_KEY_BLUE", TestAPIKeyBlue)
	os.Setenv("PERSON_API_KEY_GREEN", TestAPIKeyGreen)

	queries := db.New(pool)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Setup handlers
	healthHandler := health.NewHealthCheckHandler(queries)
	keyValueHandler := key_value.NewKeyValueHandler(queries)
	personAttributesHandler := person_attributes.NewPersonAttributesHandler(queries)

	// Setup routes (same as main.go)
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

	return &TestServer{
		Echo:    e,
		Pool:    pool,
		Queries: queries,
	}
}

// Request executes an HTTP request against the test server
func (ts *TestServer) Request(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonBytes)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	ts.Echo.ServeHTTP(rec, req)
	return rec
}

// GET executes a GET request
func (ts *TestServer) GET(path string, headers map[string]string) *httptest.ResponseRecorder {
	return ts.Request(http.MethodGet, path, nil, headers)
}

// POST executes a POST request
func (ts *TestServer) POST(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return ts.Request(http.MethodPost, path, body, headers)
}

// PUT executes a PUT request
func (ts *TestServer) PUT(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return ts.Request(http.MethodPut, path, body, headers)
}

// DELETE executes a DELETE request
func (ts *TestServer) DELETE(path string, headers map[string]string) *httptest.ResponseRecorder {
	return ts.Request(http.MethodDelete, path, nil, headers)
}

// WithAPIKey returns headers with the blue API key
func WithAPIKey() map[string]string {
	return map[string]string{
		"x-api-key": TestAPIKeyBlue,
	}
}

// WithGreenAPIKey returns headers with the green API key
func WithGreenAPIKey() map[string]string {
	return map[string]string{
		"x-api-key": TestAPIKeyGreen,
	}
}

// WithCustomAPIKey returns headers with a custom API key
func WithCustomAPIKey(key string) map[string]string {
	return map[string]string{
		"x-api-key": key,
	}
}

// ParseJSONResponse parses the response body as JSON into the given target
func ParseJSONResponse(rec *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(rec.Body.Bytes(), target)
}

// GetJSONField extracts a field from a JSON response
func GetJSONField(rec *httptest.ResponseRecorder, field string) (interface{}, error) {
	var result map[string]interface{}
	if err := ParseJSONResponse(rec, &result); err != nil {
		return nil, err
	}
	return result[field], nil
}
