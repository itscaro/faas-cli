GOLANG_IMAGE ?= golang:1.8.3
TEST_OPTS ?=

.PHONY: build build_redist test-unit test-unit-docker

build:
	./build.sh

build_redist:
	./build_redist.sh

test-unit:
	go test $(shell go list ./... | grep -v -e /vendor/ -e /template/ -e /build/ -e /sample/) $(TEST_OPTS) -cover
test-unit-docker:
	docker run --rm \
	-v "$(shell pwd):/go/src/github.com/openfaas/faas-cli" \
	-w /go/src/github.com/openfaas/faas-cli \
	$(GOLANG_IMAGE) go test $(shell go list ./... | grep -v -e /vendor/ -e /template/ -e /build/ -e /sample/) $(TEST_OPTS) -cover
