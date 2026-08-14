.PHONY: test coverage cover unit lint format

GOTESTSUM := $(shell command -v gotestsum 2> /dev/null)
GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)
GOFUMPT := $(shell command -v gofumpt 2> /dev/null)

install-gotestsum:
ifndef GOTESTSUM
	go install gotest.tools/gotestsum@latest
endif

install-golangci-lint:
ifndef GOLANGCI_LINT
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
endif

test: install-gotestsum
	gotestsum --format pkgname -- ./...

unit: test

coverage: install-gotestsum
	gotestsum --format pkgname -- -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

cover: coverage

lint: install-golangci-lint
	golangci-lint run ./...

install-gofumpt:
ifndef GOFUMPT
	go install mvdan.cc/gofumpt@latest
endif

format: install-gofumpt
	gofumpt -w .
