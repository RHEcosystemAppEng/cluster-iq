# Development Setup

This guide describes how to build and deploy [ClusterIQ](https://github.com/RHEcosystemAppEng/cluster-iq) in a development environment. The setup uses container compose files and is intended for development purposes only.

ClusterIQ consists of two repositories:

* [Console Repo](https://github.com/RHEcosystemAppEng/cluster-iq-console) contains the web user interface.
* [Main Repo](https://github.com/RHEcosystemAppEng/cluster-iq-console) contains the API and Scanner components.

Each repository requires separate configuration and management.

## Prerequisites

Before you begin:

* Ensure you have the necessary cloud account credentials.
* Ensure you have access to `registry.redhat.io` to download the required container images.
* If you experience file mounting issues with local files (such as `init.psql` or `credentials`), verify your SELinux settings. SELinux in enforcing mode can prevent container runtime from binding files to containers.

To temporarily disable SELinux:

```sh
sudo setenforce 0
```

[!NOTE] Use this command with caution and only in development environments.

## Build dependencies

* [go v1.19](https://go.dev/dl/)
* [podman](https://podman.io/docs/installation) or [docker](https://docs.docker.com/engine/install)
* [podman-compose](https://github.com/containers/podman-compose?tab=readme-ov-file#installation) or [docker-compose](https://docs.docker.com/compose/install/)
* [swag](https://github.com/swaggo/swag?tab=readme-ov-file#getting-started)

## Build

Follow these steps to build the ClusterIQ components:

1. Create and navigate to a common folder for both repos:

    ```sh
    WORKDIR=$(pwd)/cluster-iq-repos
    mkdir -p $WORKDIR && cd $WORKDIR
    ```

2. Clone the repositories:

    ```sh
    git clone git@github.com:RHEcosystemAppEng/cluster-iq.git
    git clone git@github.com:RHEcosystemAppEng/cluster-iq-console.git
    ```

3. Validate required dependencies:

    If you encounter an error, please ensure that you have installed all the necessary dependencies before proceeding.

    ```sh
    cd ${WORKDIR}/cluster-iq
    make check-dependencies
    ```

4. Build the container images:

    ```sh
    git checkout main
    make build
    ```

    ```sh
    cd ${WORKDIR}/cluster-iq-console
    git checkout main
    make build
    ```

5. Verify the container images:

   You should see the following images `cluster-iq-api`, `cluster-iq-scanner`, `cluster-iq-console`

    ```sh
    CONTAINER_ENGINE=$(which podman >/dev/null 2>&1 && echo podman || echo docker)
    $CONTAINER_ENGINE images | grep cluster-iq
    ```

## Deployment

To manage your development environment:

1. Change the working directory to `cluster-iq` repo

   ```sh
   cd ${WORKDIR}/cluster-iq
   ```

2. Configure your [cloud account credentials](../README.md#accounts-configuration).
3. Start the environment:

    ```sh
    make start-dev
    ```

4. Stop the environment:

    ```sh
    make stop-dev
    ```

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
