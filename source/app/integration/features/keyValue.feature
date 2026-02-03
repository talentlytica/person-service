Feature: Key Value integration tests
  As a system administrator
  I want to test the key value feature
  So that I can ensure the key value feature is working properly

  Background:
    Given the service is running

  Scenario: Key-Value table
    When I insert key "a key" and value "a value" directly to database
    Then it should return key "a key" and value "a value" and created_at and updated_at should be current timestamp

  Scenario: Key-Value api
    When I call the key-value api with key "api key" and value "api value"
    Then it should respond with key "api key" and value "api value" and created_at and updated_at should be current timestamp
      And exist row in key_value table with key "api key" and value "api value"
    