package key_value

import (
	"context"
	"errors"
	"net/http"
	db "person-service/internal/db/generated"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// SetValueRequest represents the request body for setting a key-value pair
type SetValueRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// KeyValueHandler handles KeyValue
type KeyValueHandler struct {
	queries *db.Queries
}

// KeyValueHandler creates a new instance of KeyValueHandler with injected queries
func NewKeyValueHandler(queries *db.Queries) *KeyValueHandler {
	return &KeyValueHandler{
		queries: queries,
	}
}

// SetValue handles POST /api/key_value - sets or updates a key-value pair
func (h *KeyValueHandler) SetValue(c echo.Context) error {
	// Parse request body
	var req SetValueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Key == "" || req.Value == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Key and value are required",
		})
	}

	// Set value in database
	ctx := context.Background()
	err := h.queries.SetValue(ctx, db.SetValueParams{
		Key:   req.Key,
		Value: req.Value,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to set value",
		})
	}

	// Retrieve the full record with timestamps
	record, err := h.queries.GetKeyValue(ctx, req.Key)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve value",
		})
	}

	// Return success with the full key-value record
	response := map[string]interface{}{
		"key":   record.Key,
		"value": record.Value,
	}

	// Add timestamps if they are valid
	if record.CreatedAt.Valid {
		response["created_at"] = record.CreatedAt.Time
	}
	if record.UpdatedAt.Valid {
		response["updated_at"] = record.UpdatedAt.Time
	}

	return c.JSON(http.StatusOK, response)
}

// GetValue handles GET /api/key_value/:key - retrieves a value by key
func (h *KeyValueHandler) GetValue(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Key parameter is required",
		})
	}

	// Get full record from database
	ctx := context.Background()
	record, err := h.queries.GetKeyValue(ctx, key)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Key not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve value",
		})
	}

	// Return the full key-value record
	response := map[string]interface{}{
		"key":   record.Key,
		"value": record.Value,
	}

	// Add timestamps if they are valid
	if record.CreatedAt.Valid {
		response["created_at"] = record.CreatedAt.Time
	}
	if record.UpdatedAt.Valid {
		response["updated_at"] = record.UpdatedAt.Time
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteValue handles DELETE /api/key_value/:key - deletes a key-value pair
func (h *KeyValueHandler) DeleteValue(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Key parameter is required",
		})
	}

	// Delete value from database
	ctx := context.Background()
	err := h.queries.DeleteValue(ctx, key)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete value",
		})
	}

	// Return success message
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Key deleted successfully",
	})
}
