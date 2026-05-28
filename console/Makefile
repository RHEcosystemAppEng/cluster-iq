#
# ClusterIQ Console Makefile
#
# This Makefile handles console-specific targets (npm/TypeScript).
# Container image builds and compose orchestration are managed by the root Makefile.
# From the repo root, use: make console-install, make console-build, make console-lint, etc.
################################################################################

# Standard targets
all: ## Installs dependencies and builds locally
all: local-install local-build

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
