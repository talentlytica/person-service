Feature: Person Attributes Management
  As a user of the Person Service API
  I want to manage attributes for persons
  So that I can store and retrieve dynamic person information

  Background:
    Given the persons and attributes table is empty
    And the service is running
    And I have a valid API key

  # Happy Path Scenarios

  Scenario: Add a single attribute to a person
    Given a person exists with the following details:
      | name      | clientId   |
      | John Doe  | 1234567890 |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key   | value                  |
      | email | john.doe@example.com   |
    And the request meta contains:
      | caller  | reason              | traceId                              |
      | user123 | add email attribute | 550e8400-e29b-41d4-a716-446655440000 |
    Then the response status should be 201
    And the response should contain an attribute with:
      | key   | value                  |
      | email | john.doe@example.com   |
    And the attribute should have an "id"
    And the attribute should have "createdAt" timestamp
    And the attribute should have "updatedAt" timestamp

  Scenario: Add multiple attributes to a person
    Given a person exists with the following details:
      | name       | clientId   |
      | Jane Smith | 9876543210 |
    When I add the following attributes to the person:
      | key     | value                    |
      | email   | jane.smith@example.com   |
      | phone   | +1234567890              |
      | address | 123 Main St, City, State |
    Then all attributes should be added successfully
    And the person should have 3 attributes

  Scenario: Get all attributes for a person
    Given a person exists with the following details:
      | name       | clientId   |
      | Bob Wilson | 5555555555 |
    And the person has the following attributes:
      | key        | value                  |
      | email      | bob.wilson@example.com |
      | department | Engineering            |
      | position   | Senior Developer       |
    When I send a GET request to "/persons/{personId}/attributes"
    Then the response status should be 200
    And the response should contain 3 attributes
    And the attributes should include:
      | key        | value                  |
      | email      | bob.wilson@example.com |
      | department | Engineering            |
      | position   | Senior Developer       |

  Scenario: Update an existing attribute
    Given a person exists with the following details:
      | name        | clientId   |
      | Alice Brown | 1111111111 |
    And the person has an attribute:
      | key   | value                   |
      | email | alice.brown@example.com |
    When I send a PUT request to "/persons/{personId}/attributes/{attributeId}" with:
      | key   | value                     |
      | email | alice.updated@example.com |
    And the request meta contains:
      | caller  | reason       | traceId                              |
      | user456 | update email | 660e8400-e29b-41d4-a716-446655440001 |
    Then the response status should be 200
    And the attribute "email" should have value "alice.updated@example.com"
    And the attribute "updatedAt" timestamp should be updated

  Scenario: Delete an attribute from a person
    Given a person exists with the following details:
      | name          | clientId   |
      | Charlie Davis | 2222222222 |
    And the person has the following attributes:
      | key   | value                     |
      | email | charlie.davis@example.com |
      | phone | +9876543210               |
    When I send a DELETE request to "/persons/{personId}/attributes/{attributeId}" for attribute "email"
    And the request meta contains:
      | caller  | reason       | traceId                              |
      | user789 | remove email | 770e8400-e29b-41d4-a716-446655440002 |
    Then the response status should be 200
    And the response should indicate success
    And the person should have 1 attribute remaining
    And the remaining attribute should be "phone"

  Scenario: Get attributes for a person with no attributes
    Given a person exists with the following details:
      | name      | clientId   |
      | Dan Evans | 3333333333 |
    And the person has no attributes
    When I send a GET request to "/persons/{personId}/attributes"
    Then the response status should be 200
    And the response should contain an empty attributes array

  # Error Scenarios

  Scenario: Attempt to add attribute to non-existent person
    When I send a POST request to "/persons/99999/attributes" with:
      | key   | value            |
      | email | test@example.com |
    And the request meta contains:
      | caller  | reason              | traceId                              |
      | user123 | add email attribute | 880e8400-e29b-41d4-a716-446655440003 |
    Then the response status should be 404
    And the error message should indicate "Person not found"

  Scenario: Attempt to add attribute with missing required fields
    Given a person exists with the following details:
      | name      | clientId   |
      | Test User | 4444444444 |
    When I send a POST request to "/persons/{personId}/attributes" with invalid data:
      | key |
      |     |
    And the request meta contains:
      | caller  | reason       | traceId                              |
      | user123 | invalid test | 990e8400-e29b-41d4-a716-446655440004 |
    Then the response status should be 400
    And the error should contain validation details

  Scenario: Attempt to update non-existent attribute
    Given a person exists with the following details:
      | name       | clientId   |
      | Test User2 | 5555555555 |
    When I send a PUT request to "/persons/{personId}/attributes/99999" with:
      | key   | value            |
      | email | test@example.com |
    And the request meta contains:
      | caller  | reason      | traceId                              |
      | user123 | update test | 101e8400-e29b-41d4-a716-446655440005 |
    Then the response status should be 404
    And the error message should indicate "Attribute not found"

  Scenario: Attempt to delete non-existent attribute
    Given a person exists with the following details:
      | name       | clientId   |
      | Test User3 | 6666666666 |
    When I send a DELETE request to "/persons/{personId}/attributes/99999"
    And the request meta contains:
      | caller  | reason      | traceId                              |
      | user123 | delete test | 111e8400-e29b-41d4-a716-446655440006 |
    Then the response status should be 404
    And the error message should indicate "Attribute not found"

  Scenario: Get attributes for non-existent person
    When I send a GET request to "/persons/99999/attributes"
    Then the response status should be 404
    And the error message should indicate "Person not found"

  Scenario: Attempt to add attribute without meta information
    Given a person exists with the following details:
      | name       | clientId   |
      | Test User4 | 7777777777 |
    When I send a POST request to "/persons/{personId}/attributes" without meta:
      | key   | value            |
      | email | test@example.com |
    Then the response status should be 400
    And the error should indicate missing required field "meta"

  # API Key Authentication Scenarios

  Scenario: Attempt to add attribute without API key
    Given a person exists with the following details:
      | name         | clientId   |
      | No Key User  | 8080808080 |
    When I send a POST request to "/persons/{personId}/attributes" without API key:
      | key   | value            |
      | email | test@example.com |
    And the request meta contains:
      | caller  | reason       | traceId                              |
      | user123 | no key test  | 221e8400-e29b-41d4-a716-446655440017 |
    Then the response status should be 401
    And the error message should indicate "Missing required header"

  Scenario: Attempt to add attribute with invalid API key format
    Given a person exists with the following details:
      | name           | clientId   |
      | Bad Key User   | 8181818181 |
    When I send a POST request to "/persons/{personId}/attributes" with invalid API key "invalid-key-format":
      | key   | value            |
      | email | test@example.com |
    And the request meta contains:
      | caller  | reason          | traceId                              |
      | user123 | bad format test | 231e8400-e29b-41d4-a716-446655440018 |
    Then the response status should be 401
    And the error message should indicate "Invalid API key format"

  Scenario: Attempt to add attribute with wrong API key
    Given a person exists with the following details:
      | name          | clientId   |
      | Wrong Key User | 8282828282 |
    When I send a POST request to "/persons/{personId}/attributes" with invalid API key "person-service-key-00000000-0000-0000-0000-000000000000":
      | key   | value            |
      | email | test@example.com |
    And the request meta contains:
      | caller  | reason         | traceId                              |
      | user123 | wrong key test | 241e8400-e29b-41d4-a716-446655440019 |
    Then the response status should be 401
    And the error message should indicate "Invalid API key"

  Scenario: Successfully use green API key
    Given a person exists with the following details:
      | name           | clientId   |
      | Green Key User | 8383838383 |
    When I send a POST request to "/persons/{personId}/attributes" using green API key with:
      | key   | value                   |
      | email | greenkey@example.com    |
    And the request meta contains:
      | caller  | reason         | traceId                              |
      | user123 | green key test | 251e8400-e29b-41d4-a716-446655440020 |
    Then the response status should be 201
    And the response should contain an attribute with:
      | key   | value                   |
      | email | greenkey@example.com    |

  # Attribute Value Update Scenarios

  Scenario: Update only the value of an attribute keeping the key same
    Given a person exists with the following details:
      | name        | clientId   |
      | Emma Wilson | 8888888888 |
    And the person has an attribute:
      | key    | value  |
      | status | active |
    When I send a PUT request to "/persons/{personId}/attributes/{attributeId}" with:
      | key    | value    |
      | status | inactive |
    And the request meta contains:
      | caller   | reason          | traceId                              |
      | admin001 | deactivate user | 121e8400-e29b-41d4-a716-446655440007 |
    Then the response status should be 200
    And the attribute "status" should have value "inactive"

  Scenario: Update attribute key and value
    Given a person exists with the following details:
      | name         | clientId   |
      | Frank Miller | 9999999999 |
    And the person has an attribute:
      | key       | value      |
      | temp_role | contractor |
    When I send a PUT request to "/persons/{personId}/attributes/{attributeId}" with:
      | key  | value     |
      | role | full_time |
    And the request meta contains:
      | caller | reason      | traceId                              |
      | hr001  | update role | 131e8400-e29b-41d4-a716-446655440008 |
    Then the response status should be 200
    And the attribute should have key "role" and value "full_time"

  # Attribute Lifecycle Scenario

  Scenario: Complete attribute lifecycle - Create, Read, Update, Delete
    Given a person exists with the following details:
      | name         | clientId   |
      | Grace Taylor | 1010101010 |
    When I add an attribute to the person:
      | key      | value    |
      | location | New York |
    Then the attribute should be created successfully
    When I retrieve all attributes for the person
    Then I should see the "location" attribute with value "New York"
    When I update the "location" attribute to:
      | key      | value       |
      | location | Los Angeles |
    Then the attribute should be updated successfully
    And the "location" attribute should have value "Los Angeles"
    When I delete the "location" attribute
    Then the attribute should be deleted successfully
    And the person should have no attributes

  # Edge Cases

  Scenario: Add attribute with special characters in value
    Given a person exists with the following details:
      | name      | clientId   |
      | Henry Lee | 1212121212 |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key | value                                      |
      | bio | Software Engineer @ Company™ <test@email> |
    And the request meta contains:
      | caller  | reason  | traceId                              |
      | user123 | add bio | 141e8400-e29b-41d4-a716-446655440009 |
    Then the response status should be 201
    And the attribute value should be stored correctly with special characters

  Scenario: Add attribute with empty string value
    Given a person exists with the following details:
      | name     | clientId   |
      | Ivy Chen | 1313131313 |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key   | value |
      | notes |       |
    And the request meta contains:
      | caller  | reason          | traceId                              |
      | user123 | add empty notes | 151e8400-e29b-41d4-a716-446655440010 |
    Then the response status should be 201
    And the attribute should be created with empty value

  Scenario: Add multiple attributes with same key to same person (upsert behavior)
    Given a person exists with the following details:
      | name       | clientId   |
      | Jack Brown | 1414141414 |
    And the person has an attribute:
      | key   | value            |
      | email | jack@example.com |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key   | value             |
      | email | jack2@example.com |
    And the request meta contains:
      | caller  | reason           | traceId                              |
      | user123 | add second email | 161e8400-e29b-41d4-a716-446655440011 |
    Then the response status should be 201
    And the person should have 1 attribute with key "email" and value "jack2@example.com"

  # Database Encryption Verification

  Scenario: Attribute value is stored encrypted in database
    Given a person exists with the following details:
      | name           | clientId   |
      | Secure User    | 1515151515 |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key      | value                 |
      | password | my-secret-password123 |
    And the request meta contains:
      | caller  | reason          | traceId                              |
      | user123 | add password    | 171e8400-e29b-41d4-a716-446655440012 |
    Then the response status should be 201
    And the raw database value for attribute "password" should not equal "my-secret-password123"
    And the raw database value should be encrypted bytes

  Scenario: Key version is stored correctly for each attribute
    Given a person exists with the following details:
      | name            | clientId   |
      | Version User    | 1616161616 |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key   | value             |
      | email | version@test.com  |
    And the request meta contains:
      | caller  | reason       | traceId                              |
      | user123 | add email    | 181e8400-e29b-41d4-a716-446655440013 |
    Then the response status should be 201
    And the database should have key_version 1 for the attribute "email"

  Scenario: Decryption works correctly on retrieval
    Given a person exists with the following details:
      | name             | clientId   |
      | Decrypt User     | 1717171717 |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key    | value                          |
      | secret | this-is-sensitive-information  |
    And the request meta contains:
      | caller  | reason     | traceId                              |
      | user123 | add secret | 191e8400-e29b-41d4-a716-446655440014 |
    Then the response status should be 201
    When I send a GET request to "/persons/{personId}/attributes"
    Then the response status should be 200
    And the attribute "secret" should have value "this-is-sensitive-information"
    And the raw database value for attribute "secret" should not equal "this-is-sensitive-information"


  # Audit Verification

  Scenario: Verify audit log creation
    Given a person exists with the following details:
      | name       | clientId   |
      | Audit User | 1818181818 |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key    | value     |
      | status | audited   |
    And the request meta contains:
      | caller  | reason      | traceId                              |
      | user123 | add audited | 201e8400-e29b-41d4-a716-446655440015 |
    Then the response status should be 201
    And an audit record should be created for traceId "201e8400-e29b-41d4-a716-446655440015"
    And the audit record should contain caller "user123" and reason "add audited"

  # Idempotency Verification

  Scenario: Idempotency of request with same traceId
    Given a person exists with the following details:
      | name            | clientId   |
      | Idempotent User | 1919191919 |
    When I send a POST request to "/persons/{personId}/attributes" with:
      | key   | value      |
      | token | unique-123 |
    And the request meta contains:
      | caller  | reason    | traceId                              |
      | user123 | add token | 211e8400-e29b-41d4-a716-446655440016 |
    Then the response status should be 201
    When I send the same POST request again with traceId "211e8400-e29b-41d4-a716-446655440016"
    Then the response status should be 201
    And the attribute should be created only once
