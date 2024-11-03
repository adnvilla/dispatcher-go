SHELL := /bin/bash

default: build

test:
	go test ./... -cover

mock_gen:	
	mockery