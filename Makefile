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

# run tests with coverage
test-coverage:
	cd specs && npm run test:coverage

# serve the coverage report
serve-coverage:
	@echo "Serving coverage report at http://localhost:8181/report.html"
	cd specs/coverage && python3 -m http.server 8181

# build the docker image and npm install in the specs folder. Make this the default target
build:
	cd source && docker compose build
	cd specs && npm install
