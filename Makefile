# ClusterIQ Makefile
################################################################################

# Global Vars
SHORT_COMMIT_HASH := $(shell git rev-parse --short=7 HEAD)

# Binary vars
CONTAINER_ENGINE ?= $(shell which podman >/dev/null 2>&1 && echo podman || echo docker)
K8S_CLI ?= $(shell which oc >/dev/null 2>&1 && echo oc || echo kubectl)
GO ?= go
GO_LINTER ?= golangci-lint
GO_FMT ?= gofmt
SWAGGER ?= swag
PROTOC ?= protoc

# Required binaries
REQUIRED_BINS := $(CONTAINER_ENGINE) $(CONTAINER_ENGINE)-compose $(K8S_CLI) $(GO) $(GO_LINTER) $(GO_FMT) $(SWAGGER) $(PROTOC)

# Container image registry vars
REGISTRY ?= quay.io
PROJECT_NAME ?= cluster-iq
REGISTRY_REPO ?= ecosystem-appeng
COMPOSE_NETWORK ?= compose_cluster_iq

# Building vars
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(SHORT_COMMIT_HASH)"

# Project directories
TEST_DIR ?= ./test
BUILD_DIR ?= ./build
BIN_DIR ?= $(BUILD_DIR)/bin
CMD_DIR ?= ./cmd
PKG_DIR ?= ./internal
DEPLOYMENTS_DIR ?= ./deployments
GENERATED_DIR ?= ./generated

# Images
API_IMG_NAME ?= $(PROJECT_NAME)-api
API_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/$(API_IMG_NAME)
API_CONTAINERFILE ?= ./$(DEPLOYMENTS_DIR)/containerfiles/Containerfile-api
SCANNER_IMG_NAME ?= $(PROJECT_NAME)-scanner
SCANNER_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/$(SCANNER_IMG_NAME)
SCANNER_CONTAINERFILE ?= ./$(DEPLOYMENTS_DIR)/containerfiles/Containerfile-scanner
AGENT_IMG_NAME ?= $(PROJECT_NAME)-agent
AGENT_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/$(AGENT_IMG_NAME)
AGENT_CONTAINERFILE ?= ./$(DEPLOYMENTS_DIR)/containerfiles/Containerfile-agent
AGENT_PROTO_PATH ?= ./cmd/agent/proto/agent.proto
PGSQL_IMG_NAME ?= $(PROJECT_NAME)-pgsql
PGSQL_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/$(PGSQL_IMG_NAME)
PGSQL_CONTAINERFILE ?= ./$(DEPLOYMENTS_DIR)/containerfiles/Containerfile-pgsql

# Standard targets
all: ## Stop, build and start the development environment based on containers
all: check-dependencies stop-dev build start-dev

.PHONY: check-dependencies
check-dependencies:
	@$(foreach bin,$(REQUIRED_BINS),\
		$(if $(shell command -v $(bin) 2> /dev/null),,\
			$(error ✗ $(bin) is required but not installed)))

# Building in local environment
local-clean:
	@echo "### [Cleaning local builds] ###"
	@rm -Rf $(BIN_DIR)
	@rm -Rf $(GENERATED_DIR)

local-build: local-build-scanner local-build-api local-build-agent ## Build all local binaries

local-build-api: swagger-doc ## Build the API binary
	@echo "### [Building API] ###"
	@$(GO) build -o $(BIN_DIR)/api/api $(LDFLAGS) ./cmd/api/

local-build-scanner: ## Build the scanner binary
	@echo "### [Building Scanner] ###"
	@$(GO) build -o $(BIN_DIR)/scanners/scanner $(LDFLAGS) ./cmd/scanner

local-build-agent: ## Build the agent binary
	@echo "### [Building Agent] ###"
	@[ ! -d $(GENERATED_DIR) ] && { mkdir $(GENERATED_DIR); } || { exit 0; }
	@protoc --go_out=$(GENERATED_DIR) --go-grpc_out=$(GENERATED_DIR) $(AGENT_PROTO_PATH)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_DIR)/agent/agent $(LDFLAGS) ./cmd/agent


# Container based working targets
clean: ## Remove the container images
	@echo "### [Cleaning Container images] ###"
	@-$(CONTAINER_ENGINE) images | grep -e $(SCANNER_IMAGE) -e $(API_IMAGE) -e $(AGENT_IMAGE) -e $(PGSQL_IMAGE) | awk '{print $$3}' | xargs $(CONTAINER_ENGINE) rmi -f

build: ## Build all container images
build: build-api build-scanner build-agent build-pgsql
build-api: ## Build the API container image
	@echo "### [Building API container image] ###"
	@$(CONTAINER_ENGINE) build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(SHORT_COMMIT_HASH) \
		-t $(API_IMAGE):latest -f $(API_CONTAINERFILE) .
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(SHORT_COMMIT_HASH)
	@echo "Build Successful"

build-scanner: ## Build the scanner container image
	@echo "### [Building Scanner container image] ###"
	@$(CONTAINER_ENGINE) build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(SHORT_COMMIT_HASH) \
		-t $(SCANNER_IMAGE):latest -f $(SCANNER_CONTAINERFILE) .
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(SHORT_COMMIT_HASH)
	@echo "Build Successful"

