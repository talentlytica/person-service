package middleware

import (
	"net/http"
	"os"
	"strings"

	errs "person-service/errors"

	"github.com/labstack/echo/v4"
)

// BearerMiddleware creates a middleware that validates the Authorization: Bearer <token> header
// against PERSON_API_KEY_BLUE and PERSON_API_KEY_GREEN environment variables.
func BearerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
					Message:   "Missing required Authorization header",
					ErrorCode: errs.ErrMissingBearerToken,
				})
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
					Message:   "Invalid Authorization header format, expected Bearer token",
					ErrorCode: errs.ErrInvalidBearerFormat,
				})
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			if !apiKeyPattern.MatchString(token) {
				return c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
					Message:   "Invalid token format",
					ErrorCode: errs.ErrInvalidBearerFormat,
				})
			}

			apiKeyBlue := os.Getenv("PERSON_API_KEY_BLUE")
			apiKeyGreen := os.Getenv("PERSON_API_KEY_GREEN")

			blueActive := apiKeyPattern.MatchString(apiKeyBlue)
			greenActive := apiKeyPattern.MatchString(apiKeyGreen)

			if !blueActive && !greenActive {
				return c.JSON(http.StatusServiceUnavailable, errs.ErrorResponse{
					Message:   "API keys are not properly configured",
					ErrorCode: errs.ErrAPIKeysNotConfigured,
				})
			}

			keyValid := false
			if blueActive && token == apiKeyBlue {
				keyValid = true
			}
			if greenActive && token == apiKeyGreen {
				keyValid = true
			}

			if !keyValid {
				return c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
					Message:   "Invalid token",
					ErrorCode: errs.ErrInvalidAPIKey,
				})
			}

			return next(c)
		}
	}
}
