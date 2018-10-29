SHELL := /bin/bash

all: imports test

def:
	@echo "Available commands: imports, test"

imports:
	@echo "Restoring imports..."
	@dep ensure

test:
	@echo "Testing..."
	@DEEPL_TEST_AUTH_KEY=$$AUTHKEY  go test -race $(shell go list ./... | grep -v /vendor/ | grep -v /cmd/)