# Mutation Testing

This project uses [Gremlins](https://github.com/go-gremlins/gremlins) for mutation testing of Go code.

## What is Mutation Testing?

Mutation testing is a technique to evaluate the quality of test suites by introducing small changes (mutations) to the code and checking if the tests catch these changes. A mutation that is not detected by tests is called a "surviving mutant" and indicates a gap in test coverage.

## Installation

Install Gremlins:

```bash
go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
```

## Running Mutation Tests

Run mutation tests using the Makefile:

```bash
make test-mutation
```

Or directly with gremlins:

```bash
cd source/app
gremlins unleash --exclude-files="vendor/.*" --exclude-files="internal/db/generated/.*" --exclude-files=".*_test\\.go$" --exclude-files="main\\.go" --integration --timeout-coefficient=10 .
```

## Current Status

- **Test Efficacy**: 96.15%
- **Mutator Coverage**: 100%
- **Killed Mutants**: 50
- **Survived Mutants**: 2
- **Timed Out**: 2

### Surviving Mutants

Two mutants currently survive:

1. **Line 140 in person_attributes/person_attributes.go**: Logging error handling (non-functional, best-effort logging)
2. **Line 392 in person_attributes/person_attributes.go**: Key change check in UpdateAttribute (complex integration scenario)

These surviving mutants are acceptable as they represent:
- Non-critical logging functionality
- Complex integration scenarios that would require extensive mocking to test properly

## Configuration

Mutation testing is configured to:
- Exclude vendor directories
- Exclude generated database code
- Exclude test files themselves
- Exclude the main.go entry point (tested via integration tests)
- Run in integration mode (full test suite for each mutation)
- Use a timeout coefficient of 10 to prevent false timeouts

## Understanding Results

- **KILLED**: The mutation was detected by tests (good!)
- **LIVED**: The mutation was not detected by tests (potential test gap)
- **TIMED OUT**: The test took too long (may indicate an issue)
- **NOT COVERED**: The code is not covered by tests
- **NOT VIABLE**: The mutation couldn't be applied

Test efficacy is calculated as: `KILLED / (KILLED + LIVED) * 100%`

## Improving Test Efficacy

To improve test efficacy:

1. Identify surviving mutants from the gremlins output
2. Examine the code at those locations
3. Add or improve tests to catch the specific mutations
4. Re-run mutation tests to verify improvements

Example mutations that should be caught:
- Changing `==` to `!=`
- Changing `<` to `<=`
- Changing `&&` to `||`
- Inverting negations
- Modifying arithmetic operations

## CI/CD Integration

Mutation testing is computationally expensive and time-consuming. Consider:

- Running mutation tests as a separate CI job
- Running them on a schedule rather than on every commit
- Using mutation testing as a quality gate for major releases
- Tracking test efficacy trends over time
