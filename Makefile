.DEFAULT_GOAL := help
.PHONY: help setup lint lint-code lint-repository lint-fix test audit build ci

LINTER_IMAGE := next-ip-linter:local

help: ## Show development commands
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "%-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Download and verify Go dependencies
	go mod download
	go mod verify

lint-code: ## Check formatting, module tidiness and Go diagnostics
	@test -z "$$(gofmt -l $$(git ls-files '*.go'))" || { gofmt -l $$(git ls-files '*.go'); exit 1; }
	go mod tidy -diff
	golangci-lint run --timeout=5m

test: ## Run race tests with atomic coverage
	go test -mod=readonly -race -count=1 -timeout=5m -covermode=atomic -coverprofile=coverage.out ./...

audit: ## Audit Go dependencies and the standard library
	go run golang.org/x/vuln/cmd/govulncheck@v1.10.0 ./...

build: ## Build the application without publishing
	go build -mod=readonly -trimpath -o bin/next-ip ./cmd/next-ip

lint-fix: ## Apply Go and repository formatting fixes
	go fmt ./...
	golangci-lint run --fix --timeout=5m
	$(call run_linter,-e FIX_YAML_PRETTIER=true -e FIX_MARKDOWN=true -e FIX_MARKDOWN_PRETTIER=true -e FIX_BIOME_FORMAT=true -e FIX_SHELL_SHFMT=true)

lint: lint-code lint-repository ## Run application and repository checks

lint-repository: ## Run Hoverkraft's shared repository linter
	$(call run_linter,)

ci: setup ## Run local validation without rewriting tracked files
	$(MAKE) lint test audit build

define run_linter
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --tag $(LINTER_IMAGE) .github/linters
	docker run --rm -v "$(CURDIR):$(CURDIR)" -w "$(CURDIR)" -e DEFAULT_WORKSPACE="$(CURDIR)" $(1) $(LINTER_IMAGE)
endef
