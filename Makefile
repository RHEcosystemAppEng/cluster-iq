#
# ClusterIQ Makefile
################################################################################

# Global Vars
SHORT_COMMIT_HASH := $(shell git rev-parse --short=7 HEAD)

# Binary vars
CONTAINER_ENGINE ?= $(shell which podman >/dev/null 2>&1 && echo podman || echo docker)
K8S_CLI ?= $(shell which oc >/dev/null 2>&1 && echo oc || echo kubectl)

# Container image registy vars
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

# Images
API_IMG_NAME ?= $(PROJECT_NAME)-api
API_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/$(API_IMG_NAME)
API_CONTAINERFILE ?= ./$(DEPLOYMENTS_DIR)/containerfiles/Containerfile-api
SCANNER_IMG_NAME ?= $(PROJECT_NAME)-scanner
SCANNER_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/$(SCANNER_IMG_NAME)
SCANNER_CONTAINERFILE ?= ./$(DEPLOYMENTS_DIR)/containerfiles/Containerfile-scanner

# Standard targets
all: ## Stops, build and starts the development environment based on containers
all: stop-dev build start-dev


# Local working targets
local-clean:
	@echo "### [Cleanning local building] ###"
	@rm -Rf $(BIN_DIR)

local-build: local-build-scanner local-build-api

local-build-api: swagger-doc
	@echo "### [Building API] ###"
	@go build -o $(BIN_DIR)/api/api $(LDFLAGS) ./cmd/api/

local-build-scanner:
	@echo "### [Building Scanner] ###"
	@go build -o $(BIN_DIR)/scanners/scanner $(LDFLAGS) ./cmd/scanner


# Container based working targets
clean: ## Removes the container images for the API and the Scanner
	@echo "### [Cleanning Container images] ###"
	@$(CONTAINER_ENGINE) images | grep $(REGISTRY_REPO)/$(PROJECT_NAME) | awk '{print $$3}' | xargs $(CONTAINER_ENGINE) rmi -f

build: ## Builds the container images for the API and the Scanner
build: build-api build-scanner

build-api: ## Builds the API Container image
	@echo "### [Building API container image] ###"
	@$(CONTAINER_ENGINE) build -t $(API_IMAGE):latest -f $(API_CONTAINERFILE) .
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(SHORT_COMMIT_HASH)
	@echo "Build Successful"

build-scanner: ## Builds the Scanner Container image
	@echo "### [Building Scanner container image] ###"
	@$(CONTAINER_ENGINE) build -t $(SCANNER_IMAGE):latest -f $(SCANNER_CONTAINERFILE) .
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(SHORT_COMMIT_HASH)
	@echo "Build Successful"


# Development targets
start-dev: ## Development env based on Docker/Podman Compose tool
	@echo "### [Starting dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml up -d
	@echo "### [Dev environment running] ###"
	@echo "### [API: http://localhost:8081/api/v1/healthcheck ] ###"

stop-dev: ## Stops the container based development env
	@echo "### [Stopping dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml down
	@# If there are no containers attached to the network, remove it
	@[ "$(shell $(CONTAINER_ENGINE) ps --all --filter network=$(COMPOSE_NETWORK) | tail -n +2 | wc -l)" -eq "0" ] && { $(CONTAINER_ENGINE) network rm $(COMPOSE_NETWORK); }

restart-dev: ## Restarts the container based development env
restart-dev: stop-dev start-dev


# Tests targets
test:
	@[[ -d $(TEST_DIR) ]] || mkdir $(TEST_DIR)
	@go test -race ./... -coverprofile $(TEST_DIR)/cover.out

cover: test
	@go tool cover -func $(TEST_DIR)/cover.out


# Documentation targets
swagger-editor: ## Opens web editor for modifying Swagger docs
	@echo "### [Launching Swagger editor] ###"
	@$(CONTAINER_ENGINE) run --rm -p 127.0.0.1:8082:8080 \
		-e SWAGGER_FILE=/tmp/swagger.yaml \
		-v ./cmd/api/docs/swagger.yaml:/tmp/swagger.yaml:Z \
		swaggerapi/swagger-editor
	@echo "Open your browser at http://127.0.0.1:8082"

swagger-doc: # Generates Swagger documentation for ClusterIQ API
	@echo "### [Generating Swagger Docs] ###"
	@swag fmt
	@swag init --generalInfo ./cmd/api/api_server.go --parseDependency --output ./cmd/api/docs


# Set the default target to "help"
.DEFAULT_GOAL := help
# Help
help: ## Display this help message
	@echo "### [Makefile Help] ###"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[1;36m%-30s\033[0m %s\n", $$1, $$2}'
