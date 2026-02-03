package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	validAPIKeyBlue  = "person-service-key-11111111-2222-3333-4444-555555555555"
	validAPIKeyGreen = "person-service-key-aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
)

func TestAPIKeyMiddleware_MissingKey(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := APIKeyMiddleware()
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Missing required header")
}

func TestAPIKeyMiddleware_InvalidFormat(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", "invalid-format-key")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := APIKeyMiddleware()
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid API key format")
}

func TestAPIKeyMiddleware_NoConfiguredKeys(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("PERSON_API_KEY_BLUE")
	os.Unsetenv("PERSON_API_KEY_GREEN")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", validAPIKeyBlue)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := APIKeyMiddleware()
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "not properly configured")
}

func TestAPIKeyMiddleware_InvalidKey(t *testing.T) {
	os.Setenv("PERSON_API_KEY_BLUE", validAPIKeyBlue)
	defer os.Unsetenv("PERSON_API_KEY_BLUE")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", "person-service-key-99999999-8888-7777-6666-555555555555")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := APIKeyMiddleware()
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid API key")
}

func TestAPIKeyMiddleware_ValidBlueKey(t *testing.T) {
	os.Setenv("PERSON_API_KEY_BLUE", validAPIKeyBlue)
	defer os.Unsetenv("PERSON_API_KEY_BLUE")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", validAPIKeyBlue)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := APIKeyMiddleware()
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestAPIKeyMiddleware_ValidGreenKey(t *testing.T) {
	os.Setenv("PERSON_API_KEY_GREEN", validAPIKeyGreen)
	defer os.Unsetenv("PERSON_API_KEY_GREEN")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", validAPIKeyGreen)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := APIKeyMiddleware()
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestAPIKeyMiddleware_BothKeysConfigured(t *testing.T) {
	os.Setenv("PERSON_API_KEY_BLUE", validAPIKeyBlue)
	os.Setenv("PERSON_API_KEY_GREEN", validAPIKeyGreen)
	defer os.Unsetenv("PERSON_API_KEY_BLUE")
	defer os.Unsetenv("PERSON_API_KEY_GREEN")

	// Test with blue key
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", validAPIKeyBlue)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := APIKeyMiddleware()
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Test with green key
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("x-api-key", validAPIKeyGreen)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	err2 := handler(c2)

	assert.NoError(t, err2)
	assert.Equal(t, http.StatusOK, rec2.Code)
}

func TestAPIKeyMiddleware_BlueKeyActiveGreenInvalidFormat(t *testing.T) {
	os.Setenv("PERSON_API_KEY_BLUE", validAPIKeyBlue)
	os.Setenv("PERSON_API_KEY_GREEN", "invalid-format")
	defer os.Unsetenv("PERSON_API_KEY_BLUE")
	defer os.Unsetenv("PERSON_API_KEY_GREEN")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", validAPIKeyBlue)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := APIKeyMiddleware()
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
