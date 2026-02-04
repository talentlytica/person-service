package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"person-service/integration/testutil"

	"github.com/cucumber/godog"
)

func registerKeyValueSteps(sc *godog.ScenarioContext, tc *TestContext) {
	// Variables to store timestamps for comparison
	var savedTimestamps = make(map[string]time.Time)

	sc.Step(`^I insert key "([^"]*)" and value "([^"]*)" directly to database$`, func(key, value string) error {
		return testutil.InsertKeyValueDirect(context.Background(), tc.Pool, key, value)
	})

	sc.Step(`^it should return key "([^"]*)" and value "([^"]*)" and created_at and updated_at should be current timestamp$`, func(key, value string) error {
		// Query the database directly to verify
		storedValue, err := testutil.GetKeyValueDirect(context.Background(), tc.Pool, key)
		if err != nil {
			return err
		}
		if storedValue != value {
			return fmt.Errorf("expected value %q but got %q", value, storedValue)
		}

		// Verify timestamps via direct query
		var createdAt, updatedAt time.Time
		err = tc.Pool.QueryRow(context.Background(), `
			SELECT created_at, updated_at FROM key_value WHERE key = $1
		`, key).Scan(&createdAt, &updatedAt)
		if err != nil {
			return fmt.Errorf("failed to get timestamps: %w", err)
		}

		// Timestamps should be within the last minute
		now := time.Now()
		if now.Sub(createdAt) > time.Minute {
			return fmt.Errorf("created_at is too old: %v", createdAt)
		}
		if now.Sub(updatedAt) > time.Minute {
			return fmt.Errorf("updated_at is too old: %v", updatedAt)
		}

		return nil
	})

	sc.Step(`^I call the key-value api with key "([^"]*)" and value "([^"]*)"$`, func(key, value string) error {
		body := map[string]string{
			"key":   key,
			"value": value,
		}
		tc.Response = tc.Server.POST("/api/key-value", body, nil)
		return nil
	})

	sc.Step(`^it should respond with key "([^"]*)" and value "([^"]*)" and created_at and updated_at should be current timestamp$`, func(key, value string) error {
		// Accept 200 (update) or 201 (create) as success
		if tc.Response.Code != 200 && tc.Response.Code != 201 {
			return fmt.Errorf("expected status 200 or 201 but got %d: %s", tc.Response.Code, tc.Response.Body.String())
		}

		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}

		if result["key"] != key {
			return fmt.Errorf("expected key %q but got %q", key, result["key"])
		}
		if result["value"] != value {
			return fmt.Errorf("expected value %q but got %q", value, result["value"])
		}
		if result["created_at"] == nil {
			return fmt.Errorf("created_at is missing")
		}
		if result["updated_at"] == nil {
			return fmt.Errorf("updated_at is missing")
		}

		return nil
	})

	sc.Step(`^exist row in key_value table with key "([^"]*)" and value "([^"]*)"$`, func(key, value string) error {
		storedValue, err := testutil.GetKeyValueDirect(context.Background(), tc.Pool, key)
		if err != nil {
			return err
		}
		if storedValue != value {
			return fmt.Errorf("expected value %q but got %q", value, storedValue)
		}
		return nil
	})

	// Setup steps
	sc.Step(`^a key-value pair exists with key "([^"]*)" and value "([^"]*)"$`, func(key, value string) error {
		return testutil.InsertKeyValueDirect(context.Background(), tc.Pool, key, value)
	})

	sc.Step(`^the key "([^"]*)" is updated to value "([^"]*)"$`, func(key, value string) error {
		return testutil.InsertKeyValueDirect(context.Background(), tc.Pool, key, value)
	})

	sc.Step(`^the key "([^"]*)" is deleted$`, func(key string) error {
		return testutil.DeleteKeyValueDirect(context.Background(), tc.Pool, key)
	})

	sc.Step(`^the key-value table exists$`, func() error {
		// Table existence is guaranteed by migrations in Before hook
		return nil
	})

	// Request steps - specific to key-value API paths
	sc.Step(`^I send a GET request to "/api/key-value/([^"]*)"$`, func(key string) error {
		tc.Response = tc.Server.GET("/api/key-value/"+key, nil)
		return nil
	})

	sc.Step(`^I send a DELETE request to "/api/key-value/([^"]*)"$`, func(key string) error {
		tc.Response = tc.Server.DELETE("/api/key-value/"+key, nil)
		return nil
	})

	sc.Step(`^I send a POST request to "([^"]*)" with body:$`, func(path string, body *godog.DocString) error {
		var jsonBody map[string]interface{}
		if err := json.Unmarshal([]byte(body.Content), &jsonBody); err != nil {
			return fmt.Errorf("invalid JSON in docstring: %w", err)
		}
		tc.Response = tc.Server.POST(path, jsonBody, nil)
		return nil
	})

	sc.Step(`^I send a POST request to "/api/key-value" with invalid JSON$`, func() error {
		tc.Response = tc.Server.POSTRaw("/api/key-value", "{invalid json", nil)
		return nil
	})

	// Assertion steps
	sc.Step(`^the response should contain field "([^"]*)"$`, func(field string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		if _, ok := result[field]; !ok {
			return fmt.Errorf("field %q not found in response: %v", field, result)
		}
		return nil
	})

	sc.Step(`^the response should contain field "([^"]*)" with value "([^"]*)"$`, func(field, expectedValue string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		actual, ok := result[field]
		if !ok {
			return fmt.Errorf("field %q not found in response: %v", field, result)
		}
		actualStr := fmt.Sprintf("%v", actual)
		if actualStr != expectedValue {
			return fmt.Errorf("expected %q to be %q but got %q", field, expectedValue, actualStr)
		}
		return nil
	})

	sc.Step(`^I save the "([^"]*)" timestamp$`, func(field string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		timestampStr, ok := result[field].(string)
		if !ok {
			return fmt.Errorf("field %q not found or not a string", field)
		}
		t, err := time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			t, err = time.Parse(time.RFC3339, timestampStr)
			if err != nil {
				return fmt.Errorf("failed to parse timestamp %q: %w", timestampStr, err)
			}
		}
		savedTimestamps[field] = t
		return nil
	})

	sc.Step(`^the "([^"]*)" timestamp should be more recent than the saved "([^"]*)"$`, func(field, savedField string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		timestampStr, ok := result[field].(string)
		if !ok {
			return fmt.Errorf("field %q not found or not a string", field)
		}
		t, err := time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			t, err = time.Parse(time.RFC3339, timestampStr)
			if err != nil {
				return fmt.Errorf("failed to parse timestamp %q: %w", timestampStr, err)
			}
		}
		saved, ok := savedTimestamps[savedField]
		if !ok {
			return fmt.Errorf("no saved timestamp for %q", savedField)
		}
		if !t.After(saved) {
			return fmt.Errorf("expected %q (%v) to be more recent than saved %q (%v)", field, t, savedField, saved)
		}
		return nil
	})

	sc.Step(`^the "([^"]*)" timestamp should be a valid ISO 8601 datetime$`, func(field string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		timestampStr, ok := result[field].(string)
		if !ok {
			return fmt.Errorf("field %q not found or not a string", field)
		}
		// Try parsing as RFC3339 (ISO 8601 compatible)
		_, err := time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			_, err = time.Parse(time.RFC3339, timestampStr)
			if err != nil {
				return fmt.Errorf("field %q value %q is not a valid ISO 8601 datetime", field, timestampStr)
			}
		}
		return nil
	})

	sc.Step(`^the "([^"]*)" timestamp should be equal to "([^"]*)"$`, func(field1, field2 string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		ts1, ok := result[field1].(string)
		if !ok {
			return fmt.Errorf("field %q not found or not a string", field1)
		}
		ts2, ok := result[field2].(string)
		if !ok {
			return fmt.Errorf("field %q not found or not a string", field2)
		}
		if ts1 != ts2 {
			return fmt.Errorf("expected %q (%q) to equal %q (%q)", field1, ts1, field2, ts2)
		}
		return nil
	})

	sc.Step(`^the key "([^"]*)" should not exist in the database$`, func(key string) error {
		exists, err := testutil.KeyValueExists(context.Background(), tc.Pool, key)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("expected key %q to not exist, but it does", key)
		}
		return nil
	})

	sc.Step(`^the key "([^"]*)" should exist in the database$`, func(key string) error {
		exists, err := testutil.KeyValueExists(context.Background(), tc.Pool, key)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("expected key %q to exist, but it does not", key)
		}
		return nil
	})

	sc.Step(`^the error message should contain "([^"]*)"$`, func(expected string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			// Response might not be JSON for some errors
			body := tc.Response.Body.String()
			if strings.Contains(strings.ToLower(body), strings.ToLower(expected)) {
				return nil
			}
			return fmt.Errorf("expected error message to contain %q but got %q", expected, body)
		}
		msg, ok := result["message"]
		if !ok {
			msg = result["error"]
		}
		if msg == nil {
			return fmt.Errorf("no error message in response: %v", result)
		}
		msgStr := fmt.Sprintf("%v", msg)
		if !strings.Contains(strings.ToLower(msgStr), strings.ToLower(expected)) {
			return fmt.Errorf("expected error message to contain %q but got %q", expected, msgStr)
		}
		return nil
	})

	sc.Step(`^I wait for (\d+) milliseconds$`, func(ms int) error {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return nil
	})
}
