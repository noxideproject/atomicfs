SHELL = bash

default: test

.PHONY: test
test: vet
	@echo "==> Running Tests ..."
	@go test -count=1 -v -race ./...

.PHONY: copywrite
copywrite:
	@echo "==> Checking Copywrite ..."
	copywrite --config .copywrite.hcl headers --spdx "BSD-3-Clause"

.PHONY: vet
vet:
	@echo "==> Vet Go sources ..."
	@go vet ./...

.PHONY: lint
lint: vet
	@echo "==> Lint ..."
	@golangci-lint run
