all: build test

.PHONY: build
build:
	@docker build . \
	-f Dockerfile.build \
	--target bin \
	--output bin/

.PHONY: test
test:
	@docker build . \
	-f Dockerfile.build \
	--target unit-test
