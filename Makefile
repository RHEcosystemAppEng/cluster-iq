#
# Cluster IQ Makefile
################################################################################

# Global Vars
#
VERSION := $(shell cat VERSION)
IMAGE_TAG := $(shell git rev-parse --short=7 HEAD)

# Binary vars
CONTAINER_ENGINE ?= $(shell which podman >/dev/null 2>&1 && echo podman || echo docker)
K8S_CLI ?= $(shell which oc >/dev/null 2>&1 && echo oc || echo kubectl)

# Container images vars
REGISTRY ?= quay.io
PROJECT_NAME ?= cluster-iq
REGISTRY_REPO ?= ecosystem-appeng
API_IMG_NAME ?= $(PROJECT_NAME)-api
API_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/${API_IMG_NAME}
SCANNER_IMG_NAME ?= $(PROJECT_NAME)-scanner
SCANNER_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/${SCANNER_IMG_NAME}

# Building vars
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(IMAGE_TAG)"

# Project directories
TEST_DIR ?= ./test
BUILD_DIR ?= ./build
BIN_DIR ?= $(BUILD_DIR)/bin
CMD_DIR ?= ./cmd
PKG_DIR ?= ./internal
DEPLOYMENTS_DIR ?= ./deployments

# Load the .env file and export the variables
include .env
export

# PHONY rules
.PHONY: deploy test

# Help message
define HELP_MSG
Makefile Rules:
	deploy: Deploys the application on the current context configured on Openshift/Kubernetes CLI
	clean: Removes local container images
	build: Builds every component image in the repo: (API, scanner)
	build-api: Builds API container image
	build-scanner: Builds the cluster-iq scanner container image
	local-build: Builds every component it the repo in your local environment: (API, scanner)
	local-build-api: Builds API binary in your local environment.
	local-build-scanner: Builds in your local environment the cluster-iq scanner
	push: Pushes every container image into remote repo
	push-api: Pushes API container image
	push-scanner: Pushes cluster-iq-scanner container image
	start-dev: Starts a local environment using 'docker/podman-compose'
	stop-dev: Stops the local environment using 'docker/podman-compose'
	help: Displays this message
endef
export HELP_MSG


# Deployments
deploy:
	@$(K8S_CLI) apply -f $(DEPLOYMENTS_DIR)/openshift


# Building using containers:
clean:
	@echo "### [Cleanning Docker images] ###"
	@$(CONTAINER_ENGINE) images | grep $(PROJECT_NAME) | awk '{print $$3}' | xargs $(CONTAINER_ENGINE) rmi -f

build: build-scanner build-api

build-api:
	@echo "### [Building API] ###"
	@$(CONTAINER_ENGINE) build -t $(API_IMAGE):latest -f ./$(DEPLOYMENTS_DIR)/dockerfiles/Dockerfile-api --build-arg="VERSION=${VERSION}" --build-arg="COMMIT=${IMAGE_TAG}" .
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(IMAGE_TAG)

build-scanner:
	@echo "### [Building Scanner] ###"
	@$(CONTAINER_ENGINE) build -t $(SCANNER_IMAGE):latest -f ./$(DEPLOYMENTS_DIR)/dockerfiles/Dockerfile-scanner --build-arg="VERSION=${VERSION}" --build-arg="COMMIT=${IMAGE_TAG}" .
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(IMAGE_TAG)


# Building in local environment
local-clean:
	@echo "### [Cleanning local building] ###"
	@rm -Rf $(BIN_DIR)

local-build: local-build-scanner local-build-api

local-build-api:
	@echo "### [Building API] ###"
	@go build -o $(BIN_DIR)/api/api $(LDFLAGS) ./cmd/api/

local-build-scanner:
	@echo "### [Building Scanner] ###"
	@go build -o $(BIN_DIR)/scanners/scanner $(LDFLAGS) ./cmd/scanner

# Publish images
push: push-api push-scanner

push-api:
	$(CONTAINER_ENGINE) push $(API_IMAGE):latest
	$(CONTAINER_ENGINE) push $(API_IMAGE):$(VERSION)
	$(CONTAINER_ENGINE) push $(API_IMAGE):$(IMAGE_TAG)

push-scanner:
	@$(CONTAINER_ENGINE) push $(SCANNER_IMAGE):latest
	@$(CONTAINER_ENGINE) push $(SCANNER_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) push $(SCANNER_IMAGE):$(IMAGE_TAG)


# Development env based on Docker/Podman Compose tool
start-dev:
	@echo "### [Starting dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/docker-compose/docker-compose.yaml up &

stop-dev:
	@echo "### [Stopping dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/docker-compose/docker-compose.yaml down


# Tests
test:
	@[[ -d $(TEST_DIR) ]] || mkdir $(TEST_DIR)
	@go test -race ./... -coverprofile $(TEST_DIR)/cover.out

cover:
	@go tool cover -func $(TEST_DIR)/cover.out


# Help
# Set the default target to "help"
.DEFAULT_GOAL := help
help:
	@echo "$$HELP_MSG"
