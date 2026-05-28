# Development Setup

This guide describes how to build and deploy [ClusterIQ](https://github.com/RHEcosystemAppEng/cluster-iq) in a development environment. The setup uses container compose files and is intended for development purposes only.

ClusterIQ is a monorepo containing both the backend (Go) and the web console (React/TypeScript) under the `console/` directory.

## Prerequisites

Before you begin:

* Ensure you have the necessary cloud account credentials.
* Ensure you have access to `registry.redhat.io` to download the required container images.
* If you experience file mounting issues with local files (such as `init.psql` or `credentials`), verify your SELinux settings. SELinux in enforcing mode can prevent container runtime from binding files to containers.

To temporarily disable SELinux:

```sh
sudo setenforce 0
```

> [!NOTE] Use this command with caution and only in development environments.

## Build dependencies

* [Go v1.25](https://go.dev/dl/)
* [Node.js 18.x](https://nodejs.org/) and npm
* [podman](https://podman.io/docs/installation) or [docker](https://docs.docker.com/engine/install)
* [podman-compose](https://github.com/containers/podman-compose?tab=readme-ov-file#installation) or [docker-compose](https://docs.docker.com/compose/install/)
* [swag](https://github.com/swaggo/swag?tab=readme-ov-file#getting-started)

## Build

Follow these steps to build the ClusterIQ components:

1. Clone the repository:

    ```sh
    git clone git@github.com:RHEcosystemAppEng/cluster-iq.git
    cd cluster-iq
    ```

2. Validate required dependencies:

    If you encounter an error, please ensure that you have installed all the necessary dependencies before proceeding.

    ```sh
    make check-dependencies
    ```

3. Build the container images (backend + console):

    ```sh
    make build
    ```

4. Verify the container images:

   You should see `cluster-iq-api`, `cluster-iq-scanner`, `cluster-iq-agent`, `cluster-iq-pgsql`, and `cluster-iq-console`.

    ```sh
    CONTAINER_ENGINE=$(which podman >/dev/null 2>&1 && echo podman || echo docker)
    $CONTAINER_ENGINE images | grep cluster-iq
    ```

## Deployment

To manage your development environment:

1. Configure your [cloud account credentials](../../README.md#accounts-configuration).

2. Start the environment:

    ```sh
    make start-dev
    ```

    This starts all services (API, Scanner, Agent, Console, PostgreSQL) via compose.
    - API: http://localhost:8081/api/v1/healthcheck
    - Console: http://localhost:8080

3. Stop the environment:

    ```sh
    make stop-dev
    ```

## Console Development

For working on the console frontend locally (with hot-reload):

```sh
make console-install     # Install npm dependencies
make console-start-dev   # Start Vite dev server (port 3000, proxies API to localhost:8081)
make console-lint        # Run prettier + eslint + tsc
```

See `console/README.md` for more details.

## API Documentation

### Generating Swagger Documentation

To generate the API documentation from the source code:

```sh
make swagger-doc
```

### Running Swagger Editor

To view and edit the OpenAPI specification in a browser-based editor:

```sh
make swagger-editor
```
