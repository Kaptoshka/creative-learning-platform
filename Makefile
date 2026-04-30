SHELL := /bin/bash

GO_LINT := golangci-lint
HADOLINT := hadolint

DOCKERFILES := $(shell find . -type f \( -name "Dockerfile" -o -name "Dockerfile.dev" \))
GO_FILES := $(shell find . -name "*.go")
YAML_FILES := $(shell find . -name "*.yml" -o -name "*.yaml")
NIX_FILES := $(shell find . -name "*.nix")
PROTO_DIR := libs/protos

.PHONY: lint
lint:
	@make -j lint-go lint-nix lint-yaml lint-proto lint-docker lint-compose

.PHONY: lint-go
lint-go:
	@echo "==> Linting Go..."
	@$(GO_LINT) run --config=.golangci.yml

.PHONY: lint-nix
lint-nix:
	@echo "==> Linting Nix..."
	@nix run nixpkgs#nixpkgs-fmt -- --check .
	@nix run nixpkgs#statix -- check .
	@nix run nixpkgs#deadnix -- --fail .

.PHONY: lint-yaml
lint-yaml:
		@echo "==> Linting YAML..."
		@yamllint -c .yamllint.yml .

.PHONY: lint-proto
lint-proto:
	@echo "==> Linting Protobuf..."
	@buf lint $(PROTO_DIR)

.PHONY: lint-docker
lint-docker:
	@echo "==> Linting Dockerfiles..."
	@set -e; \
	for file in $(DOCKERFILES); do \
		echo "Linting $$file"; \
		$(HADOLINT) --config .hadolint.yaml $$file; \
	done

.PHONY: lint-compose
lint-compose:
	@echo "==> Validating docker-compose..."
	@docker compose -f docker-compose.yml config
	@docker compose -f docker-compose.dev.yml config

.PHONY: gen-proto

gen-proto:
	buf generate libs/protos --template libs/protos/buf.gen.yaml
