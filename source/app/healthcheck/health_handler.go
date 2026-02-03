package health

import (
	"context"
	"net/http"

	errs "person-service/errors"
	db "person-service/internal/db/generated"

	"github.com/labstack/echo/v4"
)

// HealthCheckService handles health check operations
type HealthCheckHandler struct {
	queries *db.Queries
}

// NewHealthCheckHandler creates a new instance of HealthCheckService with injected queries
func NewHealthCheckHandler(queries *db.Queries) *HealthCheckHandler {
	return &HealthCheckHandler{
		queries: queries,
	}
}

// Check performs a health check on the database and returns the result as an echo handler
func (h *HealthCheckHandler) Check(c echo.Context) error {
	// Call HealthCheck directly
	err := h.queries.HealthCheck(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Message:   err.Error(),
			ErrorCode: errs.ErrHealthCheckFailed,
		})
	}

	// Return healthy status
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "healthy",
	})
}
