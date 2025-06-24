#
# ClusterIQ Makefile
################################################################################

# Global Vars
SHORT_COMMIT_HASH := $(shell git rev-parse --short=7 HEAD)

# Binary vars
CONTAINER_ENGINE ?= $(shell which podman >/dev/null 2>&1 && echo podman || echo docker)
K8S_CLI ?= $(shell which oc >/dev/null 2>&1 && echo oc || echo kubectl)

# Required binaries
REQUIRED_BINS := $(CONTAINER_ENGINE) $(CONTAINER_ENGINE)-compose $(K8S_CLI) swag protoc

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

# Standard targets
all: ## Stop, build and start the development environment based on containers
all: stop-dev build start-dev

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
	@go build -o $(BIN_DIR)/api/api $(LDFLAGS) ./cmd/api/

local-build-scanner: ## Build the scanner binary
	@echo "### [Building Scanner] ###"
	@go build -o $(BIN_DIR)/scanners/scanner $(LDFLAGS) ./cmd/scanner

local-build-agent: ## Build the agent binary
	@echo "### [Building Agent] ###"
	@[ ! -d $(GENERATED_DIR) ] && { mkdir $(GENERATED_DIR); } || { exit 0; }
	@protoc --go_out=$(GENERATED_DIR) --go-grpc_out=$(GENERATED_DIR) $(AGENT_PROTO_PATH)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/agent/agent $(LDFLAGS) ./cmd/agent


# Container based working targets
clean: ## Remove the container images
	@echo "### [Cleaning Container images] ###"
	@-$(CONTAINER_ENGINE) images | grep -e $(SCANNER_IMAGE) -e $(API_IMAGE) -e $(AGENT_IMAGE) | awk '{print $$3}' | xargs $(CONTAINER_ENGINE) rmi -f

build: build-api build-scanner build-agent ## Build all container images
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
.PHONY: test
tests_unit_tests: ## Runs the Unit tests for this project internal packages
	@[[ -d $(TEST_DIR) ]] || mkdir $(TEST_DIR)
	@go test -v -race ./internal/inventory -coverprofile $(TEST_DIR)/cover-unit-tests.out

tests_integration_tests: ## Runs the Integration tests for this project
tests_integration_tests: restart-dev
	@podman stop scanner
	@echo -e "\n\n### [Running Integration tests] ###"
	@go test -v -race ./test/integration -coverprofile $(TEST_DIR)/cover-integration-tests.out

tests: ## Runs every test
tests: tests_unit_tests tests_integration_tests

tests_cover: ## Runs every test and reports about the coverage percentage
tests_cover: tests
	@go tool cover -func $(TEST_DIR)/cover-unit-tests.out
	@go tool cover -func $(TEST_DIR)/cover-integration-tests.out


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
	@swag fmt
	@swag init --generalInfo ./cmd/api/server.go --parseDependency --output ./cmd/api/docs


# Set the default target to "help"
.DEFAULT_GOAL := help
# Help
help: ## Display this help message
	@echo "### [Makefile Help] ###"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[1;36m%-30s\033[0m %s\n", $$1, $$2}'
