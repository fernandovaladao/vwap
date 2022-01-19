.DEFAULT_GOAL := install
IMAGE=zero-hash-vwap
VERSION=v1.0.0
CONTAINER=vwap
DOCKER_BUILDKIT=DOCKER_BUILDKIT=1

install: build test

.PHONY: clean
clean:
	docker rmi $(IMAGE):$(VERSION)
	
.PHONY: build
build:
	$(DOCKER_BUILDKIT) docker build . \
		-t $(IMAGE):$(VERSION) \
		--target bin

.PHONY: test
test:
	$(DOCKER_BUILDKIT) docker build . \
		--target unit-test
	$(DOCKER_BUILDKIT) docker build . \
		--target integration-test

.PHONY: run
run:
	@docker run \
		--name $(CONTAINER) \
		--rm \
		$(IMAGE):$(VERSION)
