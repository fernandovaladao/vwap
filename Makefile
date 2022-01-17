.DEFAULT_GOAL := install

install: build test

.PHONY: build
build:
	DOCKER_BUILDKIT=1 docker build . \
	-t zero-hash-vwap:v1.0.0 \
	--target bin

.PHONY: test
test:
	DOCKER_BUILDKIT=1 docker build . \
	--target unit-test

.PHONY: run
run:
	docker run \
	--name zero-hash-vwap \
	zero-hash-vwap:v1.0.0

