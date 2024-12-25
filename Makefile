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
\033[1;37mMakefile Rules\033[0m:
	\033[1;36mall:\033[0m                 \033[0;37m Stops the devel env, re-build the images, and starts the devel env again
	\033[1;36mdeploy:\033[0m              \033[0;37m Deploys the application on the current context configured on Openshift/Kubernetes CLI
	\033[1;36mclean:\033[0m               \033[0;37m Removes local container images
	\033[1;36mbuild:\033[0m               \033[0;37m Builds every component image in the repo:\033[0m \033[0;37m (API, scanner)
	\033[1;36mbuild-api:\033[0m           \033[0;37m Builds API container image
	\033[1;36mbuild-scanner:\033[0m       \033[0;37m Builds the cluster-iq scanner container image
	\033[1;36mlocal-build:\033[0m         \033[0;37m Builds every component it the repo in your local environment:\033[0m \033[0;37m (API, scanner)
	\033[1;36mlocal-build-api:\033[0m     \033[0;37m Builds API binary in your local environment.
	\033[1;36mlocal-build-scanner:\033[0m \033[0;37m Builds in your local environment the cluster-iq scanner
	\033[1;36mpush:\033[0m                \033[0;37m Pushes every container image into remote repo
	\033[1;36mpush-api:\033[0m            \033[0;37m Pushes API container image
	\033[1;36mpush-scanner:\033[0m        \033[0;37m Pushes cluster-iq-scanner container image
	\033[1;36mstart-dev:\033[0m           \033[0;37m Starts a local environment using 'docker/podman-compose' and initializes the Database with some fake data
	\033[1;36mdeploy-compose:\033[0m      \033[0;37m Starts a local environment using 'docker/podman-compose'
	\033[1;36mstop-dev:\033[0m            \033[0;37m Stops the local environment using 'docker/podman-compose'
	\033[1;36mswagger-editor:\033[0m      \033[0;37m Starts Swagger Editor using a docker container
	\033[1;36mswagger-doc:\033[0m         \033[0;37m generates Swagger Documentation of the API
	\033[1;36mhelp:\033[0m                \033[0;37m Displays this message
	\033[0m
endef
export HELP_MSG


all: stop-dev build start-dev

# Deployments
deploy:
	@$(K8S_CLI) apply -f $(DEPLOYMENTS_DIR)/openshift


# Building using containers:
clean:
	@echo "### [Cleanning Docker images] ###"
	@$(CONTAINER_ENGINE) images | grep $(PROJECT_NAME) | awk '{print $$3}' | xargs $(CONTAINER_ENGINE) rmi -f

build: build-scanner swagger-doc build-api

build-api:
	@echo "### [Building API] ###"
	@$(CONTAINER_ENGINE) build -t $(API_IMAGE):latest -f ./$(DEPLOYMENTS_DIR)/containerfiles/Containerfile-api .
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(IMAGE_TAG)

build-scanner:
	@echo "### [Building Scanner] ###"
	@$(CONTAINER_ENGINE) build -t $(SCANNER_IMAGE):latest -f ./$(DEPLOYMENTS_DIR)/containerfiles/Containerfile-scanner .
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(IMAGE_TAG)


# Building in local environment
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

# Publish images
push: push-api push-scanner

push-api:
	$(CONTAINER_ENGINE) push $(API_IMAGE):latest
	$(CONTAINER_ENGINE) push $(API_IMAGE):$(VERSION)
	$(CONTAINER_ENGINE) push $(API_IMAGE):$(IMAGE_TAG)

push-scanner:
	$(CONTAINER_ENGINE) push $(SCANNER_IMAGE):latest
	$(CONTAINER_ENGINE) push $(SCANNER_IMAGE):$(VERSION)
	$(CONTAINER_ENGINE) push $(SCANNER_IMAGE):$(IMAGE_TAG)


# Development env based on Docker/Podman Compose tool
start-dev:
	@echo "### [Starting dev environment] ###" 
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml pull
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml up -d

stop-dev:
	@echo "### [Stopping dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f $(DEPLOYMENTS_DIR)/compose/compose-devel.yaml down

restart-dev: stop-dev start-dev

init-psql-test-data:
	export PGPASSWORD=password ; psql -h localhost -p 5432 -U user -d clusteriq < db/sql/test_data.sql


# Tests
test:
	@[[ -d $(TEST_DIR) ]] || mkdir $(TEST_DIR)
	@go test -race ./... -coverprofile $(TEST_DIR)/cover.out

cover: test
	@go tool cover -func $(TEST_DIR)/cover.out

# Swagger
swagger-editor:
	@echo "Open your browser at http://127.0.0.1:8082"
	@$(CONTAINER_ENGINE) run --rm -p 127.0.0.1:8082:8080 \
		-e SWAGGER_FILE=/tmp/swagger.yaml \
		-v ./cmd/api/docs/swagger.yaml:/tmp/swagger.yaml:Z \
		swaggerapi/swagger-editor

swagger-doc:
	@echo "### [Generating Swagger Docs] ###"
	@swag fmt
	@swag init --generalInfo ./cmd/api/api_server.go --parseDependency --output ./cmd/api/docs




# Help
# Set the default target to "help"
.DEFAULT_GOAL := help
help:
	@echo -e "$$HELP_MSG"
