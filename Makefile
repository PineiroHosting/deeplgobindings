SHELL := /bin/bash

all: test

def:
	@echo "Available commands: imports, test"

test:
	@echo "Testing..."
	@DEEPL_TEST_AUTH_KEY=$$AUTHKEY  go test -race $(shell go list ./... | grep -v /vendor/ | grep -v /cmd/)
