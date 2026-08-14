.PHONY: test coverage show-coverage cover unit

GOTESTSUM := $(shell command -v gotestsum 2> /dev/null)

install-gotestsum:
ifndef GOTESTSUM
	go install gotest.tools/gotestsum@latest
endif

test: install-gotestsum
	gotestsum --format pkgname -- ./...

unit: test

coverage: install-gotestsum
	gotestsum --format pkgname -- -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

cover: coverage
