# Error Codes Reference

This document outlines all unique error codes used throughout the person-service application. Each error has a specific code for pinpointing issues during debugging and monitoring.

## Error Response Format

All error responses now follow a standardized format:

```json
{
  "message": "Human-readable error message",
  "error_code": "UNIQUE_ERROR_CODE"
}
```

## Error Code Categories

### Person Attributes Endpoints (PA_*)

#### Validation Errors (PA_001-PA_006)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| PA_001_INVALID_PERSON_ID | 404/400 | Invalid person ID format in path parameter |
| PA_002_INVALID_ATTRIBUTE_ID | 400 | Invalid attribute ID format in path parameter |
| PA_003_INVALID_REQUEST_BODY | 400 | Request body is malformed or invalid JSON |
| PA_004_MISSING_KEY | 400 | Required "key" field is missing in request body |
| PA_005_MISSING_META | 400 | Required "meta" field is missing in request body |
| PA_006_INVALID_ATTRIBUTE_ID_FORMAT | 400 | Attribute ID cannot be parsed as integer |

#### Resource Not Found Errors (PA_101-PA_102)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| PA_101_PERSON_NOT_FOUND | 404 | Specified person ID does not exist in database |
| PA_102_ATTRIBUTE_NOT_FOUND | 404 | Specified attribute ID does not exist for person |

#### Database Operation Errors (PA_201-PA_208)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| PA_201_FAILED_VERIFY_PERSON | 500 | Error verifying if person exists in database |
| PA_202_FAILED_CREATE_ATTRIBUTE | 500 | Error creating/updating person attribute in database |
| PA_203_FAILED_RETRIEVE_ATTRIBUTE | 500 | Error retrieving single attribute after creation |
| PA_204_FAILED_RETRIEVE_ATTRIBUTES | 500 | Error retrieving all attributes for person |
| PA_205_FAILED_UPDATE_ATTRIBUTE | 500 | Error updating attribute value in database |
| PA_206_FAILED_RETRIEVE_UPDATED | 500 | Error retrieving attribute after update |
| PA_207_FAILED_DELETE_ATTRIBUTE | 500 | Error deleting attribute from database |
| PA_208_FAILED_UPDATE_KEY | 500 | Error updating attribute key name |

#### Audit Logging Errors (PA_301-PA_301)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| PA_301_FAILED_AUDIT_LOG | None* | Error logging request to audit trail (non-blocking) |

*Non-blocking error - operation continues if audit log fails

---

### Key-Value Endpoints (KV_*)

#### Validation Errors (KV_001-KV_003)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| KV_001_INVALID_REQUEST_BODY | 400 | Request body is malformed or invalid JSON |
| KV_002_MISSING_KEY_OR_VALUE | 400 | Required "key" or "value" field is missing |
| KV_003_MISSING_KEY_PARAM | 400 | Required "key" path parameter is empty |

#### Resource Not Found Errors (KV_101-KV_101)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| KV_101_KEY_NOT_FOUND | 404 | Specified key does not exist in database |

#### Database Operation Errors (KV_201-KV_203)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| KV_201_FAILED_SET_VALUE | 500 | Error setting or updating key-value pair in database |
| KV_202_FAILED_RETRIEVE_VALUE | 500 | Error retrieving key-value pair from database |
| KV_203_FAILED_DELETE_VALUE | 500 | Error deleting key-value pair from database |

---

### API Key Middleware (API_*)

#### Authentication Errors (API_001-API_004)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| API_001_MISSING_API_KEY | 401 | Required "x-api-key" header is missing |
| API_002_INVALID_API_KEY_FORMAT | 401 | API key does not match expected format |
| API_003_KEYS_NOT_CONFIGURED | 503 | No valid API keys configured in environment |
| API_004_INVALID_API_KEY | 401 | API key provided does not match configured keys |

---

### Health Check (HC_*)

#### Health Check Errors (HC_001-HC_001)
| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| HC_001_HEALTH_CHECK_FAILED | 500 | Database health check failed |

---

### Database Setup (DB_*)

#### Database Connection Errors (DB_001-DB_007)
| Error Code | Status | Description |
|-----------|--------|-------------|
| DB_001_URL_NOT_SET | Fatal | DATABASE_URL environment variable is not set |
| DB_002_INVALID_PORT | Fatal | PORT environment variable is not a valid integer |
| DB_003_FAILED_PARSE_URL | Fatal | Unable to parse DATABASE_URL configuration |
| DB_004_FAILED_CREATE_POOL | Fatal | Failed to create database connection pool |
| DB_005_PING_FAILED | Fatal | Database ping test failed - cannot connect |
| DB_006_FAILED_START_SERVER | Error | Server failed to start on configured port |
| DB_007_FAILED_SHUTDOWN_SERVER | Fatal | Server failed to shutdown gracefully |

---

## Implementation Details

### Updated Files

1. **[errors/errors.go](errors/errors.go)** - Central location for all error codes and ErrorResponse struct
2. **[person_attributes/person_attributes.go](person_attributes/person_attributes.go)** - All endpoints now return error codes
3. **[key_value/key_value.go](key_value/key_value.go)** - All endpoints now return error codes
4. **[middleware/api_key.go](middleware/api_key.go)** - Authentication errors now include error codes
5. **[healthcheck/health_handler.go](healthcheck/health_handler.go)** - Health check errors now include error codes
6. **[main.go](main.go)** - Database setup and shutdown errors now include error codes

### Error Code Naming Convention

Error codes follow the pattern: `PREFIX_SEQUENCE_DESCRIPTION`

- **PREFIX**: 2-letter module identifier (PA, KV, API, HC, DB)
- **SEQUENCE**: 3-digit category and sequence number
  - First digit: Category (0=validation, 1=not found, 2=database ops, 3=other)
  - Last two digits: Sequential number within category
- **DESCRIPTION**: Uppercase descriptive name

### Usage Example

When an error occurs, the client receives:

```json
{
  "message": "Person not found",
  "error_code": "PA_101_PERSON_NOT_FOUND"
}
```

The error code can be used to:
- Quickly identify the issue in logs
- Implement client-side error handling
- Correlate errors across distributed systems
- Build error dashboards and alerts

---

## Monitoring and Debugging

### Log Parsing

Search logs by error code to find all occurrences:

```bash
grep "PA_101_PERSON_NOT_FOUND" application.log
```

### Error Tracking

Each error code should be tracked in monitoring systems for:
- Frequency of occurrence
- Response time impact
- Correlation with other errors
- Trend analysis over time

### Client Integration

Clients can parse the error_code field to:
- Show localized error messages
- Retry with appropriate backoff strategies
- Route errors to appropriate handlers
- Log structured error data
