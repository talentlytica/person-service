Feature: Health Check Integration Tests
  As a system administrator
  I want to monitor the service health
  So that I can ensure the service is running properly

  Background:
    Given the service is running

  Scenario: Health endpoint returns 200 OK
    When I send a GET request to "/health"
    Then the response status should be 200
    And the response should contain "status" with value "healthy"

  Scenario: Health endpoint returns valid JSON
    When I send a GET request to "/health"
    Then the response content type should be "application/json"
    And the response should be valid JSON
    And the response should be an object

  Scenario: Health endpoint includes service metadata
    When I send a GET request to "/health"
    Then the response should have field "status"
    And the field "status" should be one of "healthy,unhealthy"

  Scenario: Service is responsive to multiple health checks
    When I send 3 concurrent GET requests to "/health"
    Then all responses should have status 200

  Scenario: Health check does not timeout
    When I send a GET request to "/health" with 5000ms timeout
    Then the response status should be 200
    And the request should complete within timeout

