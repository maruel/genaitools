# Copyright 2026 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Format, verify, build and test the Go tools and agent file index.
.DEFAULT_GOAL := help
.PHONY: build fix git-hooks help test verify

build:
	@go build ./...

fix:
	@go tool golangci-lint fmt
	@uv run --locked ruff format scripts
	@uv run --locked python scripts/update_agents_file_index.py

verify:
	@go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
	@go tool golangci-lint run --show-stats=false ./...
	@go vet ./...
	@go test -run '^$$' ./...
	@go run github.com/google/addlicense@v1.2.0 -ignore default_config.yml -ignore '**/testdata/**' -ignore pyproject.toml -ignore '.venv/**' -ignore '.ruff_cache/**' -check .
	@uv run --locked python scripts/lint_binaries.py
	@uv run --locked python scripts/update_agents_file_index.py --check
	@uv run --locked ruff check scripts
	@uv run --locked ruff format --check scripts

test:
	@go test -timeout=600s -v -covermode=count -coverprofile coverage.txt -bench=. -benchtime=1x ./...

git-hooks:
	@./scripts/install-git-hooks.sh

help:
	@printf '  %-14s - %s\n' 'make help' 'Show available targets'
	@printf '  %-14s - %s\n' 'make fix' 'Format sources and refresh the agent index'
	@printf '  %-14s - %s\n' 'make verify' 'Run static checks without tests'
	@printf '  %-14s - %s\n' 'make test' 'Run tests with coverage and benchmark smoke tests'
	@printf '  %-14s - %s\n' 'make build' 'Build all Go packages'
	@printf '  %-14s - %s\n' 'make git-hooks' 'Install Git hooks'
