VERSION := $(shell cat VERSION)
IMAGE_TAG := $(shell git rev-parse --short=7 HEAD)
CONTAINER_ENGINE ?= $(shell which podman >/dev/null 2>&1 && echo podman || echo docker)
REGISTRY ?= quay.io
PROJECT_NAME ?= cluster-iq
REGISTRY_REPO ?= ecosystem-appeng
API_IMG_NAME ?= $(PROJECT_NAME)-api
API_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/${API_IMG_NAME}
SCANNER_IMG_NAME ?= $(PROJECT_NAME)-aws-scanner
SCANNER_IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/${SCANNER_IMG_NAME}

# Load the .env file and export the variables
include .env
export

# Help message
define HELP_MSG
Makefile Rules:
	clean: Removes local container images
	build: Builds every component it the repo: (API, AWS-Scanner)
	build-api: Builds every component it the repo: (API, AWS-Scanner)
	build-scanners: Builds the scanners for each supported cloud provider
	build-aws-scanner: Builds the scanner for AWS
	push: Pushes every container image into remote repo
	push-api: Pushes API container image
	push-scanner: Pushes every supported scanner image
	push-aws-scanner: Pushes AWS scanner image
	start-dev: Starts a local environment using 'docker/podman-compose'
	stop-dev: Stops the local environment using 'docker/podman-compose'
	help: Displays this message
endef
export HELP_MSG

clean:
	@echo "### [Cleanning Docker images] ###"
	@$(CONTAINER_ENGINE) images | grep $(PROJECT_NAME) | awk '{print $3}' | xargs $(CONTAINER_ENGINE) rmi -f

build: build-scanners build-api

build-api:
	@echo "### [Building API] ###"
	@$(CONTAINER_ENGINE) build -t $(API_IMAGE):latest -f src/api/Dockerfile --build-arg="VERSION=${VERSION}" .
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) tag $(API_IMAGE):latest $(API_IMAGE):$(IMAGE_TAG)

build-scanners: build-aws-scanner

build-aws-scanner:
	@echo "### [Building AWS Scanner] ###"
	@$(CONTAINER_ENGINE) build -t $(SCANNER_IMAGE):latest -f src/aws_scanner/Dockerfile --build-arg="VERSION=${VERSION}" .
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(VERSION)
	@$(CONTAINER_ENGINE) tag $(SCANNER_IMAGE):latest $(SCANNER_IMAGE):$(IMAGE_TAG)

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

start-dev: build
	@echo "### [Starting dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f manifests/docker-compose.yaml up &

stop-dev:
	@echo "### [Stopping dev environment] ###"
	@$(CONTAINER_ENGINE)-compose -f manifests/docker-compose.yaml down

help:
	@echo "$$HELP_MSG"
