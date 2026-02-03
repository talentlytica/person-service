package middleware

import (
	"net/http"
	"os"
	"regexp"

	errs "person-service/errors"

	"github.com/labstack/echo/v4"
)

// API key format: person-service-key-<UUID>
// UUID format: 8-4-4-4-12 hexadecimal characters
var apiKeyPattern = regexp.MustCompile(`^person-service-key-[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// APIKeyMiddleware creates a middleware that validates the x-api-key header
// against PERSON_API_KEY_BLUE and PERSON_API_KEY_GREEN environment variables.
// The API key must follow the format: person-service-key-<UUID>
func APIKeyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("x-api-key")

			// Check if API key is provided
			if apiKey == "" {
				return c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
					Message:   "Missing required header \"x-api-key\"",
					ErrorCode: errs.ErrMissingAPIKey,
				})
			}

			// Validate the format of the provided API key
			if !apiKeyPattern.MatchString(apiKey) {
				return c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
					Message:   "Invalid API key format",
					ErrorCode: errs.ErrInvalidAPIKeyFormat,
				})
			}

			// Get the configured API keys from environment
			apiKeyBlue := os.Getenv("PERSON_API_KEY_BLUE")
			apiKeyGreen := os.Getenv("PERSON_API_KEY_GREEN")

			// Check if blue key is active (valid format)
			blueActive := apiKeyPattern.MatchString(apiKeyBlue)
			// Check if green key is active (valid format)
			greenActive := apiKeyPattern.MatchString(apiKeyGreen)

			// If neither key is active (properly configured), reject the request
			if !blueActive && !greenActive {
				return c.JSON(http.StatusServiceUnavailable, errs.ErrorResponse{
					Message:   "API keys are not properly configured",
					ErrorCode: errs.ErrAPIKeysNotConfigured,
				})
			}

			// Validate the provided key against active keys
			keyValid := false
			if blueActive && apiKey == apiKeyBlue {
				keyValid = true
			}
			if greenActive && apiKey == apiKeyGreen {
				keyValid = true
			}

			if !keyValid {
				return c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
					Message:   "Invalid API key",
					ErrorCode: errs.ErrInvalidAPIKey,
				})
			}

			return next(c)
		}
	}
}
