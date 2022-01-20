.DEFAULT_GOAL := install
IMAGE=zero-hash-vwap
VERSION=v1.0.0
E2E-IMAGE=zero-hash-vwap-e2e-tests
E2E-VERSION=latest
CONTAINER=vwap
DOCKER_BUILDKIT=DOCKER_BUILDKIT=1

.PHONY: install
install: build test

.PHONY: clean
clean:
	@docker rmi $(IMAGE):$(VERSION)
	
.PHONY: build
build:
	@$(DOCKER_BUILDKIT) docker build . \
		-t $(IMAGE):$(VERSION) \
		--target bin

.PHONY: test
test: unit-test integration-test e2e-test

.PHONY: unit-test
unit-test:
	@$(DOCKER_BUILDKIT) docker build . \
		--target unit-test
	@docker image prune \
		-f \
		--filter label=stage=unit-test

.PHONY: integration-test
integration-test:
	@$(DOCKER_BUILDKIT) docker build . \
		--target integration-test
	@docker image prune \
		-f \
		--filter label=stage=integration-test

.PHONY: e2e-test
e2e-test:
	@DOCKER_BUILDKIT=1 docker build . \
		-t $(E2E-IMAGE):$(E2E-VERSION) \
		--target e2e-test
	@docker run -it \
		-v /var/run/docker.sock:/var/run/docker.sock \
		--rm \
		$(E2E-IMAGE):$(E2E-VERSION) \
		go test -v -tags e2e

.PHONY: run
run:
	@docker run \
		--name $(CONTAINER) \
		--rm \
		$(IMAGE):$(VERSION) $(TRADE_PAIRS)
