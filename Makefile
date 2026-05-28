#
# ClusterIQ Console Makefile
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

# Project directories
DEPLOYMENTS_DIR ?= ./deployments

# Images
CONSOLE_IMG_NAME ?= $(PROJECT_NAME)-console
CONSOLE_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/$(CONSOLE_IMG_NAME)


# Standard targets
all: ## Installs and run the Console on Local
all: install build-local start-dev

# Local working targets
local-clean: ## Cleans project built (./dist)
	@echo "### [Cleanning local building] ###"
	@npm run clean

local-install: ## Installs project dependencies in local
	@echo "### [Installing dependencies] ###"
	@npm install

local-build: ## Builds the project locally
	@echo "### [Building project] ###"
	@npm run build

local-start-dev: export VITE_CIQ_API_URL = http://localhost:8081
local-start-dev: ## Starts the project in development mode
	@echo "### [Starting project DEV MODE] ###"
	@echo "### [API-URL: $$VITE_CIQ_API_URL] ###"
	@npm run start


# Container based working targets
clean: ## Removes the container image for the Console
	@echo "### [Cleanning Container images] ###"
	@$(CONTAINER_ENGINE) images | grep $(CONSOLE_IMAGE) | awk '{print $$3}' | xargs $(CONTAINER_ENGINE) rmi -f

build: ## Builds the Console container image
	@echo "### [Building Console container image] ###"
	@$(CONTAINER_ENGINE) build -t $(CONSOLE_IMAGE):latest -f $(DEPLOYMENTS_DIR)/containerfiles/Containerfile .
	@$(CONTAINER_ENGINE) tag $(CONSOLE_IMAGE):latest $(CONSOLE_IMAGE):$(SHORT_COMMIT_HASH)
	@echo "Build Successful"


# Development targets
start-dev: ## Development env based on Docker/Podman Compose tool
	@echo "### [Starting dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml up --build -d
	@echo "### [Console: http://localhost:8080] ###"

stop-dev: ## Stops the container based development env
	@echo "### [Stopping dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml down

restart-dev: ## Restarts the container based env
restart-dev: stop-dev start-dev

ts-prettier: ## Runs code prettier
	@echo "### [Running Prettier] ###"
	@npx prettier --log-level=warn --check ./src

ts-prettier-fix: ## Runs code prettier fixing
	@echo "### [Running Prettier] ###"
	@npx prettier --log-level=warn --check --write ./src && git diff ./src

ts-eslint: ## Runs Linter
	@echo "### [Running EsLinter] ###"
	@npx eslint .

ts-tsc: ## Runs Typescript type test
	@echo "### [Running TSC test] ###"
	@npx tsc --noEmit

# Tests targets
ts-test: ## Runs code tests and coverage
ts-test: ts-prettier ts-eslint ts-tsc


# Set the default target to "help"
.DEFAULT_GOAL := help
# Help
help: ## Display this help message
	@echo "### [Makefile Help] ###"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[1;36m%-30s\033[0m %s\n", $$1, $$2}'
