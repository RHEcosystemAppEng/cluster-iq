#
# Cluster IQ Makefile
################################################################################

# Global Vars
VERSION := $(shell cat VERSION)
IMAGE_TAG := $(shell git rev-parse --short=7 HEAD)
CONTAINER_ENGINE ?= $(shell which podman >/dev/null 2>&1 && echo podman || echo docker)
K8S_CLI ?= $(shell which oc >/dev/null 2>&1 && echo oc || echo kubectl)
REGISTRY ?= quay.io
PROJECT_NAME ?= cluster-iq
REGISTRY_REPO ?= ecosystem-appeng
API_IMG_NAME ?= $(PROJECT_NAME)-api
API_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/${API_IMG_NAME}
SCANNER_IMG_NAME ?= $(PROJECT_NAME)-aws-scanner
SCANNER_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/${SCANNER_IMG_NAME}
BUILD_DIR ?= ./build
BIN_DIR ?= $(BUILD_DIR)/bin
CMD_DIR ?= ./cmd
PKG_DIR ?= ./internal
DEPLOYMENTS_DIR ?= ./deployments

# Load the .env file and export the variables
include .env
export

# PHONY rules
.PHONY: deploy

# Help message
define HELP_MSG
Makefile Rules:
	deploy: Deploys the application on the current context configured on Openshift/Kubernetes CLI
	clean: Removes local container images
	build: Builds every component it the repo: (API, AWS-Scanner)
	build-api: Builds every component it the repo: (API, AWS-Scanner)
	build-scanners: Builds the scanners for each supported cloud provider
	build-aws-scanner: Builds the scanner for AWS
	local-build: Builds every component it the repo in your local environment: (API, AWS-Scanner)
	local-build-api: Builds every component it the repo in your local environment: (API, AWS-Scanner)
	local-build-scanners: Builds in your local environment the scanners for each supported cloud provider
	local-build-aws-scanner: Builds in your local environment the scanner for AWS
	push: Pushes every container image into remote repo
	push-api: Pushes API container image
	push-scanner: Pushes every supported scanner image
	push-aws-scanner: Pushes AWS scanner image
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
	@$(CONTAINER_ENGINE) images | grep $(PROJECT_NAME) | awk '{print $3}' | xargs $(CONTAINER_ENGINE) rmi -f

build: build-scanners build-api

build-api:
	@echo "### [Building API] ###"
	@$(CONTAINER_ENGINE) build -t $(API_IMAGE):latest -f ./$(DEPLOYMENTS_DIR)/dockerfiles/Dockerfile-api-arg="VERSION=${VERSION}" .
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(IMAGE_TAG)

build-scanners: build-aws-scanner

build-aws-scanner:
	@echo "### [Building AWS Scanner] ###"
	@$(CONTAINER_ENGINE) build -t $(SCANNER_IMAGE):latest -f ./$(DEPLOYMENTS_DIR)/dockerfiles/Dockerfile-scanner-aws --build-arg="VERSION=${VERSION}" .
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(IMAGE_TAG)


# Building in local environment
local-clean:
	@echo "### [Cleanning local building] ###"
	@rm -Rf $(BIN_DIR)

local-build: local-build-scanners local-build-api

local-build-api:
	@echo "### [Building API] ###"
	@cd ./cmd/api/ ; go build -o ../../$(BIN_DIR)/api/api ; cd ../..

local-build-scanners: local-build-api local-build-aws-scanner

local-build-aws-scanner:
	@echo "### [Building AWS Scanner] ###"
	@cd ./cmd/api/ ; go build -o ../../$(BIN_DIR)/scanners/aws_scanner ; cd ../..


# Publish images
push: push-api push-scanner

push-api:
	$(CONTAINER_ENGINE) push $(API_IMAGE):latest
	$(CONTAINER_ENGINE) push $(API_IMAGE):$(VERSION)
	$(CONTAINER_ENGINE) push $(API_IMAGE):$(IMAGE_TAG)

push-scanners: push-aws-scanner

push-aws-scanner:
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

help:
	@echo "$$HELP_MSG"
