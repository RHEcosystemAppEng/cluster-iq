# Helm Chart installation guide

This guide provides step-by-step instructions for installing the `cluster-iq` Helm chart from the repository directory.

## Directory Structure

Here's a simplified structure of the `cluster-iq` Helm chart:

```text
deployments/helm/cluster-iq/
├── charts
├── Chart.yaml
├── README.md
├── templates
│   ├── agent
│   ├── api
│   ├── console
│   ├── database
│   ├── _helpers.tpl
│   └── scanner
└── values.yaml
```

## Prerequisites

Before installing the chart, ensure the following resources are created in your OpenShift cluster:

### OpenShift project

```bash
oc new-project cluster-iq
```

### Cloud Providers credentials file

The credentials file required for scraping cloud resources by the scanner.
The secret name must be `credentials` and stored as a OpenShift secret.

Example of a file with credentials

```text
[ACCOUNT_NAME]
provider = {aws/gcp/azure}
user = XXXXXXX
key = YYYYYYY
billing_enabled = {true/false}
```

### ImagePullSecrets for the database

- The `database.imagePullSecrets` value in `values.yaml` must point to a pre-created OpenShift secret. This secret contains the credentials required to pull the database image.

Please refer to the links provided for more information.

- [Red Hat Ecosystem Catalog](https://catalog.redhat.com/software/containers/rhel8/postgresql-12/5db133bd5a13461646df330b?container-tabs=gti&gti-tabs=red-hat-login)
- [Registry Service Accounts](https://access.redhat.com/terms-based-registry)
- [OpenShift documentation](https://docs.openshift.com/container-platform/4.17/openshift_images/managing_images/using-image-pull-secrets.html#images-allow-pods-to-reference-images-from-secure-registries_using-image-pull-secrets)

Update the `database.imagePullSecrets` in `values.yaml`

```yaml
database:
  imagePullSecrets:
    - name: my-database-pull-secret
```

## Installation Steps

1. Navigate to the Chart directory

    ```bash
    cd deployments/helm/cluster-iq
    ```

2. Review and update `values.yaml` if needed
3. Install the Chart

    ```bash
    helm upgrade --install cluster-iq . --namespace cluster-iq -f values.yaml
    ```

4. Verify the Installation

    ```bash
    oc get all -n cluster-iq
    ```

5. To avoid waiting for the first launch of the scanner, you can start the job manually

    ```bash
    oc create job --from=cronjob/scanner scanner-init -n cluster-iq
    ```

## Uninstallation

To remove the `cluster-iq` chart and all associated resources:

  ```bash
  helm uninstall cluster-iq --namespace cluster-iq
  ```

This will clean up all resources created by the chart.
