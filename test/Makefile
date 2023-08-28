SHELL             := /usr/bin/env bash -Eeu -o pipefail
REPO_ROOT         := $(shell git rev-parse --show-toplevel)
MAKEFILE_DIR      := $(shell { cd "$(subst /,,$(dir $(lastword ${MAKEFILE_LIST})))" && pwd; } || pwd)
DOTLOCAL_DIR      := ${MAKEFILE_DIR}/.local

export PATH := ${DOTLOCAL_DIR}/bin:${REPO_ROOT}/.bin:${PATH}

.DEFAULT_GOAL := help
.PHONY: help
help:  ## Display this help documents
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint:  ## Run secretlint, go mod tidy, golangci-lint
	# tidy
	go mod tidy
	git diff --exit-code go.mod go.sum
	# golangci-lint
	# ref. https://golangci-lint.run/usage/linters/
	golangci-lint run -c "${REPO_ROOT}/.golangci.yml" --fix --sort-results --verbose --timeout=5m
	# diff
	git diff --exit-code
	# ref. https://github.com/secretlint/secretlint
	docker run -v "`pwd`:`pwd`" -w "`pwd`" --rm secretlint/secretlint secretlint --secretlintignore=.gitignore "**/*"

.PHONY: test
test:  ## Run go test and display coverage
	@[[ -x "${DOTLOCAL_DIR}/bin/godotnev" ]] || GOBIN="${DOTLOCAL_DIR}/bin" go install github.com/joho/godotenv/cmd/godotenv@latest
	# test
	godotenv -f .test.env go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt ./...
	go tool cover -func=./coverage.txt

.PHONY: ci
ci: lint test ## CI command set
