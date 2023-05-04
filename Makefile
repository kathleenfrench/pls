WHOAMI ?= $(shell whoami)
CWD := $(shell pwd)
NAME := pls
BIN_NAME := pls
INSTALL_LOCATION := /usr/local/bin
COMMIT := $(shell git rev-parse --short HEAD)
TODAY = $(shell date +%Y-%m-%d)
VERSION := $(COMMIT)-$(TODAY)

BUILD_OUTPUT_DIR := $(CWD)/build
BINARY_LOCATION := $(BUILD_OUTPUT_DIR)/$(BIN_NAME)
MODULE := $(shell go list -m)
CMD_MODULE := $(MODULE)/cmd/$(BIN_NAME)

${BUILD_OUTPUT_DIR}:
	@mkdir -p $(BUILD_OUTPUT_DIR)

UNAME_S := $(shell uname -s)
ifeq ($(PLATFORM),)
ifeq ($(UNAME_S),Darwin)
PLATFORM ?= darwin
endif
ifeq ($(UNAME_S),Linux)
PLATFORM ?= linux
endif
endif

GOOS = $(PLATFORM)
GOARCH ?= amd64

GO := $(shell command -v go 2>/dev/null)
GO_LINKER_FLAGS = -X $(CMD_MODULE).Builder=$(WHOAMI) -X $(CMD_MODULE).Version=$(VERSION) -X $(CMD_MODULE).Commit=$(COMMIT) -X $(CMD_MODULE).Date=$(TODAY)
GO_BUILD_FLAGS = -trimpath -tags 'osusergo netgo' -ldflags "$(GO_LINKER_FLAGS)" -o $(BINARY_LOCATION)

.PHONY: build
build: ${BUILD_OUTPUT_DIR} ## build the pls binary
	@echo "compiling ${NAME}..."
	@export GOOS=$(GOOS) GOARCH=$(GOARCH) && \
		export CGO_ENABLED=0 && \
		$(GO) build $(GO_BUILD_FLAGS)
	@echo "${NAME} bin compiled!"

.PHONY: install
install: build ## install the pls binary to /usr/local/bin
	@echo "installing pls to ${INSTALL_LOCATION}"
	@cp ${BINARY_LOCATION} ${INSTALL_LOCATION}
	@chmod 755 ${INSTALL_LOCATION}/${BIN_NAME}
	@echo "installation complete"

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: lint-fix
lint-fix: ## run linter and fix issues if possible
	@echo "running golangci..."
	@golangci-lint run --fix

.PHONY: lint
lint: ## go linter and shadow tool
	@$(GO) get -u golang.org/x/lint/golint
	@$(GO) get -u golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	@$(GO) vet $(shell go list ./...)
	@$(GO) vet -vettool=$(shell which shadow) ./...

.PHONY: test
test: lint ## run linter and unit tests
	@echo "running tests..."
	@$(GO) test ./...

.PHONY: clean
clean: ## delete the build binary
	@rm -rf ${BUILD_OUTPUT_DIR}
	@echo "removed ${BUILD_OUTPUT_DIR}..."

.PHONY: pages
pages: ## run pages site locally
	@cd docs && \
		bundle && \
		bundle install && \
		bundle exec jekyll serve && \
		cd ..

.PHONY: help
help: ## lists some available makefile commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
