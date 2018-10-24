GO_FILES?=$$(find . -name '*.go' |grep -v vendor)
TAG?=latest

.PHONY: build
build:
	./build.sh

.PHONY: build_redist
build_redist:
	./build_redist.sh

.PHONY: build_samples
build_samples:
	./build_samples.sh

.PHONY: local-fmt
local-fmt:
	gofmt -l -d $(GO_FILES)

.PHONY: local-goimports
local-goimports:
	goimports -w $(GO_FILES)

.PHONY: test-unit
test-unit:
	go test $(shell go list ./... | grep -v /vendor/ | grep -v /template/ | grep -v build) -cover

ci-armhf-push:
	(docker push openfaas/faas-cli:$(TAG)-armhf)
ci-armhf-build:
	(./build.sh $(TAG)-armhf)

.PHONY: test-templating
PORT?=38080
FUNCTION?=templating-test-func
FUNCTION_UP_TIMEOUT?=30
.EXPORT_ALL_VARIABLES:
test-templating:
	./build_integration_test.sh

.PHONY: check-license
check-license:
	@/usr/bin/license-check -path ./ --verbose=false "Alex Ellis" "OpenFaaS Author(s)" "OpenFaaS Project"

.PHONY: check-gofmt
check-gofmt:
	@test -z "$$(gofmt -l $$(find . -type f -name '*.go' -not -path "./vendor/*"))" || \
		{ echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }
