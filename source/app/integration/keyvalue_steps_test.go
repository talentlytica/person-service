package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"person-service/integration/testutil"

	"github.com/cucumber/godog"
)

func registerKeyValueSteps(sc *godog.ScenarioContext, tc *TestContext) {
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
		if tc.Response.Code != 200 {
			return fmt.Errorf("expected status 200 but got %d: %s", tc.Response.Code, tc.Response.Body.String())
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
}
