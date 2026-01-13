# install all dependencies required to run docker compose up
install:
	@echo "Checking and installing dependencies..."
	@command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Visit https://docs.docker.com/get-docker/"; exit 1; }
	@command -v docker compose >/dev/null 2>&1 || { echo "Docker Compose is required but not installed. Visit https://docs.docker.com/compose/install/"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "Go 1.24+ is required but not installed. Visit https://golang.org/dl/"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "Node.js is required but not installed. Visit https://nodejs.org/"; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "npm is required but not installed. Visit https://docs.npmjs.com/getting-started"; exit 1; }
	@command -v python3 >/dev/null 2>&1 || { echo "Python3 is required but not installed. Visit https://www.python.org/downloads/"; exit 1; }
	@echo "✓ All system dependencies are installed"
	cd source/app && go mod download
	cd specs && npm install
	@echo "✓ Project dependencies installed successfully"
	@echo ""
	@echo "Dependencies installed:"
	@echo "  - Docker & Docker Compose"
	@echo "  - Go modules"
	@echo "  - npm packages"
	@echo ""
	@echo "Ready to run: make build"

# run rpm test from specs folder
test:
	make build
	cd specs && npm run test:no-build

# run rpm test from specs folder without building the docker image
test-only:
	cd specs && npm run test:no-build

# run unit tests
test-unit:
	cd source/app && go test -v ./...

# run unit tests with coverage
test-unit-coverage:
	cd source/app && go test -coverprofile=../../specs/coverage/unit.out ./...
	cd source/app && go tool cover -func=../../specs/coverage/unit.out

# run mutation tests (requires gremlins to be installed: go install github.com/go-gremlins/gremlins/cmd/gremlins@latest)
test-mutation:
	cd source/app && $(shell go env GOPATH)/bin/gremlins unleash --exclude-files="vendor/.*" --exclude-files="internal/db/generated/.*" --exclude-files=".*_test\\.go$$" --exclude-files="main\\.go" --integration --timeout-coefficient=10 .

# run tests with coverage (integration + unit tests merged)
test-coverage:
	cd specs && npm run test:coverage

# merge coverage from unit tests and integration tests (run after test-coverage has collected go-data)
merge-coverage:
	cd specs && npm run coverage:merge

# serve the coverage report
serve-coverage:
	@echo "Serving coverage report at http://localhost:8181/report.html"
	cd specs/coverage && python3 -m http.server 8181

# build the docker image and npm install in the specs folder. Make this the default target
build:
	cd source && docker compose build
	cd specs && npm install

# Start the backend and the web UI (nginx proxy) for local development
.PHONY: webui webui-down webui-logs
webui:
	docker compose -f source/webui/docker-compose.yml up --build -d

webui-down:
	docker compose -f source/webui/docker-compose.yml down

webui-logs:
	docker compose -f source/webui/docker-compose.yml logs -f