build-agent: ## Build the agent container image
	@echo "### [Building Agent container image] ###"
	@$(CONTAINER_ENGINE) build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(SHORT_COMMIT_HASH) \
		-t $(AGENT_IMAGE):latest -f $(AGENT_CONTAINERFILE) .
	@$(CONTAINER_ENGINE) tag $(AGENT_IMAGE):latest $(AGENT_IMAGE):$(SHORT_COMMIT_HASH)
	@echo "Build Successful"

build-pgsql: ## Build the PGSQL container image
	@echo "### [Building PGSQL container image] ###"
	@$(CONTAINER_ENGINE) build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(SHORT_COMMIT_HASH) \
		-t $(PGSQL_IMAGE):latest -f $(PGSQL_CONTAINERFILE) .
	@$(CONTAINER_ENGINE) tag $(PGSQL_IMAGE):latest $(PGSQL_IMAGE):$(SHORT_COMMIT_HASH)
	@echo "Build Successful"


# Development targets
start-dev: ## Start the container-based development environment
	@echo "### [Starting dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml up -d
	@echo "### [Running dev environment] ###"
	@echo "### [API: http://localhost:8081/api/v1/healthcheck ] ###"

stop-dev: ## Stop the container-based development environment
	@echo "### [Stopping dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml down
	@# If there are no containers attached to the network, remove it
	@[ "$(shell $(CONTAINER_ENGINE) ps --all --filter network=$(COMPOSE_NETWORK) | tail -n +2 | wc -l)" -eq "0" ] && { $(CONTAINER_ENGINE) network rm $(COMPOSE_NETWORK); } || { exit 0; }

restart-dev: ## Restart the container-based development environment
restart-dev: stop-dev start-dev


# Tests targets
go-setup-tests:
	@[ -d $(TEST_DIR) ] || mkdir $(TEST_DIR)

go-unit-tests: ## Runs go unit tests
go-unit-tests: go-setup-tests
	@$(GO) test -v -race ./internal/inventory -coverprofile $(TEST_DIR)/cover-unit-tests.out | sed -e 's/PASS/\x1b[32mPASS\x1b[0m/' -e 's/FAIL/\x1b[31mFAIL\x1b[0m/' -e 's/RUN/\x1b[33mRUN\x1b[0m/'
	@$(GO) tool cover -func $(TEST_DIR)/cover-unit-tests.out

go-integration-tests: ## Runs the Integration tests for this project
go-integration-tests: go-setup-tests
	@echo -e "### [Running Integration tests] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-integration-tests.yaml down
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-integration-tests.yaml up -d
	@sleep 10 # let init-pgsql to do its job
	@$(GO) test -v -race $(TEST_DIR)/integration -coverprofile $(TEST_DIR)/cover-integration-tests.out | sed -e 's/PASS/\x1b[32mPASS\x1b[0m/' -e 's/FAIL/\x1b[31mFAIL\x1b[0m/' -e 's/RUN/\x1b[33mRUN\x1b[0m/'
	@$(GO) tool cover -func $(TEST_DIR)/cover-integration-tests.out
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-integration-tests.yaml down

go-tests: ## Runs every test
go-tests: go-unit-tests go-integration-tests

lint: ## Runs go linter tools against the whole project
	@$(GO_LINTER) run

lint-staged: ## Runs go linter tools against staged files
	@echo "### [Running Linter against staged files] ###"
	@STAGED_FILES=$$(git diff --name-only --staged -- '*.go' ':(exclude)*.pb.go'); \
	if [ -z "$$STAGED_FILES" ]; then \
		echo "No staged Go files to lint."; \
	else \
		$(GO_LINTER) run --new-from-patch=<(git diff --staged -- '*.go' ':(exclude)*.pb.go') --whole-files; \
	fi

go-fmt: ## Runs go formatting tools
	@WRONG_LINES="$$($(GO_FMT) -l . | wc -l)"; \
	if [[ $$WRONG_LINES -gt 0 ]]; then \
		echo "The following files are not properly formatted: $$WRONG_LINES"; \
		$(GO_FMT) -d -e . ; \
		exit 1; \
	fi

# Documentation targets
swagger-editor: ## Open web editor for modifying Swagger docs
	@echo "### [Launching Swagger editor] ###"
	@$(CONTAINER_ENGINE) run --rm -p 127.0.0.1:8082:8080 \
		-e SWAGGER_FILE=/tmp/swagger.yaml \
		-v ./cmd/api/docs/swagger.yaml:/tmp/swagger.yaml:Z \
		swaggerapi/swagger-editor
	@echo "Open your browser at http://127.0.0.1:8082"

swagger-doc: ## Generate Swagger documentation for ClusterIQ API
	@echo "### [Generating Swagger Docs] ###"
	@$(SWAGGER) fmt
	@$(SWAGGER) init --generalInfo ./cmd/api/server.go --parseDependency --output ./cmd/api/docs


# Set the default target to "help"
.DEFAULT_GOAL := help
# Help
help: ## Display this help message
	@echo "### [Makefile Help] ###"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[1;36m%-30s\033[0m %s\n", $$1, $$2}'
