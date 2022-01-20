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
	DOCKER_BUILDKIT=1 docker build . \
		-f Dockerfile.functional \
		-t zero-hash-vwap-functional-tests:latest \
		--target functional-test
	docker run -it \
		-v /var/run/docker.sock:/var/run/docker.sock \
		--rm \
		zero-hash-vwap-functional-tests:latest \
		go test -v --tags functional
	docker rmi zero-hash-vwap-functional-tests:latest

.PHONY: run
run:
	@docker run \
		--name $(CONTAINER) \
		--rm \
		$(IMAGE):$(VERSION) BTC-USD ETH-BTC ETH-USD
