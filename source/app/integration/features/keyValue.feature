Feature: Key Value integration tests
  As a system administrator
  I want to test the key value feature
  So that I can ensure the key value feature is working properly

  Background:
    Given the service is running

  # ============================================
  # POST /api/key-value scenarios
  # ============================================

  Scenario: Key-Value table
    When I insert key "a key" and value "a value" directly to database
    Then it should return key "a key" and value "a value" and created_at and updated_at should be current timestamp

  Scenario: Key-Value api
    When I call the key-value api with key "api key" and value "api value"
    Then it should respond with key "api key" and value "api value" and created_at and updated_at should be current timestamp
      And exist row in key_value table with key "api key" and value "api value"

  Scenario: Create a new key-value pair
    When I send a POST request to "/api/key-value" with body:
      """
      {"key": "test-key", "value": "test-value"}
      """
    Then the response status should be 201
      And the response should contain field "key" with value "test-key"
      And the response should contain field "value" with value "test-value"
      And the response should contain field "created_at"
      And the response should contain field "updated_at"
      And the "created_at" timestamp should be a valid ISO 8601 datetime
      And the "created_at" timestamp should be equal to "updated_at"

  Scenario: Update an existing key-value pair
    Given a key-value pair exists with key "existing-key" and value "original-value"
    When I send a POST request to "/api/key-value" with body:
      """
      {"key": "existing-key", "value": "updated-value"}
      """
    Then the response status should be 200
      And the response should contain field "key" with value "existing-key"
      And the response should contain field "value" with value "updated-value"
    When I send a GET request to "/api/key-value/existing-key"
    Then the response status should be 200
      And the response should contain field "value" with value "updated-value"

  Scenario: Create key-value with special characters
    When I send a POST request to "/api/key-value" with body:
      """
      {"key": "special!@#$%^&*()key", "value": "special!@#$%^&*()value"}
      """
    Then the response status should be 201
      And the response should contain field "key" with value "special!@#$%^&*()key"
      And the response should contain field "value" with value "special!@#$%^&*()value"

  Scenario: Missing required field - key
    When I send a POST request to "/api/key-value" with body:
      """
      {"value": "test-value"}
      """
    Then the response status should be 400
      And the error message should contain "key"

  Scenario: Missing required field - value
    When I send a POST request to "/api/key-value" with body:
      """
      {"key": "test-key"}
      """
    Then the response status should be 400
      And the error message should contain "value"

  Scenario: Empty key
    When I send a POST request to "/api/key-value" with body:
      """
      {"key": "", "value": "test-value"}
      """
    Then the response status should be 400
      And the error message should contain "key"

  Scenario: Empty value
    When I send a POST request to "/api/key-value" with body:
      """
      {"key": "test-key", "value": ""}
      """
    Then the response status should be 400
      And the error message should contain "value"

  Scenario: Invalid JSON body
    When I send a POST request to "/api/key-value" with invalid JSON
    Then the response status should be 400

  # ============================================
  # GET /api/key-value/:key scenarios
  # ============================================

  Scenario: Get an existing key-value pair
    Given a key-value pair exists with key "get-test-key" and value "get-test-value"
    When I send a GET request to "/api/key-value/get-test-key"
    Then the response status should be 200
      And the response should contain field "key" with value "get-test-key"
      And the response should contain field "value" with value "get-test-value"
      And the response should contain field "created_at"
      And the response should contain field "updated_at"

  Scenario: Get non-existent key
    When I send a GET request to "/api/key-value/non-existent-key"
    Then the response status should be 404

  Scenario: Get key with special characters
    Given a key-value pair exists with key "special!@#key" and value "special-value"
    When I send a GET request to "/api/key-value/special!@#key"
    Then the response status should be 200
      And the response should contain field "value" with value "special-value"

  Scenario: Get key returns latest value after update
    Given a key-value pair exists with key "updated-key" and value "original-value"
      And the key "updated-key" is updated to value "new-value"
    When I send a GET request to "/api/key-value/updated-key"
    Then the response status should be 200
      And the response should contain field "value" with value "new-value"

  Scenario: Get multiple different keys
    Given a key-value pair exists with key "key-1" and value "value-1"
      And a key-value pair exists with key "key-2" and value "value-2"
    When I send a GET request to "/api/key-value/key-1"
    Then the response status should be 200
      And the response should contain field "value" with value "value-1"
    When I send a GET request to "/api/key-value/key-2"
    Then the response status should be 200
      And the response should contain field "value" with value "value-2"

  Scenario: Get key with long value
    Given a key-value pair exists with key "long-value-key" and value "This is a very long value that contains more than a hundred characters to test that the system can handle longer values properly without any issues."
    When I send a GET request to "/api/key-value/long-value-key"
    Then the response status should be 200
      And the response should contain field "value" with value "This is a very long value that contains more than a hundred characters to test that the system can handle longer values properly without any issues."

  Scenario: Get key after it was deleted
    Given a key-value pair exists with key "to-be-deleted" and value "some-value"
      And the key "to-be-deleted" is deleted
    When I send a GET request to "/api/key-value/to-be-deleted"
    Then the response status should be 404

  Scenario: Verify timestamps are valid
    Given a key-value pair exists with key "timestamp-key" and value "timestamp-value"
    When I send a GET request to "/api/key-value/timestamp-key"
    Then the response status should be 200
      And the "created_at" timestamp should be a valid ISO 8601 datetime
      And the "updated_at" timestamp should be a valid ISO 8601 datetime

  # ============================================
  # DELETE /api/key-value/:key scenarios
  # ============================================

  Scenario: Delete an existing key-value pair
    Given a key-value pair exists with key "delete-key" and value "delete-value"
    When I send a DELETE request to "/api/key-value/delete-key"
    Then the response status should be 200

  Scenario: Verify key is deleted from database
    Given a key-value pair exists with key "verify-delete-key" and value "verify-delete-value"
    When I send a DELETE request to "/api/key-value/verify-delete-key"
    Then the response status should be 200
      And the key "verify-delete-key" should not exist in the database

  Scenario: Delete non-existent key
    When I send a DELETE request to "/api/key-value/non-existent-key"
    Then the response status should be 404

  Scenario: Delete key with special characters
    Given a key-value pair exists with key "special!@#delete" and value "special-value"
    When I send a DELETE request to "/api/key-value/special!@#delete"
    Then the response status should be 200
      And the key "special!@#delete" should not exist in the database

  Scenario: Delete same key twice
    Given a key-value pair exists with key "double-delete" and value "some-value"
    When I send a DELETE request to "/api/key-value/double-delete"
    Then the response status should be 200
    When I send a DELETE request to "/api/key-value/double-delete"
    Then the response status should be 404

  Scenario: Delete multiple keys sequentially
    Given a key-value pair exists with key "multi-delete-1" and value "value-1"
      And a key-value pair exists with key "multi-delete-2" and value "value-2"
    When I send a DELETE request to "/api/key-value/multi-delete-1"
    Then the response status should be 200
    When I send a DELETE request to "/api/key-value/multi-delete-2"
    Then the response status should be 200
      And the key "multi-delete-1" should not exist in the database
      And the key "multi-delete-2" should not exist in the database

  Scenario: Create, Read, Update, Delete lifecycle
    # Create
    When I send a POST request to "/api/key-value" with body:
      """
      {"key": "lifecycle-key", "value": "initial-value"}
      """
    Then the response status should be 201
      And the response should contain field "key" with value "lifecycle-key"
    # Read
    When I send a GET request to "/api/key-value/lifecycle-key"
    Then the response status should be 200
      And the response should contain field "value" with value "initial-value"
    # Update
    When I send a POST request to "/api/key-value" with body:
      """
      {"key": "lifecycle-key", "value": "updated-value"}
      """
    Then the response status should be 200
      And the response should contain field "value" with value "updated-value"
    When I send a GET request to "/api/key-value/lifecycle-key"
    Then the response status should be 200
      And the response should contain field "value" with value "updated-value"
    # Delete
    When I send a DELETE request to "/api/key-value/lifecycle-key"
    Then the response status should be 200
    When I send a GET request to "/api/key-value/lifecycle-key"
    Then the response status should be 404
