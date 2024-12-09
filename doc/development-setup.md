# Development Setup
This document explains how to build and deploy
[ClusterIQ](https://github.com/RHEcosystemAppEng/cluster-iq) for development
purposes. For creating a dev environment, this project uses compose files. Not
recommended for production

ClusterIQ has two different repos:
* [Console Repo](https://github.com/RHEcosystemAppEng/cluster-iq-console) for the Web User Interface
* [Main Repo](https://github.com/RHEcosystemAppEng/cluster-iq-console) containing the API and the Scanner.


Both repos must be managed by separate.

## Prerequisites
Check you have configured your [cloud account
credentials](../README.md#accounts-configuration) before continuing:

:warning: Make sure you have access to `registry.redhat.io` for downloading
required images.

:warning: If you're having issues mounting your local files (like init.psql or
the credentials file) check if your SELinux is enforcing. This could prevent
podman to bind these files into the containers.
```sh
# Run this carefully and under your own responsability
sudo setenforce 0
```

## Building
First, we need to build the images for every component.

1. Create a common folder for both repos
```sh
mkdir cluster-iq
```

2. Download both repos
```sh
#Clonning Console Repo
git clone git@github.com:RHEcosystemAppEng/cluster-iq.git
```
```sh
#Clonning Console Repo
git clone git@github.com:RHEcosystemAppEng/cluster-iq-console.git
```

3. Build container images
```sh
cd cluster-iq
## OPTIONAL: switch to another branch version if you need it, but we recomend to use `main`.
git checkout main
make build
cd ..
```
```sh
cd cluster-iq-console
## OPTIONAL: switch to another branch version if you need it, but we recomend to use `main`.
git checkout main
make build
cd ..
```

4. Verify all the images are correctly built:
```sh
# You should see three images, cluster-iq-api, cluster-iq-scanner, cluster-iq-console
podman images | grep cluster-iq
```

## Deployment
1. Use Compose files for deploying the platform
```sh
make start-dev
```

2. Stop dev environment
```sh
make stop-dev
```
