SHELL             := /usr/bin/env bash -Eeu -o pipefail
REPO_ROOT         := $(shell git rev-parse --show-toplevel)
MAKEFILE_DIR      := $(shell { cd "$(subst /,,$(dir $(lastword ${MAKEFILE_LIST})))" && pwd; } || pwd)
DOTLOCAL_DIR      := ${MAKEFILE_DIR}/.local
PRE_PUSH          := ${REPO_ROOT}/.git/hooks/pre-push

export PATH := ${DOTLOCAL_DIR}/bin:${REPO_ROOT}/.bin:${PATH}

.DEFAULT_GOAL := help
.PHONY: help
help: githooks ## Display this help documents
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: githooks
githooks:
	@[ -f "${PRE_PUSH}" ] || cp -aiv "${REPO_ROOT}/.githooks/pre-push" "${PRE_PUSH}"

.PHONY: setup
setup: githooks ## Setup tools for development
	# == SETUP =====================================================
	# versenv
	make versenv
	# --------------------------------------------------------------

.PHONY: versenv
versenv:
	# direnv
	direnv allow .
	# golangci-lint
	golangci-lint --version

.PHONY: clean
clean:  ## Clean up cache, etc
	# reset tmp
	rm -rf ${MAKEFILE_DIR}/.tmp
	mkdir -p ${MAKEFILE_DIR}/.tmp
	# go build cache
	go env GOCACHE
	go clean -x -cache -testcache -modcache -fuzzcache
	# golangci-lint cache
	golangci-lint cache status
	golangci-lint cache clean

.PHONY: lint
lint: githooks ## Run secretlint, go mod tidy, golangci-lint
	# tidy
	go-mod-tidy.sh
	git diff --exit-code go.mod go.sum
	# golangci-lint
	# ref. https://golangci-lint.run/usage/linters/
	golangci-lint run -c "${REPO_ROOT}/.golangci.yml" --fix --sort-results --verbose --timeout=5m
	# diff
	git diff --exit-code
	# ref. https://github.com/secretlint/secretlint
	docker run -v "${REPO_ROOT}:${REPO_ROOT}:ro" -w "`pwd`" --rm secretlint/secretlint secretlint --secretlintignore="${REPO_ROOT}/.gitignore" "**/*"

.PHONY: test
test: githooks ## Run go test and display coverage
	@[ -x "${DOTLOCAL_DIR}/bin/godotnev" ] || GOBIN="${DOTLOCAL_DIR}/bin" go install github.com/joho/godotenv/cmd/godotenv@latest
	# Unit testing
	godotenv -f .test.env go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt ./... && go tool cover -func=./coverage.txt
	# Unit testing (with external modules)
	cd grpc && godotenv -f .test.env go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt ./... && go tool cover -func=./coverage.txt
	# Integration testing
	cd integrationtest && godotenv -f .test.env go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt ./... && go tool cover -func=./coverage.txt

.PHONY: bench
bench: ## Run benchmarks
	cd integrationtest && go test -run "^NoSuchTestForBenchmark" -benchmem -bench . github.com/kunitsucom/util.go/integrationtest/database/sql -v -trimpath -race -p=4 -parallel=8 -timeout=30s

.PHONY: ci
ci: lint test ## CI command set

.PHONY: git-tag-go-mod
git-tag-go-mod:  ## git tag per go modules
	${REPO_ROOT}/.bin/git-tag-go-mod.sh
