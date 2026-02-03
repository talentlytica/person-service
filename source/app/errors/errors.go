package errors

// ErrorResponse represents a standard error response with a unique error code
type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

// Error codes for Person Attributes endpoints
const (
	// Validation errors (1000-1099)
	ErrInvalidPersonID          = "PA_001_INVALID_PERSON_ID"
	ErrInvalidAttributeID       = "PA_002_INVALID_ATTRIBUTE_ID"
	ErrInvalidRequestBody       = "PA_003_INVALID_REQUEST_BODY"
	ErrMissingRequiredFieldKey  = "PA_004_MISSING_KEY"
	ErrMissingRequiredFieldMeta = "PA_005_MISSING_META"
	ErrInvalidAttributeIDFormat = "PA_006_INVALID_ATTRIBUTE_ID_FORMAT"

	// Resource not found errors (1100-1199)
	ErrPersonNotFound    = "PA_101_PERSON_NOT_FOUND"
	ErrAttributeNotFound = "PA_102_ATTRIBUTE_NOT_FOUND"

	// Database operation errors (1200-1299)
	ErrFailedVerifyPerson        = "PA_201_FAILED_VERIFY_PERSON"
	ErrFailedCreateAttribute     = "PA_202_FAILED_CREATE_ATTRIBUTE"
	ErrFailedRetrieveAttribute   = "PA_203_FAILED_RETRIEVE_ATTRIBUTE"
	ErrFailedRetrieveAttributes  = "PA_204_FAILED_RETRIEVE_ATTRIBUTES"
	ErrFailedUpdateAttribute     = "PA_205_FAILED_UPDATE_ATTRIBUTE"
	ErrFailedRetrieveUpdatedAttr = "PA_206_FAILED_RETRIEVE_UPDATED"
	ErrFailedDeleteAttribute     = "PA_207_FAILED_DELETE_ATTRIBUTE"
	ErrFailedUpdateAttributeKey  = "PA_208_FAILED_UPDATE_KEY"

	// Audit logging errors (1300-1399)
	ErrFailedAuditLog = "PA_301_FAILED_AUDIT_LOG"
)

// Error codes for Key-Value endpoints
const (
	// Validation errors (2000-2099)
	ErrKVInvalidRequestBody = "KV_001_INVALID_REQUEST_BODY"
	ErrKVMissingKeyOrValue  = "KV_002_MISSING_KEY_OR_VALUE"
	ErrKVMissingKeyParam    = "KV_003_MISSING_KEY_PARAM"

	// Resource not found errors (2100-2199)
	ErrKVKeyNotFound = "KV_101_KEY_NOT_FOUND"

	// Database operation errors (2200-2299)
	ErrKVFailedSetValue      = "KV_201_FAILED_SET_VALUE"
	ErrKVFailedRetrieveValue = "KV_202_FAILED_RETRIEVE_VALUE"
	ErrKVFailedDeleteValue   = "KV_203_FAILED_DELETE_VALUE"
)

// Error codes for API Key middleware
const (
	// Authentication errors (3000-3099)
	ErrMissingAPIKey        = "API_001_MISSING_API_KEY"
	ErrInvalidAPIKeyFormat  = "API_002_INVALID_API_KEY_FORMAT"
	ErrAPIKeysNotConfigured = "API_003_KEYS_NOT_CONFIGURED"
	ErrInvalidAPIKey        = "API_004_INVALID_API_KEY"
)

// Error codes for Health Check
const (
	// Health check errors (4000-4099)
	ErrHealthCheckFailed = "HC_001_HEALTH_CHECK_FAILED"
)

// Error codes for Database Setup
const (
	// Database connection errors (5000-5099)
	ErrDatabaseURLNotSet    = "DB_001_URL_NOT_SET"
	ErrInvalidPort          = "DB_002_INVALID_PORT"
	ErrFailedParseDBURL     = "DB_003_FAILED_PARSE_URL"
	ErrFailedCreateConnPool = "DB_004_FAILED_CREATE_POOL"
	ErrDatabasePingFailed   = "DB_005_PING_FAILED"
	ErrFailedStartServer    = "DB_006_FAILED_START_SERVER"
	ErrFailedShutdownServer = "DB_007_FAILED_SHUTDOWN_SERVER"
)
