SHELL := /bin/bash

ULE_DIRS = .
GOLANGCI_VERSION=1.61.0

.PHONY: test
test:
	go test ./... -cover

.PHONY: cover
cover:
	go test $(go list ./... | grep -v '^./mock') -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: mock_gen
mock_gen:	
	mockery

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[lint] golangci-lint: $(mod)" && \
		go run github.com/golangci/golangci-lint/cmd/golangci-lint@v${GOLANGCI_VERSION} run $(ARGS) --path-prefix $(mod) ./...) &&) true
