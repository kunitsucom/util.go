SHELL    := /usr/bin/env bash -Eeu -o pipefail
GITROOT  := $(shell git rev-parse --show-toplevel || pwd || echo '.')
PRE_PUSH := ${GITROOT}/.git/hooks/pre-push

.DEFAULT_GOAL := help
.PHONY: help
help: githooks ## display this help documents
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: githooks
githooks:
	@[[ -f "${PRE_PUSH}" ]] || cp -ai "${GITROOT}/.githooks/pre-push" "${PRE_PUSH}"

.PHONY: setup
setup:  ## Setup tools for development
	# direnv
	./.bin/direnv allow
	# golangci-lint
	./.bin/golangci-lint --version

.PHONY: lint
lint:  ## Run secretlint, go mod tidy, golangci-lint
	# ref. https://github.com/secretlint/secretlint
	docker run -v `pwd`:`pwd` -w `pwd` --rm secretlint/secretlint secretlint "**/*"
	# tidy
	go mod tidy
	git diff --exit-code go.mod go.sum
	# lint
	# ref. https://golangci-lint.run/usage/linters/
	./.bin/golangci-lint run --fix --sort-results
	git diff --exit-code

.PHONY: test
test: githooks ## Run go test and display coverage
	# test
	go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt ./...
	go tool cover -func=./coverage.txt

.PHONY: ci
ci: lint credits test ## CI command set
