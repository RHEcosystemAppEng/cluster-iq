# New version preparation guide
This document details the steps for reviewing, preparing and releasing a new
ClusterIQ version with a standard procedure

## Procedure
In order to prepare and release a new version of ClusterIQ follow this
checklist:

* [ ] Checkout the `release-X.Y.Z` branch you want to review.
  ```sh
  git checkout release-X.Y.Z
  git pull
  ```
* [ ] Create the PR with `main` as base with `Release <VERSION>` as name.
* [ ] Ensure every PR involved on the release was correctly merged.
* [ ] Check the current branch passes `tests` and `linterns`.
  ```sh
  make go-tests
  ```
  **DO NOT CONTINUE** if there's any errors
* [ ] Update `./VERSION` file with the version number which corresponds to the release.
* [ ] Check if the Helm Template for the initialization of the DB. The files
       `./db/sql/init.sql` and `./db/sql/cron.sql` content must be updated in
       `./deployments/helm/cluster-iq/templates/database/configmap-init.yaml`
  ```sh
  # Compare and update the configmap
  vim db/sql/init.sql deployments/helm/cluster-iq/templates/database/configmap-init.yaml -O
  vim db/sql/cron.sql deployments/helm/cluster-iq/templates/database/configmap-init.yaml -O
  ```
* [ ] Update the Helm Chart version and AppVersion if it's needed
  ```sh
  vim ./deployments/helm/cluster-iq/Chart.yaml
  ```
* [ ] Check if Helm is able to process all the templates correctly
  ```sh
  helm template ./deployments/helm/cluster-iq -f ./deployments/helm/cluster-iq/values.yaml
  ```
* [ ] Review the PR on GitHub to check if the GHActions were executed correctly.
* [ ] Merge the PR (:warning: **DO NOT SQUASH** we want a merge commit to
    automatically generate the changelog)
* [ ] If the merge results in a new release, create the tag and push it
  ```sh
  helm template ./deployments/helm/cluster-iq -f ./deployments/helm/cluster-iq/values.yaml
  git checkout main
  git pull
  git tag v<VERSION> -m "Release <VERSION>"
  git push --tags
  ```
