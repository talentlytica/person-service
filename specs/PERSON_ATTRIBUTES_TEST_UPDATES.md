# Person Attributes Test Updates Summary

## Overview
Updated `specs/steps/person_attributes.steps.js` to match the actual database schema defined in `source/app/internal/db/schema.sql` and queries in `source/app/internal/db/queries.sql`.

## Key Changes Made

### 1. **Table Name Correction**
- **Before:** `persons` (plural)
- **After:** `person` (singular)
- **Impact:** All database queries now use the correct table name

### 2. **Person Table Schema Alignment**
- **Before:** Expected columns: `id`, `name`, `client_id`, `created_at`, `updated_at`
- **After:** Actual columns: `id` (UUID), `client_id`, `created_at`, `updated_at`, `deleted_at`
- **Changes:**
  - Removed `name` column handling (not in schema)
  - Person IDs are now UUIDs (generated via `uuidv7()`)
  - Added support for `deleted_at` (soft delete)

### 3. **Person Attributes Table Schema Alignment**
- **Before:** Expected columns: `id`, `person_id`, `key`, `value`, `meta`
- **After:** Actual columns: `id`, `person_id`, `attribute_key`, `encrypted_value`, `key_version`
- **Changes:**
  - `key` → `attribute_key`
  - `value` → `encrypted_value` (stored as encrypted BYTEA)
  - Removed `meta` column (not in schema)
  - Added `key_version` for encryption key rotation

### 4. **Encryption Implementation**
Added PostgreSQL `pgcrypto` extension support for encrypting/decrypting attribute values:

```javascript
// Encryption key configuration
const ENCRYPTION_KEY = process.env.ENCRYPTION_KEY || 'test-encryption-key-12345';
const KEY_VERSION = 1; // Integer representing encryption key version
```

**Helper Functions Updated:**

#### `createAttribute()` - Now encrypts data
```sql
INSERT INTO person_attributes (
  person_id, 
  attribute_key, 
  encrypted_value,  -- Using pgp_sym_encrypt()
  key_version
)
VALUES ($1, $2, pgp_sym_encrypt($3, $4), $5)
RETURNING id, person_id, attribute_key, 
          pgp_sym_decrypt(encrypted_value, $4) AS attribute_value,
          key_version, created_at, updated_at
```

#### `getPersonAttributes()` - Now decrypts data
```sql
SELECT 
  id,
  person_id,
  attribute_key,
  pgp_sym_decrypt(encrypted_value, $2) AS attribute_value,  -- Decrypting
  key_version,
  created_at,
  updated_at
FROM person_attributes 
WHERE person_id = $1 
ORDER BY attribute_key
```

### 5. **Unique Constraint Handling**
The schema has `UNIQUE(person_id, attribute_key)` constraint, which means:
- Cannot have duplicate attributes with the same key for one person
- `CreateOrUpdatePersonAttribute` query uses `ON CONFLICT DO UPDATE`
- Adding an attribute with existing key **updates** it instead of creating duplicate

**Test Updated:**
- "Add multiple attributes with same key to same person" scenario
- Now expects 1 attribute (updated) instead of 2 attributes (duplicated)

### 6. **Column Name Updates in Tests**
Updated all test assertions to use correct column names:
- `attr.key` → `attr.attribute_key`
- `attr.value` → `attr.attribute_value` (when reading from DB)
- API responses may still use `key` and `value` (mapped by API layer)

## Environment Variables

### New Environment Variable
- `ENCRYPTION_KEY`: Encryption key for attribute values (defaults to `test-encryption-key-12345`)

## Database Dependencies

### Required PostgreSQL Extensions
1. **pgcrypto** - For encryption/decryption functions
   - `pgp_sym_encrypt()` - Encrypts data
   - `pgp_sym_decrypt()` - Decrypts data

2. **citext** - Case-insensitive text for attribute keys

3. **uuidv7()** - UUID v7 generation (PostgreSQL 17+)

## Migration Notes

### Database Setup Order
1. Enable extensions (`pgcrypto`, `citext`)
2. Create tables with proper constraints
3. Set up indexes for performance

### Test Data Changes
- Feature file still references `name` in test tables (will be ignored by `createPerson()`)
- Only `client_id` is actually stored in the database
- Person IDs are UUIDs, not integers

## Testing Considerations

### Encryption in Tests
- All attribute values are encrypted before storage
- Decryption happens automatically in helper queries
- Encryption key is configurable via environment variable
- Key version tracks which key was used (supports key rotation)

### Performance Notes
- Encryption/decryption adds processing overhead
- Large-scale searches on encrypted values are expensive
- Consider caching or indexing strategies for production

## API vs Database Layer

### Important Distinction
- **Database Layer:** Uses `attribute_key` and `encrypted_value`
- **API Layer:** May use `key` and `value` in JSON responses
- Test helpers abstract this difference appropriately

## Files Modified

1. `specs/steps/person_attributes.steps.js` - Complete rewrite of database interaction layer

## Files That May Need Updates

1. `specs/features/person_attributes.feature` - Feature file expects behaviors that conflict with schema:
   - Last scenario expects 2 attributes with same key (impossible due to UNIQUE constraint)
   - Test data includes `name` field (not in person table)

## Backward Compatibility

### Breaking Changes
- Tests now require PostgreSQL with `pgcrypto` extension
- Encryption key must be available (via env var or default)
- Person IDs are UUIDs (may affect ID format expectations)
- Cannot create duplicate attributes with same key

## Next Steps

### Recommended Actions
1. Update feature file to match actual schema behavior
2. Configure encryption key in test environment
3. Verify all tests pass with new schema
4. Document encryption key management for production
5. Consider adding tests for key rotation scenarios

### Optional Enhancements
1. Add tests for encryption key rotation
2. Add tests for soft delete functionality (`deleted_at`)
3. Add performance tests for encrypted attribute queries
4. Add tests for UUID person ID format validation

## Verification Checklist

- [x] Table names corrected (`persons` → `person`)
- [x] Column names aligned with schema
- [x] Encryption/decryption implemented
- [x] Unique constraint behavior handled
- [x] Helper functions updated
- [x] No linter errors
- [ ] All tests passing (run test suite to verify)
- [ ] Feature file aligned with schema
- [ ] Documentation updated

## Additional Notes

### Security Considerations
- Encryption key should be stored securely (not in code)
- Use environment variables or secret management
- Consider key rotation strategy
- Audit logging for sensitive operations

### Performance Considerations
- Encrypted values cannot use standard indexes
- Consider using functional indexes if needed
- Monitor query performance on encrypted columns
- Cache decrypted values when appropriate

