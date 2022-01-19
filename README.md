# vwap

export DOCKER_BUILDKIT=1

docker run -v /etc/ssl/certs:/etc/ssl/certs zero-hash-vwap:v1.0.0

.PHONY: run
run:
	@docker run  \
	-v `pwd`/certs:/etc/ssl/certs \
	--name zero-hash-vwap \
	zero-hash-vwap:v1.0.0

		DOCKER_BUILDKIT=1 docker build . \
		-f Dockerfile.functional \
		-t zero-hash-vwap-functional-tests:latest \
		--target functional-test
	docker run -it \
	-v /var/run/docker.sock:/var/run/docker.sock \
	zero-hash-vwap-functional-tests:latest \
	go test -v --tags functional