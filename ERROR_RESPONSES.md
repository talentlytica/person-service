# API Error Response Examples

This document provides real-world examples of error responses for each endpoint in the person-service application.

## Person Attributes API

### POST /persons/:personId/attributes

#### Missing x-api-key header
**Status:** 401
```json
{
  "message": "Missing required header \"x-api-key\"",
  "error_code": "API_001_MISSING_API_KEY"
}
```

#### Invalid person ID format
**Status:** 404
```json
{
  "message": "Person not found",
  "error_code": "PA_001_INVALID_PERSON_ID"
}
```

#### Invalid request body
**Status:** 400
```json
{
  "message": "Invalid request body",
  "error_code": "PA_003_INVALID_REQUEST_BODY"
}
```

#### Missing required key field
**Status:** 400
```json
{
  "message": "Key is required",
  "error_code": "PA_004_MISSING_KEY"
}
```

#### Missing required meta field
**Status:** 400
```json
{
  "message": "Missing required field \"meta\"",
  "error_code": "PA_005_MISSING_META"
}
```

#### Person not found in database
**Status:** 404
```json
{
  "message": "Person not found",
  "error_code": "PA_101_PERSON_NOT_FOUND"
}
```

#### Database error during attribute creation
**Status:** 500
```json
{
  "message": "Failed to create attribute",
  "error_code": "PA_202_FAILED_CREATE_ATTRIBUTE"
}
```

#### Success Response
**Status:** 201
```json
{
  "id": 42,
  "key": "email",
  "value": "john@example.com",
  "createdAt": "2026-02-02T14:30:00Z",
  "updatedAt": "2026-02-02T14:30:00Z"
}
```

---

### GET /persons/:personId/attributes/:attributeId

#### Invalid attribute ID format
**Status:** 400
```json
{
  "message": "Invalid attribute ID format",
  "error_code": "PA_006_INVALID_ATTRIBUTE_ID_FORMAT"
}
```

#### Attribute not found
**Status:** 404
```json
{
  "message": "Attribute not found",
  "error_code": "PA_102_ATTRIBUTE_NOT_FOUND"
}
```

#### Success Response
**Status:** 200
```json
{
  "id": 42,
  "key": "email",
  "value": "john@example.com",
  "createdAt": "2026-02-02T14:30:00Z",
  "updatedAt": "2026-02-02T14:30:00Z"
}
```

---

### PUT /persons/:personId/attributes/:attributeId

#### Failed to update attribute key
**Status:** 500
```json
{
  "message": "Failed to update attribute key",
  "error_code": "PA_208_FAILED_UPDATE_KEY"
}
```

#### Failed to retrieve updated attribute
**Status:** 500
```json
{
  "message": "Failed to retrieve updated attribute",
  "error_code": "PA_206_FAILED_RETRIEVE_UPDATED"
}
```

#### Success Response
**Status:** 200
```json
{
  "id": 42,
  "key": "phone",
  "value": "+1-555-0123",
  "createdAt": "2026-02-02T14:30:00Z",
  "updatedAt": "2026-02-02T14:35:00Z"
}
```

---

### DELETE /persons/:personId/attributes/:attributeId

#### Failed to delete attribute
**Status:** 500
```json
{
  "message": "Failed to delete attribute",
  "error_code": "PA_207_FAILED_DELETE_ATTRIBUTE"
}
```

#### Success Response
**Status:** 200
```json
{
  "message": "Attribute deleted successfully"
}
```

---

## Key-Value API

### POST /api/key-value

#### Invalid request body
**Status:** 400
```json
{
  "message": "Invalid request body",
  "error_code": "KV_001_INVALID_REQUEST_BODY"
}
```

#### Missing key or value
**Status:** 400
```json
{
  "message": "Key and value are required",
  "error_code": "KV_002_MISSING_KEY_OR_VALUE"
}
```

#### Database error
**Status:** 500
```json
{
  "message": "Failed to set value",
  "error_code": "KV_201_FAILED_SET_VALUE"
}
```

#### Success Response
**Status:** 200
```json
{
  "key": "config_version",
  "value": "2.0.1",
  "created_at": "2026-02-02T14:30:00Z",
  "updated_at": "2026-02-02T14:30:00Z"
}
```

---

### GET /api/key-value/:key

#### Missing key parameter
**Status:** 400
```json
{
  "message": "Key parameter is required",
  "error_code": "KV_003_MISSING_KEY_PARAM"
}
```

#### Key not found
**Status:** 404
```json
{
  "message": "Key not found",
  "error_code": "KV_101_KEY_NOT_FOUND"
}
```

#### Database error
**Status:** 500
```json
{
  "message": "Failed to retrieve value",
  "error_code": "KV_202_FAILED_RETRIEVE_VALUE"
}
```

#### Success Response
**Status:** 200
```json
{
  "key": "config_version",
  "value": "2.0.1",
  "created_at": "2026-02-02T14:30:00Z",
  "updated_at": "2026-02-02T14:30:00Z"
}
```

---

### DELETE /api/key-value/:key

#### Database error
**Status:** 500
```json
{
  "message": "Failed to delete value",
  "error_code": "KV_203_FAILED_DELETE_VALUE"
}
```

#### Success Response
**Status:** 200
```json
{
  "message": "Key deleted successfully"
}
```

---

## Health Check API

### GET /health

#### Database health check failed
**Status:** 500
```json
{
  "message": "pq: connection refused",
  "error_code": "HC_001_HEALTH_CHECK_FAILED"
}
```

#### Success Response
**Status:** 200
```json
{
  "status": "healthy"
}
```

---

## API Key Authentication Errors

### Invalid API key format
**Status:** 401
```json
{
  "message": "Invalid API key format",
  "error_code": "API_002_INVALID_API_KEY_FORMAT"
}
```

### Keys not configured
**Status:** 503
```json
{
  "message": "API keys are not properly configured",
  "error_code": "API_003_KEYS_NOT_CONFIGURED"
}
```

### Invalid API key (doesn't match configured keys)
**Status:** 401
```json
{
  "message": "Invalid API key",
  "error_code": "API_004_INVALID_API_KEY"
}
```

---

## Debugging Error Responses

### How to Parse Error Codes in Clients

**JavaScript/TypeScript Example:**
```javascript
const response = await fetch('/persons/123/attributes/1');
const data = await response.json();

switch (data.error_code) {
  case 'PA_001_INVALID_PERSON_ID':
    console.error('The person ID format is invalid');
    break;
  case 'PA_101_PERSON_NOT_FOUND':
    console.error('The person does not exist');
    break;
  case 'API_001_MISSING_API_KEY':
    console.error('API key header is missing');
    break;
  default:
    console.error(`Unknown error: ${data.message}`);
}
```

**Python Example:**
```python
import requests

response = requests.get(
    '/persons/123/attributes/1',
    headers={'x-api-key': 'person-service-key-...'}
)

if response.status_code >= 400:
    error_data = response.json()
    error_code = error_data.get('error_code')
    error_message = error_data.get('message')
    
    logger.error(f"API Error [{error_code}]: {error_message}")
```

### Monitoring and Alerting

Set up alerts for specific error codes:

- **PA_101_PERSON_NOT_FOUND**: May indicate invalid client requests or data issues
- **KV_202_FAILED_RETRIEVE_VALUE**: May indicate database connectivity issues
- **HC_001_HEALTH_CHECK_FAILED**: Critical - service is unhealthy
- **DB_005_PING_FAILED**: Critical - database is unavailable
