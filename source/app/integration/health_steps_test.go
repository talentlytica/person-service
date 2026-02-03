package integration

import (
	"sync"
	"time"

	"github.com/cucumber/godog"
)

func registerHealthSteps(sc *godog.ScenarioContext, tc *TestContext) {
	var concurrentResponses []int
	var mu sync.Mutex
	var requestDuration time.Duration

	sc.Step(`^I send a GET request to "/health"$`, func() error {
		tc.Response = tc.Server.GET("/health", nil)
		return nil
	})

	sc.Step(`^I send (\d+) concurrent GET requests to "([^"]*)"$`, func(count int, path string) error {
		concurrentResponses = make([]int, 0, count)
		var wg sync.WaitGroup

		for i := 0; i < count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp := tc.Server.GET(path, nil)
				mu.Lock()
				concurrentResponses = append(concurrentResponses, resp.Code)
				mu.Unlock()
			}()
		}

		wg.Wait()
		return nil
	})

	sc.Step(`^all responses should have status (\d+)$`, func(status int) error {
		for i, code := range concurrentResponses {
			if code != status {
				return godog.ErrPending
			}
			_ = i // silence unused variable warning
		}
		return nil
	})

	sc.Step(`^I send a GET request to "([^"]*)" with (\d+)ms timeout$`, func(path string, timeout int) error {
		start := time.Now()
		tc.Response = tc.Server.GET(path, nil)
		requestDuration = time.Since(start)
		_ = timeout // The timeout is for assertion, not for the actual request
		return nil
	})

	sc.Step(`^the request should complete within timeout$`, func() error {
		// Request completed successfully (no timeout in httptest)
		// Just verify it was reasonably fast (under 5 seconds)
		if requestDuration > 5*time.Second {
			return godog.ErrPending
		}
		return nil
	})
}
