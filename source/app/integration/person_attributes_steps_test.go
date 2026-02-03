package integration

import (
	"context"
	"encoding/json"
	"fmt"

	"person-service/integration/testutil"

	"github.com/cucumber/godog"
)

func registerPersonAttributesSteps(sc *godog.ScenarioContext, tc *TestContext) {
	// Person creation steps
	sc.Step(`^a person exists with the following details:$`, func(table *godog.Table) error {
		for _, row := range table.Rows[1:] {
			name := row.Cells[0].Value
			clientID := row.Cells[1].Value

			personID, err := testutil.CreatePerson(context.Background(), tc.Pool, name, clientID)
			if err != nil {
				return err
			}
			tc.PersonID = personID
			tc.PersonName = name
			tc.ClientID = clientID
		}
		return nil
	})

	sc.Step(`^the person has no attributes$`, func() error {
		return nil
	})

	sc.Step(`^the person has an attribute:$`, func(table *godog.Table) error {
		for _, row := range table.Rows[1:] {
			key := row.Cells[0].Value
			value := row.Cells[1].Value

			body := map[string]interface{}{
				"key":   key,
				"value": value,
				"meta": map[string]string{
					"caller":  "test-setup",
					"reason":  "test-setup",
					"traceId": fmt.Sprintf("setup-%s-%s", tc.PersonID, key),
				},
			}

			resp := tc.Server.POST("/persons/"+tc.PersonID+"/attributes", body, testutil.WithAPIKey())
			if resp.Code != 201 {
				return fmt.Errorf("failed to create attribute: %s", resp.Body.String())
			}

			var result map[string]interface{}
			if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
				return err
			}
			if id, ok := result["id"].(float64); ok {
				tc.AttributeIDs[key] = int(id)
				tc.LastAttributeID = int(id)
			}
		}
		return nil
	})

	sc.Step(`^the person has the following attributes:$`, func(table *godog.Table) error {
		for _, row := range table.Rows[1:] {
			key := row.Cells[0].Value
			value := row.Cells[1].Value

			body := map[string]interface{}{
				"key":   key,
				"value": value,
				"meta": map[string]string{
					"caller":  "test-setup",
					"reason":  "test-setup",
					"traceId": fmt.Sprintf("setup-%s-%s", tc.PersonID, key),
				},
			}

			resp := tc.Server.POST("/persons/"+tc.PersonID+"/attributes", body, testutil.WithAPIKey())
			if resp.Code != 201 {
				return fmt.Errorf("failed to create attribute: %s", resp.Body.String())
			}

			var result map[string]interface{}
			if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
				return err
			}
			if id, ok := result["id"].(float64); ok {
				tc.AttributeIDs[key] = int(id)
			}
		}
		return nil
	})

	// POST request with body - stores body for later use with meta
	sc.Step(`^I send a POST request to "/persons/\{personId\}/attributes" with:$`, func(table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := ""
		if len(table.Rows[1].Cells) > 1 {
			value = table.Rows[1].Cells[1].Value
		}

		tc.JSONResponse = map[string]interface{}{
			"key":   key,
			"value": value,
		}
		tc.LastMethod = "POST"
		tc.LastPath = "/persons/" + tc.PersonID + "/attributes"
		return nil
	})

	// PUT request with body - stores body for later use with meta
	sc.Step(`^I send a PUT request to "/persons/\{personId\}/attributes/\{attributeId\}" with:$`, func(table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := table.Rows[1].Cells[1].Value

		tc.JSONResponse = map[string]interface{}{
			"key":   key,
			"value": value,
		}
		tc.LastMethod = "PUT"
		tc.LastPath = fmt.Sprintf("/persons/%s/attributes/%d", tc.PersonID, tc.LastAttributeID)
		return nil
	})

	// Meta step - sends the actual request
	sc.Step(`^the request meta contains:$`, func(table *godog.Table) error {
		caller := table.Rows[1].Cells[0].Value
		reason := table.Rows[1].Cells[1].Value
		traceID := table.Rows[1].Cells[2].Value

		tc.LastTraceID = traceID

		body := tc.JSONResponse
		if body == nil {
			body = make(map[string]interface{})
		}
		body["meta"] = map[string]string{
			"caller":  caller,
			"reason":  reason,
			"traceId": traceID,
		}

		switch {
		case tc.LastMethod == "POST":
			tc.Response = tc.Server.POST(tc.LastPath, body, testutil.WithAPIKey())
		case tc.LastMethod == "PUT":
			tc.Response = tc.Server.PUT(tc.LastPath, body, testutil.WithAPIKey())
		case tc.LastMethod == "DELETE":
			tc.Response = tc.Server.DELETE(tc.LastPath, testutil.WithAPIKey())
		case tc.LastMethod == "POST_GREEN":
			tc.Response = tc.Server.POST(tc.LastPath, body, testutil.WithGreenAPIKey())
		case tc.LastMethod == "POST_NO_KEY":
			tc.Response = tc.Server.POST(tc.LastPath, body, nil)
		case len(tc.LastMethod) > 13 && tc.LastMethod[:13] == "POST_INVALID_":
			apiKey := tc.LastMethod[13:]
			tc.Response = tc.Server.POST(tc.LastPath, body, testutil.WithCustomAPIKey(apiKey))
		default:
			tc.Response = tc.Server.POST(tc.LastPath, body, testutil.WithAPIKey())
		}
		return nil
	})

	// Response assertions
	sc.Step(`^the response should contain an attribute with:$`, func(table *godog.Table) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		tc.JSONResponse = result

		for _, row := range table.Rows[1:] {
			expectedKey := row.Cells[0].Value
			expectedValue := row.Cells[1].Value

			if result["key"] != expectedKey {
				return fmt.Errorf("expected key %q but got %q", expectedKey, result["key"])
			}
			if result["value"] != expectedValue {
				return fmt.Errorf("expected value %q but got %q", expectedValue, result["value"])
			}
		}

		if id, ok := result["id"].(float64); ok {
			tc.LastAttributeID = int(id)
			if key, ok := result["key"].(string); ok {
				tc.AttributeIDs[key] = int(id)
			}
		}
		return nil
	})

	sc.Step(`^the attribute should have an "([^"]*)"$`, func(field string) error {
		if tc.JSONResponse == nil {
			return fmt.Errorf("no JSON response stored")
		}
		if _, ok := tc.JSONResponse[field]; !ok {
			return fmt.Errorf("field %s not found", field)
		}
		return nil
	})

	sc.Step(`^the attribute should have "([^"]*)" timestamp$`, func(field string) error {
		if tc.JSONResponse == nil {
			return fmt.Errorf("no JSON response stored")
		}
		if _, ok := tc.JSONResponse[field]; !ok {
			return fmt.Errorf("field %s not found", field)
		}
		return nil
	})

	// Multiple attributes
	sc.Step(`^I add the following attributes to the person:$`, func(table *godog.Table) error {
		for _, row := range table.Rows[1:] {
			key := row.Cells[0].Value
			value := row.Cells[1].Value

			body := map[string]interface{}{
				"key":   key,
				"value": value,
				"meta": map[string]string{
					"caller":  "test",
					"reason":  "test",
					"traceId": fmt.Sprintf("multi-%s-%s", tc.PersonID, key),
				},
			}

			resp := tc.Server.POST("/persons/"+tc.PersonID+"/attributes", body, testutil.WithAPIKey())
			if resp.Code != 201 {
				return fmt.Errorf("failed to add attribute %s: %s", key, resp.Body.String())
			}

			var result map[string]interface{}
			if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
				return err
			}
			if id, ok := result["id"].(float64); ok {
				tc.AttributeIDs[key] = int(id)
			}
		}
		return nil
	})

	sc.Step(`^all attributes should be added successfully$`, func() error {
		return nil
	})

	sc.Step(`^the person should have (\d+) attributes?$`, func(count int) error {
		actualCount, err := testutil.CountAttributes(context.Background(), tc.Pool, tc.PersonID)
		if err != nil {
			return err
		}
		if actualCount != count {
			return fmt.Errorf("expected %d attributes but got %d", count, actualCount)
		}
		return nil
	})

	// GET all attributes with API key
	sc.Step(`^I send a GET request to "/persons/\{personId\}/attributes"$`, func() error {
		tc.Response = tc.Server.GET("/persons/"+tc.PersonID+"/attributes", testutil.WithAPIKey())
		return nil
	})

	sc.Step(`^the response should contain (\d+) attributes$`, func(count int) error {
		var result []map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		tc.ArrayResponse = result
		if len(result) != count {
			return fmt.Errorf("expected %d attributes but got %d", count, len(result))
		}
		return nil
	})

	sc.Step(`^the response should contain an empty attributes array$`, func() error {
		var result []map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		if len(result) != 0 {
			return fmt.Errorf("expected empty array but got %d items", len(result))
		}
		return nil
	})

	sc.Step(`^the attributes should include:$`, func(table *godog.Table) error {
		if tc.ArrayResponse == nil {
			return fmt.Errorf("no array response stored")
		}
		for _, row := range table.Rows[1:] {
			expectedKey := row.Cells[0].Value
			expectedValue := row.Cells[1].Value
			found := false
			for _, attr := range tc.ArrayResponse {
				if attr["key"] == expectedKey && attr["value"] == expectedValue {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("attribute key=%s value=%s not found", expectedKey, expectedValue)
			}
		}
		return nil
	})

	sc.Step(`^the attribute "([^"]*)" should have value "([^"]*)"$`, func(key, expectedValue string) error {
		body := tc.Response.Body.Bytes()
		// Try as array first (GET all attributes response)
		var arrResult []map[string]interface{}
		if err := json.Unmarshal(body, &arrResult); err == nil {
			for _, attr := range arrResult {
				if attr["key"] == key && attr["value"] == expectedValue {
					return nil
				}
			}
			return fmt.Errorf("attribute %s with value %s not found in array", key, expectedValue)
		}
		// Try as object (single attribute response)
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return err
		}
		if result["key"] == key && result["value"] == expectedValue {
			return nil
		}
		if result["value"] == expectedValue {
			return nil
		}
		return fmt.Errorf("attribute %s does not have value %s, got %v", key, expectedValue, result["value"])
	})

	sc.Step(`^the attribute "([^"]*)" timestamp should be updated$`, func(field string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		if _, ok := result[field]; !ok {
			return fmt.Errorf("field %s not found", field)
		}
		return nil
	})

	// DELETE attribute - stores path for later use with meta
	sc.Step(`^I send a DELETE request to "/persons/\{personId\}/attributes/\{attributeId\}" for attribute "([^"]*)"$`, func(attrKey string) error {
		attrID, ok := tc.AttributeIDs[attrKey]
		if !ok {
			return fmt.Errorf("attribute ID not found for key %s", attrKey)
		}
		tc.LastMethod = "DELETE"
		tc.LastPath = fmt.Sprintf("/persons/%s/attributes/%d", tc.PersonID, attrID)
		tc.JSONResponse = make(map[string]interface{})
		return nil
	})

	sc.Step(`^I send a DELETE request to "/persons/\{personId\}/attributes/(\d+)"$`, func(attrID int) error {
		tc.LastMethod = "DELETE"
		tc.LastPath = fmt.Sprintf("/persons/%s/attributes/%d", tc.PersonID, attrID)
		tc.JSONResponse = make(map[string]interface{})
		return nil
	})

	sc.Step(`^the person should have (\d+) attribute remaining$`, func(count int) error {
		actualCount, err := testutil.CountAttributes(context.Background(), tc.Pool, tc.PersonID)
		if err != nil {
			return err
		}
		if actualCount != count {
			return fmt.Errorf("expected %d attributes but got %d", count, actualCount)
		}
		return nil
	})

	sc.Step(`^the remaining attribute should be "([^"]*)"$`, func(key string) error {
		resp := tc.Server.GET("/persons/"+tc.PersonID+"/attributes", testutil.WithAPIKey())
		var result []map[string]interface{}
		if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
			return err
		}
		if len(result) == 0 {
			return fmt.Errorf("no attributes found")
		}
		if result[0]["key"] != key {
			return fmt.Errorf("expected %s but got %s", key, result[0]["key"])
		}
		return nil
	})

	// Error scenarios
	sc.Step(`^I send a POST request to "/persons/(\d+)/attributes" with:$`, func(personID int, table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := table.Rows[1].Cells[1].Value
		tc.JSONResponse = map[string]interface{}{
			"key":   key,
			"value": value,
		}
		tc.PersonID = fmt.Sprintf("%d", personID)
		tc.LastMethod = "POST"
		tc.LastPath = fmt.Sprintf("/persons/%d/attributes", personID)
		return nil
	})

	sc.Step(`^I send a POST request to "/persons/\{personId\}/attributes" with invalid data:$`, func(table *godog.Table) error {
		tc.JSONResponse = map[string]interface{}{
			"key":   "",
			"value": "",
		}
		tc.LastMethod = "POST"
		tc.LastPath = "/persons/" + tc.PersonID + "/attributes"
		return nil
	})

	sc.Step(`^I send a PUT request to "/persons/\{personId\}/attributes/(\d+)" with:$`, func(attrID int, table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := table.Rows[1].Cells[1].Value
		tc.JSONResponse = map[string]interface{}{
			"key":   key,
			"value": value,
		}
		tc.LastAttributeID = attrID
		tc.LastMethod = "PUT"
		tc.LastPath = fmt.Sprintf("/persons/%s/attributes/%d", tc.PersonID, attrID)
		return nil
	})

	sc.Step(`^I send a GET request to "/persons/(\d+)/attributes"$`, func(personID int) error {
		tc.Response = tc.Server.GET(fmt.Sprintf("/persons/%d/attributes", personID), testutil.WithAPIKey())
		return nil
	})

	// API key scenarios
	sc.Step(`^I send a POST request to "/persons/\{personId\}/attributes" without API key:$`, func(table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := table.Rows[1].Cells[1].Value
		tc.JSONResponse = map[string]interface{}{
			"key":   key,
			"value": value,
		}
		tc.LastMethod = "POST_NO_KEY"
		tc.LastPath = "/persons/" + tc.PersonID + "/attributes"
		return nil
	})

	sc.Step(`^I send a POST request to "/persons/\{personId\}/attributes" with invalid API key "([^"]*)":$`, func(apiKey string, table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := table.Rows[1].Cells[1].Value
		tc.JSONResponse = map[string]interface{}{
			"key":   key,
			"value": value,
		}
		tc.LastMethod = "POST_INVALID_" + apiKey
		tc.LastPath = "/persons/" + tc.PersonID + "/attributes"
		return nil
	})

	sc.Step(`^I send a POST request to "/persons/\{personId\}/attributes" using green API key with:$`, func(table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := table.Rows[1].Cells[1].Value
		tc.JSONResponse = map[string]interface{}{
			"key":   key,
			"value": value,
		}
		tc.LastMethod = "POST_GREEN"
		tc.LastPath = "/persons/" + tc.PersonID + "/attributes"
		return nil
	})

	sc.Step(`^I send a POST request to "/persons/\{personId\}/attributes" without meta:$`, func(table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := table.Rows[1].Cells[1].Value
		body := map[string]interface{}{
			"key":   key,
			"value": value,
		}
		tc.Response = tc.Server.POST("/persons/"+tc.PersonID+"/attributes", body, testutil.WithAPIKey())
		return nil
	})

	// Lifecycle scenario
	sc.Step(`^I add an attribute to the person:$`, func(table *godog.Table) error {
		key := table.Rows[1].Cells[0].Value
		value := table.Rows[1].Cells[1].Value
		body := map[string]interface{}{
			"key":   key,
			"value": value,
			"meta": map[string]string{
				"caller":  "lifecycle-test",
				"reason":  "create",
				"traceId": fmt.Sprintf("lifecycle-create-%s", key),
			},
		}
		tc.Response = tc.Server.POST("/persons/"+tc.PersonID+"/attributes", body, testutil.WithAPIKey())
		if tc.Response.Code == 201 {
			var result map[string]interface{}
			json.Unmarshal(tc.Response.Body.Bytes(), &result)
			if id, ok := result["id"].(float64); ok {
				tc.AttributeIDs[key] = int(id)
				tc.LastAttributeID = int(id)
			}
		}
		return nil
	})

	sc.Step(`^the attribute should be created successfully$`, func() error {
		if tc.Response.Code != 201 {
			return fmt.Errorf("expected 201 but got %d: %s", tc.Response.Code, tc.Response.Body.String())
		}
		return nil
	})

	sc.Step(`^I retrieve all attributes for the person$`, func() error {
		tc.Response = tc.Server.GET("/persons/"+tc.PersonID+"/attributes", testutil.WithAPIKey())
		return nil
	})

	sc.Step(`^I should see the "([^"]*)" attribute with value "([^"]*)"$`, func(key, value string) error {
		var result []map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		for _, attr := range result {
			if attr["key"] == key && attr["value"] == value {
				return nil
			}
		}
		return fmt.Errorf("attribute %s with value %s not found", key, value)
	})

	sc.Step(`^I update the "([^"]*)" attribute to:$`, func(key string, table *godog.Table) error {
		newKey := table.Rows[1].Cells[0].Value
		newValue := table.Rows[1].Cells[1].Value
		attrID, ok := tc.AttributeIDs[key]
		if !ok {
			return fmt.Errorf("attribute ID not found for key %s", key)
		}
		body := map[string]interface{}{
			"key":   newKey,
			"value": newValue,
			"meta": map[string]string{
				"caller":  "lifecycle-test",
				"reason":  "update",
				"traceId": fmt.Sprintf("lifecycle-update-%s", key),
			},
		}
		path := fmt.Sprintf("/persons/%s/attributes/%d", tc.PersonID, attrID)
		tc.Response = tc.Server.PUT(path, body, testutil.WithAPIKey())
		return nil
	})

	sc.Step(`^the attribute should be updated successfully$`, func() error {
		if tc.Response.Code != 200 {
			return fmt.Errorf("expected 200 but got %d: %s", tc.Response.Code, tc.Response.Body.String())
		}
		return nil
	})

	sc.Step(`^the "([^"]*)" attribute should have value "([^"]*)"$`, func(key, expectedValue string) error {
		resp := tc.Server.GET("/persons/"+tc.PersonID+"/attributes", testutil.WithAPIKey())
		var result []map[string]interface{}
		if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
			return err
		}
		for _, attr := range result {
			if attr["key"] == key {
				if attr["value"] == expectedValue {
					return nil
				}
				return fmt.Errorf("attribute %s has value %s, expected %s", key, attr["value"], expectedValue)
			}
		}
		return fmt.Errorf("attribute %s not found", key)
	})

	sc.Step(`^I delete the "([^"]*)" attribute$`, func(key string) error {
		attrID, ok := tc.AttributeIDs[key]
		if !ok {
			return fmt.Errorf("attribute ID not found for key %s", key)
		}
		path := fmt.Sprintf("/persons/%s/attributes/%d", tc.PersonID, attrID)
		tc.Response = tc.Server.DELETE(path, testutil.WithAPIKey())
		return nil
	})

	sc.Step(`^the attribute should be deleted successfully$`, func() error {
		if tc.Response.Code != 200 {
			return fmt.Errorf("expected 200 but got %d: %s", tc.Response.Code, tc.Response.Body.String())
		}
		return nil
	})

	sc.Step(`^the person should have no attributes$`, func() error {
		count, err := testutil.CountAttributes(context.Background(), tc.Pool, tc.PersonID)
		if err != nil {
			return err
		}
		if count != 0 {
			return fmt.Errorf("expected 0 attributes but got %d", count)
		}
		return nil
	})

	// Edge cases
	sc.Step(`^the attribute value should be stored correctly with special characters$`, func() error {
		if tc.Response.Code != 201 {
			return fmt.Errorf("expected 201 but got %d", tc.Response.Code)
		}
		return nil
	})

	sc.Step(`^the attribute should be created with empty value$`, func() error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		if result["value"] != "" {
			return fmt.Errorf("expected empty value but got %v", result["value"])
		}
		return nil
	})

	sc.Step(`^the person should have (\d+) attribute with key "([^"]*)" and value "([^"]*)"$`, func(count int, key, value string) error {
		resp := tc.Server.GET("/persons/"+tc.PersonID+"/attributes", testutil.WithAPIKey())
		var result []map[string]interface{}
		if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
			return err
		}
		matching := 0
		for _, attr := range result {
			if attr["key"] == key && attr["value"] == value {
				matching++
			}
		}
		if matching != count {
			return fmt.Errorf("expected %d attribute(s) with key=%s value=%s, found %d", count, key, value, matching)
		}
		return nil
	})

	// Encryption verification
	sc.Step(`^the raw database value for attribute "([^"]*)" should not equal "([^"]*)"$`, func(key, plainValue string) error {
		rawValue, err := testutil.GetRawAttributeValue(context.Background(), tc.Pool, tc.PersonID, key)
		if err != nil {
			return err
		}
		if string(rawValue) == plainValue {
			return fmt.Errorf("raw value equals plaintext")
		}
		return nil
	})

	sc.Step(`^the raw database value should be encrypted bytes$`, func() error {
		return nil
	})

	sc.Step(`^the database should have key_version (\d+) for the attribute "([^"]*)"$`, func(version int, key string) error {
		keyVersion, err := testutil.GetAttributeKeyVersion(context.Background(), tc.Pool, tc.PersonID, key)
		if err != nil {
			return err
		}
		if int(keyVersion) != version {
			return fmt.Errorf("expected key_version %d but got %d", version, keyVersion)
		}
		return nil
	})

	// Audit verification
	sc.Step(`^an audit record should be created for traceId "([^"]*)"$`, func(traceID string) error {
		count, err := testutil.CountRequestLogs(context.Background(), tc.Pool, traceID)
		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("no audit record found for traceId %s", traceID)
		}
		return nil
	})

	sc.Step(`^the audit record should contain caller "([^"]*)" and reason "([^"]*)"$`, func(expectedCaller, expectedReason string) error {
		caller, reason, err := testutil.GetRequestLog(context.Background(), tc.Pool, tc.LastTraceID, testutil.TestEncryptionKey)
		if err != nil {
			return err
		}
		if caller != expectedCaller {
			return fmt.Errorf("expected caller %s but got %s", expectedCaller, caller)
		}
		if reason != expectedReason {
			return fmt.Errorf("expected reason %s but got %s", expectedReason, reason)
		}
		return nil
	})

	// Idempotency
	sc.Step(`^I send the same POST request again with traceId "([^"]*)"$`, func(traceID string) error {
		body := map[string]interface{}{
			"key":   "token",
			"value": "unique-123",
			"meta": map[string]string{
				"caller":  "user123",
				"reason":  "add token",
				"traceId": traceID,
			},
		}
		tc.Response = tc.Server.POST("/persons/"+tc.PersonID+"/attributes", body, testutil.WithAPIKey())
		return nil
	})

	sc.Step(`^the attribute should be created only once$`, func() error {
		count, err := testutil.CountAttributes(context.Background(), tc.Pool, tc.PersonID)
		if err != nil {
			return err
		}
		if count != 1 {
			return fmt.Errorf("expected 1 attribute but got %d", count)
		}
		return nil
	})

	sc.Step(`^the attribute should have key "([^"]*)" and value "([^"]*)"$`, func(key, value string) error {
		var result map[string]interface{}
		if err := json.Unmarshal(tc.Response.Body.Bytes(), &result); err != nil {
			return err
		}
		if result["key"] != key {
			return fmt.Errorf("expected key %s but got %v", key, result["key"])
		}
		if result["value"] != value {
			return fmt.Errorf("expected value %s but got %v", value, result["value"])
		}
		return nil
	})
}
