# run rpm test from specs folder
test:
	make build
	cd specs && npm run test:no-build

# run rpm test from specs folder without building the docker image
test-only:
	cd specs && npm run test:no-build

# build the docker image and npm install in the specs folder. Make this the default target
build:
	cd source && docker-compose build
	cd specs && npm install
