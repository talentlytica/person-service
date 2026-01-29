# Passed Test Cases

**Total: 198**

| # | Test Name | File | Duration (ms) |
|---|-----------|------|---------------|
| 1 | SPEC: Should prevent SQL injection in attribute key | `API_Tests/specification_tests/person_attributes_security_spec.test.js` | 90 |
| 2 | SPEC: Should sanitize XSS payloads in values | `API_Tests/specification_tests/person_attributes_security_spec.test.js` | 58 |
| 3 | SPEC: Should handle extremely long attribute keys | `API_Tests/specification_tests/person_attributes_security_spec.test.js` | 8 |
| 4 | SPEC: Should handle extremely long attribute values | `API_Tests/specification_tests/person_attributes_security_spec.test.js` | 72 |
| 5 | SPEC: Should prevent accessing other person's attributes | `API_Tests/specification_tests/person_attributes_security_spec.test.js` | 30 |
| 6 | SPEC: Should handle duplicate attribute keys correctly | `API_Tests/specification_tests/person_attributes_security_spec.test.js` | 20 |
| 7 | SPEC: Should handle null and empty values appropriately | `API_Tests/specification_tests/person_attributes_security_spec.test.js` | 86 |
| 8 | Should reject request without API key | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 14 |
| 9 | Should reject request with invalid API key | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 10 | Should reject invalid UUID format for personId | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 11 | Should reject invalid UUID format for attributeId | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 12 | Should reject request without body | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 13 | Should reject empty body | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 14 | Should reject request without value field | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 3 |
| 15 | Should reject empty string as value | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 16 | Should reject null as value | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 5 |
| 17 | Should reject whitespace-only value | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 18 | Should reject number as value | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 19 | Should reject object as value | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 20 | Should reject array as value | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 21 | Should reject invalid JSON | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 6 |
| 22 | Should handle wrong Content-Type | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 23 | Should handle or reject extremely long value | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 6 |
| 24 | Should reject POST method on PUT endpoint | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 3 |
| 25 | Should reject PATCH method | `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 26 | Should reject request with missing key field | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 17 |
| 27 | Should reject request with missing value field | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 3 |
| 28 | Should reject completely empty request body | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 2 |
| 29 | Should reject empty string as key | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 1 |
| 30 | Should reject empty string as value | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 2 |
| 31 | Should reject null as key | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 2 |
| 32 | Should reject null as value | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 2 |
| 33 | Should reject number as key | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 4 |
| 34 | Should reject object as key | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 4 |
| 35 | Should reject array as key | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 3 |
| 36 | Should reject boolean as key | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 1 |
| 37 | Should reject invalid JSON | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 2 |
| 38 | Should reject wrong Content-Type | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 1 |
| 39 | Should handle or reject key with special characters | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 10 |
| 40 | Should handle extra fields in request body | `API_Tests/negative_tests/POST_api_key-value_negative.test.js` | 2 |
| 41 | Should reject request without API key | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 14 |
| 42 | Should reject request with invalid API key | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 43 | Should reject invalid UUID format for personId | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 3 |
| 44 | Should reject invalid UUID format for attributeId | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 45 | Should reject both IDs with invalid format | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 5 |
| 46 | Should reject empty personId | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 47 | Should reject empty attributeId | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 48 | Should handle SQL injection in personId | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 49 | Should handle SQL injection in attributeId | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 50 | Should prevent path traversal | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 51 | Should handle XSS payload in IDs | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 52 | Should reject POST method on DELETE endpoint | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 53 | Should reject PATCH method | `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 54 | Should reject request without API key | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 14 |
| 55 | Should reject request with invalid API key | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 56 | Should reject invalid UUID format for personId | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 3 |
| 57 | Should reject invalid UUID format for attributeId | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 58 | Should reject both IDs with invalid format | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 59 | Should reject empty personId | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 60 | Should reject empty attributeId | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 61 | Should reject POST method on GET endpoint | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 62 | Should reject PATCH method | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 63 | Should handle SQL injection in personId | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 64 | Should handle SQL injection in attributeId | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 1 |
| 65 | Should prevent path traversal | `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js` | 2 |
| 66 | Should reject request without API key | `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js` | 13 |
| 67 | Should reject request with invalid API key | `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js` | 2 |
| 68 | Should reject invalid UUID format for personId | `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js` | 2 |
| 69 | Should reject empty personId | `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js` | 2 |
| 70 | Should reject personId with special characters | `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js` | 1 |
| 71 | Should reject DELETE method on GET endpoint | `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js` | 1 |
| 72 | Should reject PATCH method | `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js` | 2 |
| 73 | Should prevent path traversal in personId | `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js` | 1 |
| 74 | Should return 404 for non-existent key | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 14 |
| 75 | Should return 404 for random UUID as key | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 3 |
| 76 | Should reject empty key parameter | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 2 |
| 77 | Should handle URL-encoded special characters in key | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 1 |
| 78 | Should handle key with slashes (path traversal attempt) | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 1 |
| 79 | Should reject PUT method on GET endpoint | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 1 |
| 80 | Should reject POST method on GET endpoint | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 1 |
| 81 | Should handle invalid Accept header | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 1 |
| 82 | Should reject request with body (GET should not have body) | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 2 |
| 83 | Should ignore query parameters on GET | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 1 |
| 84 | Should handle or reject key with SQL injection pattern | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 1 |
| 85 | Should handle key with null byte | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 5 |
| 86 | Should handle case sensitivity in key lookup | `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js` | 2 |
| 87 | SPEC: Should create audit trail for all operations | `API_Tests/specification_tests/person_attributes_business_spec.test.js` | 33 |
| 88 | SPEC: Should handle concurrent updates gracefully | `API_Tests/specification_tests/person_attributes_business_spec.test.js` | 46 |
| 89 | SPEC: Should cascade delete attributes when person is deleted | `API_Tests/specification_tests/person_attributes_business_spec.test.js` | 26 |
| 90 | SPEC: Should maintain consistent timestamps | `API_Tests/specification_tests/person_attributes_business_spec.test.js` | 1032 |
| 91 | SPEC: Should maintain data integrity on failures | `API_Tests/specification_tests/person_attributes_business_spec.test.js` | 11 |
| 92 | SPEC: Should enforce attribute key uniqueness per person | `API_Tests/specification_tests/person_attributes_business_spec.test.js` | 19 |
| 93 | 1. CREATE attribute via API and verify ENCRYPTION in database | `API_Tests/comprehensive_tests/person_attributes_full_verification.test.js` | 35 |
| 94 | 2. GET attribute via API and verify DECRYPTION | `API_Tests/comprehensive_tests/person_attributes_full_verification.test.js` | 20 |
| 95 | 3. UPDATE attribute and verify RE-ENCRYPTION in database | `API_Tests/comprehensive_tests/person_attributes_full_verification.test.js` | 20 |
| 96 | 4. DELETE attribute and verify complete removal from database | `API_Tests/comprehensive_tests/person_attributes_full_verification.test.js` | 25 |
| 97 | 5. FULL CRUD LIFECYCLE with encryption at each step | `API_Tests/comprehensive_tests/person_attributes_full_verification.test.js` | 30 |
| 98 | Should reject request without API key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 16 |
| 99 | Should reject request with invalid API key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 3 |
| 100 | Should reject request with expired/wrong format API key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 101 | Should reject request with non-existent personId | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 1 |
| 102 | Should reject request with invalid UUID format for personId | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 103 | Should reject request with empty personId | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 1 |
| 104 | Should reject request without key field | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 105 | Should reject request without meta field | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 106 | Should reject request with empty meta object | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 107 | Should reject request without entire body | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 1 |
| 108 | Should reject empty string as key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 1 |
| 109 | Should reject null as key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 1 |
| 110 | Should reject whitespace-only key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 111 | Should reject number as key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 1 |
| 112 | Should reject object as key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 113 | Should reject array as value | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 1 |
| 114 | Should reject invalid JSON | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 115 | Should reject wrong Content-Type | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 116 | Should handle or reject extremely long key | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 2 |
| 117 | Should handle or reject extremely long value | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 89 |
| 118 | Should reject GET method on POST endpoint | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 19 |
| 119 | Should reject PATCH method | `API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js` | 1 |
| 120 | 1. CREATE key-value via API and verify in database | `API_Tests/comprehensive_tests/key_value_full_verification.test.js` | 31 |
| 121 | 2. GET key-value via API and verify response matches database | `API_Tests/comprehensive_tests/key_value_full_verification.test.js` | 15 |
| 122 | 3. UPDATE key-value via API and verify in database | `API_Tests/comprehensive_tests/key_value_full_verification.test.js` | 12 |
| 123 | 4. DELETE key-value via API and verify removal from database | `API_Tests/comprehensive_tests/key_value_full_verification.test.js` | 20 |
| 124 | 5. FULL CRUD LIFECYCLE with step-by-step database verification | `API_Tests/comprehensive_tests/key_value_full_verification.test.js` | 25 |
| 125 | Delete an existing attribute | `API_Tests/tests/DELETE_persons_personId_attributes_attributeId.test.js` | 18 |
| 126 | Verify attribute is deleted from database | `API_Tests/tests/DELETE_persons_personId_attributes_attributeId.test.js` | 12 |
| 127 | Delete non-existent attribute | `API_Tests/tests/DELETE_persons_personId_attributes_attributeId.test.js` | 3 |
| 128 | Delete without API key | `API_Tests/tests/DELETE_persons_personId_attributes_attributeId.test.js` | 6 |
| 129 | Full CRUD lifecycle | `API_Tests/tests/DELETE_persons_personId_attributes_attributeId.test.js` | 20 |
| 130 | Delete an existing key-value pair | `API_Tests/tests/DELETE_api_key-value_key.test.js` | 15 |
| 131 | Verify key is deleted from database | `API_Tests/tests/DELETE_api_key-value_key.test.js` | 9 |
| 132 | Delete non-existent key | `API_Tests/tests/DELETE_api_key-value_key.test.js` | 2 |
| 133 | Delete same key twice | `API_Tests/tests/DELETE_api_key-value_key.test.js` | 5 |
| 134 | Delete multiple keys sequentially | `API_Tests/tests/DELETE_api_key-value_key.test.js` | 10 |
| 135 | CRUD lifecycle | `API_Tests/tests/DELETE_api_key-value_key.test.js` | 9 |
| 136 | Create/update attribute using PUT | `API_Tests/tests/PUT_persons_personId_attributes.test.js` | 18 |
| 137 | Update existing attribute using PUT | `API_Tests/tests/PUT_persons_personId_attributes.test.js` | 11 |
| 138 | PUT creates new attribute if not exists | `API_Tests/tests/PUT_persons_personId_attributes.test.js` | 6 |
| 139 | PUT without API key | `API_Tests/tests/PUT_persons_personId_attributes.test.js` | 6 |
| 140 | PUT without meta | `API_Tests/tests/PUT_persons_personId_attributes.test.js` | 3 |
| 141 | PUT with empty key | `API_Tests/tests/PUT_persons_personId_attributes.test.js` | 1 |
| 142 | Update attribute value only | `API_Tests/tests/PUT_persons_personId_attributes_attributeId.test.js` | 20 |
| 143 | Update non-existent attribute | `API_Tests/tests/PUT_persons_personId_attributes_attributeId.test.js` | 5 |
| 144 | Update without API key | `API_Tests/tests/PUT_persons_personId_attributes_attributeId.test.js` | 2 |
| 145 | Update without meta information (meta is optional for UPDATE) | `API_Tests/tests/PUT_persons_personId_attributes_attributeId.test.js` | 6 |
| 146 | Create a new key-value pair | `API_Tests/tests/POST_api_key-value.test.js` | 13 |
| 147 | Update an existing key-value pair | `API_Tests/tests/POST_api_key-value.test.js` | 4 |
| 148 | Create key-value with special characters | `API_Tests/tests/POST_api_key-value.test.js` | 2 |
| 149 | Missing required field - key | `API_Tests/tests/POST_api_key-value.test.js` | 6 |
| 150 | Missing required field - value | `API_Tests/tests/POST_api_key-value.test.js` | 2 |
| 151 | Empty key | `API_Tests/tests/POST_api_key-value.test.js` | 2 |
| 152 | Empty value | `API_Tests/tests/POST_api_key-value.test.js` | 2 |
| 153 | Create multiple key-value pairs | `API_Tests/tests/POST_api_key-value.test.js` | 5 |
| 154 | Get all attributes for person with multiple attributes | `API_Tests/tests/GET_persons_personId_attributes.test.js` | 24 |
| 155 | Get attributes without API key | `API_Tests/tests/GET_persons_personId_attributes.test.js` | 5 |
| 156 | Verify all attributes are decrypted | `API_Tests/tests/GET_persons_personId_attributes.test.js` | 5 |
| 157 | Get an existing key-value pair | `API_Tests/tests/GET_api_key-value_key.test.js` | 15 |
| 158 | Get non-existent key | `API_Tests/tests/GET_api_key-value_key.test.js` | 5 |
| 159 | Get key with special characters | `API_Tests/tests/GET_api_key-value_key.test.js` | 5 |
| 160 | Get key returns latest value after update | `API_Tests/tests/GET_api_key-value_key.test.js` | 4 |
| 161 | Get multiple different keys | `API_Tests/tests/GET_api_key-value_key.test.js` | 6 |
| 162 | Verify timestamps are valid | `API_Tests/tests/GET_api_key-value_key.test.js` | 3 |
| 163 | Get existing attribute by ID | `API_Tests/tests/GET_persons_personId_attributes_attributeId.test.js` | 4 |
| 164 | Get non-existent attribute | `API_Tests/tests/GET_persons_personId_attributes_attributeId.test.js` | 7 |
| 165 | Get attribute without API key | `API_Tests/tests/GET_persons_personId_attributes_attributeId.test.js` | 2 |
| 166 | Verify attribute value is decrypted | `API_Tests/tests/GET_persons_personId_attributes_attributeId.test.js` | 3 |
| 167 | Verify response includes all fields | `API_Tests/tests/GET_persons_personId_attributes_attributeId.test.js` | 3 |
| 168 | Create a single attribute | `API_Tests/tests/POST_persons_personId_attributes.test.js` | 14 |
| 169 | Create attribute without API key | `API_Tests/tests/POST_persons_personId_attributes.test.js` | 5 |
| 170 | Create attribute with invalid API key format | `API_Tests/tests/POST_persons_personId_attributes.test.js` | 2 |
| 171 | Create attribute without meta information | `API_Tests/tests/POST_persons_personId_attributes.test.js` | 2 |
| 172 | Create attribute without key | `API_Tests/tests/POST_persons_personId_attributes.test.js` | 1 |
| 173 | Should handle DELETE of non-existent key | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 11 |
| 174 | Should handle DELETE of already deleted key | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 5 |
| 175 | Should reject DELETE with empty key | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 5 |
| 176 | Should reject DELETE with whitespace-only key | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 2 |
| 177 | Should reject GET on DELETE endpoint | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 1 |
| 178 | Should reject PATCH method | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 1 |
| 179 | Should ignore request body on DELETE | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 1 |
| 180 | Should handle DELETE with SQL injection pattern in key | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 3 |
| 181 | Should handle DELETE with path traversal attempt | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 1 |
| 182 | Should handle DELETE with XSS payload in key | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 1 |
| 183 | Should handle multiple concurrent DELETEs of same key | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 5 |
| 184 | Should ignore invalid Content-Type header on DELETE | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 1 |
| 185 | Should handle DELETE without Accept header | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 1 |
| 186 | Should handle URL-encoded key in DELETE | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 1 |
| 187 | Should handle special URL characters | `API_Tests/negative_tests/DELETE_api_key-value_key_negative.test.js` | 2 |
| 188 | Health endpoint returns 200 OK | `API_Tests/tests/GET_health.test.js` | 12 |
| 189 | Health response contains status field | `API_Tests/tests/GET_health.test.js` | 2 |
| 190 | Health response is valid JSON | `API_Tests/tests/GET_health.test.js` | 1 |
| 191 | Health check responds within acceptable time | `API_Tests/tests/GET_health.test.js` | 3 |
| 192 | Multiple health checks are consistent | `API_Tests/tests/GET_health.test.js` | 8 |
| 193 | Health check does not timeout | `API_Tests/tests/GET_health.test.js` | 2 |
| 194 | Should handle invalid HTTP method (POST instead of GET) | `API_Tests/negative_tests/GET_health_negative.test.js` | 13 |
| 195 | Should handle invalid path (/healths instead of /health) | `API_Tests/negative_tests/GET_health_negative.test.js` | 2 |
| 196 | Should handle path with query parameters | `API_Tests/negative_tests/GET_health_negative.test.js` | 2 |
| 197 | Should handle excessive request timeout | `API_Tests/negative_tests/GET_health_negative.test.js` | 1 |
| 198 | Should handle malformed Accept header | `API_Tests/negative_tests/GET_health_negative.test.js` | 1 |
