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
	@$(GO_LINT) run --config=.golangci.yml ./apps/sso/... ./apps/gateway/...

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
		$(HADOLINT) --config .hadolint.yml $$file; \
	done

.PHONY: lint-compose
lint-compose:
	@echo "==> Validating docker-compose..."
	@docker compose -f docker-compose.yml config
	@docker compose -f docker-compose.dev.yml config

.PHONY: gen-proto

gen-proto:
	buf generate libs/protos --template libs/protos/buf.gen.yaml

.PHONY: compose-dev-up
compose-dev-up:
	docker compose -f docker-compose.dev.yml up --build

.PHONY: compose-dev-restart
compose-dev-restart:
	docker compose -f docker-compose.dev.yml restart

.PHONY: compose-dev-logs
compose-dev-logs:
	docker compose -f docker-compose.dev.yml logs -f

.PHONY: compose-dev-stop
compose-dev-stop:
	docker compose -f docker-compose.dev.yml down

.PHONY: format
format:
	@make -j format-go format-nix format-proto

.PHONY: format-go
format-go:
	@echo "==> Formatting Go..."
	@gofmt -w $(GO_FILES)

.PHONY: format-nix
format-nix:
	@echo "==> Formatting Nix..."
	@nix run nixpkgs#nixpkgs-fmt -- .

.PHONY: format-proto
format-proto:
	@echo "==> Formatting Protobuf..."
	@buf format $(PROTO_DIR) --write
