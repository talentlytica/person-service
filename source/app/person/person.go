package person

import (
	"errors"
	"fmt"
	"net/http"

	errs "person-service/errors"
	db "person-service/internal/db/generated"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// CreatePersonRequest represents the request body for creating a person
type CreatePersonRequest struct {
	ClientID string `json:"client_id"`
}

// UpdatePersonRequest represents the request body for updating a person
type UpdatePersonRequest struct {
	ClientID string `json:"client_id"`
}

// PersonHandler handles Person CRUD operations
type PersonHandler struct {
	queries *db.Queries
}

// NewPersonHandler creates a new instance of PersonHandler with injected queries
func NewPersonHandler(queries *db.Queries) *PersonHandler {
	return &PersonHandler{
		queries: queries,
	}
}

// formatUUID converts pgtype.UUID to a string representation
func formatUUID(u pgtype.UUID) string {
	b := u.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// parseUUID parses a UUID string into pgtype.UUID
func parseUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	err := u.Scan(s)
	return u, err
}

// CreatePerson handles POST /api/person - creates a new person
func (h *PersonHandler) CreatePerson(c echo.Context) error {
	var req CreatePersonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Message:   "Invalid request body",
			ErrorCode: errs.ErrPersonInvalidRequestBody,
		})
	}

	if req.ClientID == "" {
		return c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Message:   "client_id is required",
			ErrorCode: errs.ErrPersonMissingClientID,
		})
	}

	ctx := c.Request().Context()

	person, err := h.queries.CreatePerson(ctx, req.ClientID)
	if err != nil {
		if isDuplicateKeyError(err) {
			return c.JSON(http.StatusConflict, errs.ErrorResponse{
				Message:   "A person with this client_id already exists",
				ErrorCode: errs.ErrPersonDuplicateClientID,
			})
		}
		return c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Message:   "Failed to create person",
			ErrorCode: errs.ErrPersonFailedCreate,
		})
	}

	response := map[string]interface{}{
		"data": buildPersonResponse(person),
	}

	return c.JSON(http.StatusCreated, response)
}

// GetPerson handles GET /api/person/:id - retrieves a person by ID
func (h *PersonHandler) GetPerson(c echo.Context) error {
	id := c.Param("id")

	personID, err := parseUUID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Message:   "Invalid person ID format",
			ErrorCode: errs.ErrPersonInvalidID,
		})
	}

	ctx := c.Request().Context()

	person, err := h.queries.GetPersonById(ctx, personID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, errs.ErrorResponse{
				Message:   "Person not found",
				ErrorCode: errs.ErrPersonNotFoundCRUD,
			})
		}
		return c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Message:   "Failed to retrieve person",
			ErrorCode: errs.ErrPersonFailedRetrieve,
		})
	}

	response := map[string]interface{}{
		"data": buildPersonResponse(person),
	}

	return c.JSON(http.StatusOK, response)
}

// UpdatePerson handles PATCH /api/person/:id - updates a person's client_id
func (h *PersonHandler) UpdatePerson(c echo.Context) error {
	id := c.Param("id")

	personID, err := parseUUID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Message:   "Invalid person ID format",
			ErrorCode: errs.ErrPersonInvalidID,
		})
	}

	var req UpdatePersonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Message:   "Invalid request body",
			ErrorCode: errs.ErrPersonInvalidRequestBody,
		})
	}

	if req.ClientID == "" {
		return c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Message:   "client_id is required",
			ErrorCode: errs.ErrPersonMissingClientID,
		})
	}

	ctx := c.Request().Context()

	// Verify person exists
	_, err = h.queries.GetPersonById(ctx, personID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, errs.ErrorResponse{
				Message:   "Person not found",
				ErrorCode: errs.ErrPersonNotFoundCRUD,
			})
		}
		return c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Message:   "Failed to retrieve person",
			ErrorCode: errs.ErrPersonFailedRetrieve,
		})
	}

	// Update client_id
	err = h.queries.UpdatePersonClientId(ctx, db.UpdatePersonClientIdParams{
		NewClientID: req.ClientID,
		ID:          personID,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return c.JSON(http.StatusConflict, errs.ErrorResponse{
				Message:   "A person with this client_id already exists",
				ErrorCode: errs.ErrPersonDuplicateClientID,
			})
		}
		return c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Message:   "Failed to update person",
			ErrorCode: errs.ErrPersonFailedUpdate,
		})
	}

	// Fetch updated person
	updated, err := h.queries.GetPersonById(ctx, personID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Message:   "Failed to retrieve updated person",
			ErrorCode: errs.ErrPersonFailedRetrieve,
		})
	}

	response := map[string]interface{}{
		"data": buildPersonResponse(updated),
	}

	return c.JSON(http.StatusOK, response)
}

// DeletePerson handles DELETE /api/person/:id - soft deletes a person
func (h *PersonHandler) DeletePerson(c echo.Context) error {
	id := c.Param("id")

	personID, err := parseUUID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Message:   "Invalid person ID format",
			ErrorCode: errs.ErrPersonInvalidID,
		})
	}

	ctx := c.Request().Context()

	// Verify person exists
	_, err = h.queries.GetPersonById(ctx, personID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, errs.ErrorResponse{
				Message:   "Person not found",
				ErrorCode: errs.ErrPersonNotFoundCRUD,
			})
		}
		return c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Message:   "Failed to retrieve person",
			ErrorCode: errs.ErrPersonFailedRetrieve,
		})
	}

	// Soft delete
	err = h.queries.SoftDeletePerson(ctx, personID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Message:   "Failed to delete person",
			ErrorCode: errs.ErrPersonFailedDelete,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Person deleted successfully",
	})
}

// buildPersonResponse creates a response map from a Person model
func buildPersonResponse(p db.Person) map[string]interface{} {
	resp := map[string]interface{}{
		"id":        formatUUID(p.ID),
		"client_id": p.ClientID,
	}
	if p.CreatedAt.Valid {
		resp["created_at"] = p.CreatedAt.Time
	}
	if p.UpdatedAt.Valid {
		resp["updated_at"] = p.UpdatedAt.Time
	}
	if p.DeletedAt.Valid {
		resp["deleted_at"] = p.DeletedAt.Time
	}
	return resp
}

// isDuplicateKeyError checks if the error is a PostgreSQL unique constraint violation (23505)
func isDuplicateKeyError(err error) bool {
	return err != nil && contains(err.Error(), "23505")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	// Simple substring search
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
