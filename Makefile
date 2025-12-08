SHELL := /bin/bash

MODULE_DIRS = .
GOLANGCI_VERSION=2.7.2

setup:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(GOLANGCI_VERSION)

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
lint: 
	golangci-lint run --fix
