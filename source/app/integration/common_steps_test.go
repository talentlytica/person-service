package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"person-service/integration/testutil"

	"github.com/cucumber/godog"
)

func registerCommonSteps(sc *godog.ScenarioContext, tc *TestContext) {
	// Background steps
	sc.Step(`^the service is running$`, func() error {
		// Service is always running since we use httptest
		return nil
	})

	sc.Step(`^the persons and attributes table is empty$`, func() error {
		return testutil.TruncateTables(context.Background(), tc.Pool)
	})

	sc.Step(`^I have a valid API key$`, func() error {
		// API keys are set in testutil.NewTestServer
		return nil
	})

	// Response assertions
	sc.Step(`^the response status should be (\d+)$`, func(status int) error {
		if tc.Response.Code != status {
			return fmt.Errorf("expected status %d but got %d, body: %s", status, tc.Response.Code, tc.Response.Body.String())
		}
		return nil
	})

	sc.Step(`^the response content type should be "([^"]*)"$`, func(contentType string) error {
		ct := tc.Response.Header().Get("Content-Type")
		if !strings.Contains(ct, contentType) {
			return fmt.Errorf("expected content type %s but got %s", contentType, ct)
		}
		return nil
	})

	sc.Step(`^the response should be valid JSON$`, func() error {
		var result interface{}
		err := json.Unmarshal(tc.Response.Body.Bytes(), &result)
		if err != nil {
			return fmt.Errorf("response is not valid JSON: %v", err)
		}
		return nil
	})

	sc.Step(`^the response should be an object$`, func() error {
		var result map[string]interface{}
		err := json.Unmarshal(tc.Response.Body.Bytes(), &result)
		if err != nil {
			return fmt.Errorf("response is not a JSON object: %v", err)
		}
		tc.JSONResponse = result
		return nil
	})

	sc.Step(`^the response should contain "([^"]*)" with value "([^"]*)"$`, func(field, value string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		tc.JSONResponse = result

		actual, ok := result[field]
		if !ok {
			return fmt.Errorf("field %s not found in response", field)
		}
		if fmt.Sprintf("%v", actual) != value {
			return fmt.Errorf("expected %s to be %s but got %v", field, value, actual)
		}
		return nil
	})

	sc.Step(`^the response should have field "([^"]*)"$`, func(field string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		if _, ok := result[field]; !ok {
			return fmt.Errorf("field %s not found in response", field)
		}
		return nil
	})

	sc.Step(`^the field "([^"]*)" should be one of "([^"]*)"$`, func(field, values string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		actual, ok := result[field]
		if !ok {
			return fmt.Errorf("field %s not found", field)
		}

		allowedValues := strings.Split(values, ",")
		actualStr := fmt.Sprintf("%v", actual)
		for _, v := range allowedValues {
			if actualStr == v {
				return nil
			}
		}
		return fmt.Errorf("field %s value %v is not one of %v", field, actual, allowedValues)
	})

	sc.Step(`^the error message should indicate "([^"]*)"$`, func(expectedMsg string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		msg, ok := result["message"]
		if !ok {
			return fmt.Errorf("no message field in error response")
		}
		if !strings.Contains(fmt.Sprintf("%v", msg), expectedMsg) {
			return fmt.Errorf("expected message to contain %q but got %q", expectedMsg, msg)
		}
		return nil
	})

	sc.Step(`^the error should contain validation details$`, func() error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		// Check that there's some error-related field
		if _, ok := result["message"]; !ok {
			if _, ok := result["error"]; !ok {
				return fmt.Errorf("no error field in response")
			}
		}
		return nil
	})

	sc.Step(`^the error should indicate missing required field "([^"]*)"$`, func(field string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		msg := fmt.Sprintf("%v", result["message"])
		if !strings.Contains(strings.ToLower(msg), strings.ToLower(field)) {
			return fmt.Errorf("expected error to mention %q but got %q", field, msg)
		}
		return nil
	})

	sc.Step(`^the response should indicate success$`, func() error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		// Accept any 2xx status as success
		if tc.Response.Code < 200 || tc.Response.Code >= 300 {
			return fmt.Errorf("expected success status but got %d", tc.Response.Code)
		}
		return nil
	})
}
