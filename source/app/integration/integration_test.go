package integration

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"person-service/integration/testutil"

	"github.com/cucumber/godog"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestContext holds state shared across step definitions
type TestContext struct {
	Pool       *pgxpool.Pool
	Server     *testutil.TestServer
	Response   *httptest.ResponseRecorder
	PersonID   string
	PersonName string
	ClientID   string
	// Store attribute IDs by key for reference
	AttributeIDs map[string]int
	// Store the last created attribute ID
	LastAttributeID int
	// Store request data for idempotency tests
	LastTraceID string
	// Store JSON response for assertions
	JSONResponse map[string]interface{}
	// Store array response
	ArrayResponse []map[string]interface{}
	// Store method and path for deferred requests
	LastMethod string
	LastPath   string
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	tc := &TestContext{
		AttributeIDs: make(map[string]int),
	}

	// Setup before each scenario
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		pool, err := testutil.GetPool(ctx)
		if err != nil {
			return ctx, err
		}
		tc.Pool = pool

		// Run migrations (idempotent)
		if err := testutil.RunMigrations(ctx, pool); err != nil {
			return ctx, err
		}

		// Truncate tables for isolation
		if err := testutil.TruncateTables(ctx, pool); err != nil {
			return ctx, err
		}

		// Create test server
		tc.Server = testutil.NewTestServer(pool)

		// Reset state
		tc.Response = nil
		tc.PersonID = ""
		tc.PersonName = ""
		tc.ClientID = ""
		tc.AttributeIDs = make(map[string]int)
		tc.LastAttributeID = 0
		tc.LastTraceID = ""
		tc.JSONResponse = nil
		tc.ArrayResponse = nil

		return ctx, nil
	})

	// Register step definitions
	registerHealthSteps(sc, tc)
	registerKeyValueSteps(sc, tc)
	registerPersonAttributesSteps(sc, tc)
	registerCommonSteps(sc, tc)
}

// Cleanup after all tests
func TestMain(m *testing.M) {
	code := m.Run()

	// Cleanup container
	ctx := context.Background()
	_ = testutil.CleanupContainer(ctx)

	os.Exit(code)
}
