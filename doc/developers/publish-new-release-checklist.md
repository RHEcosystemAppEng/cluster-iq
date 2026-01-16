# New ClusterIQ Release publishing Checklist

This document defines the standard, deterministic procedure to review and
publish a new ClusterIQ release in the [Git repository](https://github.com/RHEcosystemAppEng/cluster-iq).

## Scope
This procedure covers the review and publication of a new ClusterIQ version.

This procedure does **NOT**:
- Deploy ClusterIQ
- Modify any running ClusterIQ instance
- Perform operational actions on any environment

## Preconditions
- Local Git repository is clean (`git status` shows no pending changes)
- User has push permissions to the ClusterIQ repository
- Go toolchain installed and configured
- Helm CLI installed
- GitHub Actions service available

## Procedure Checklist

* [ ] **P1** — Checkout the `release-X.Y.Z` branch to be reviewed and ensure it is up to date.
  ```sh
  export CIQ_RELEASE="release-X.Y.Z" # Set your version here
  ```
  ```sh
  git checkout $CIQ_RELEASE
  git pull <REMOTE> $CIQ_RELEASE
  git branch --show-current
  ```
  **Expected result:** current branch is `release-X.Y.Z`

* [ ] **P2** — Create a Pull Request from `release-X.Y.Z` to `main` named `Release <VERSION>`.

  **Expected result:** PR is open and base branch is `main`

* [ ] **P3** — Verify that all changes intended for this release are merged into `release-X.Y.Z`.

  **Expected result:** no pending PRs or commits expected for this release

* [ ] **P4** — Run unit tests and linters.
  ```sh
  make clean build go-tests
  ```
  **Expected result:** command exits with status `0`
  **DO NOT CONTINUE** if any error is reported

* [ ] **P5** — Update the `./VERSION` file with the release version.

  **Expected result:** `./VERSION` contains exactly `vX.Y.Z`

* [ ] **P6** — Compare `db/sql/init.sql` with the Helm database initialization ConfigMap.
  ```sh
  vim db/sql/init.sql deployments/helm/cluster-iq/templates/database/configmap-init.yaml -d
  ```
  **Expected result:** differences are identified or confirmed absent

* [ ] **P7** — Update the Helm database initialization ConfigMap if differences are found.

  **Expected result:** ConfigMap reflects the content of `db/sql/init.sql`

* [ ] **P8** — Compare `db/sql/cron.sql` with the Helm database initialization ConfigMap.
  ```sh
  vim db/sql/cron.sql deployments/helm/cluster-iq/templates/database/configmap-init.yaml -d
  ```
  **Expected result:** differences are identified or confirmed absent

* [ ] **P9** — Update the Helm database initialization ConfigMap if differences are found.

  **Expected result:** ConfigMap reflects the content of `db/sql/cron.sql`

* [ ] **P10** — Update Helm Chart `version` and `appVersion` if required.
  ```sh
  vim deployments/helm/cluster-iq/Chart.yaml
  ```
  **Expected result:** Chart metadata reflects the new release version

* [ ] **P11** — Validate Helm templates rendering.
  ```sh
  helm template deployments/helm/cluster-iq -f deployments/helm/cluster-iq/values.yaml
  echo $?
  ```
  **Expected result:** command completes without errors. Return code `0`.

* [ ] **P12** — Include these changes in a single commit describing the new
    release changes
  ```sh
  # Check the changes performed during the release preparation procedure
  git status

  # Add the changes
  git add <CHANGES>

  git commit -s -S -m "release: Changes for preparing release: $CIQ_RELEASE"
  git push <REMOTE> $CIQ_RELEASE # Remember to specify the remote name
  ```
  **Expected result:** No changes pending.

* [ ] **P13** — Verify all GitHub Actions workflows in the PR completed successfully.

  **Expected result:** all checks are green

* [ ] **P14** — Merge the Pull Request into `main`.
  **DO NOT SQUASH** (a merge commit is required for changelog generation)

  **Expected result:** merge commit created on `main`

* [ ] **P15** — Checkout `main` and ensure it is up to date.
  ```sh
  git checkout main
  git pull
  ```
  **Expected result:** local `main` matches remote `main`

* [ ] **P16** — Create and push the release tag.
  ```sh
  git tag v<VERSION> -m "Release <VERSION>"
  git push --tags
  ```
  **Expected result:** tag `v<VERSION>` exists on `main` and is pushed to the remote

* [ ] **P17** — Create the Github Release associating it to the new tag, and
    generate the changelog.
  **Expected result:** A new release available in the Github repository.


## Postconditions
- Release PR merged into `main`
- Git tag `v<VERSION>` created and pushed
- All GitHub Actions workflows completed successfully
- New Release available in Github with its changelog and tag

## Rollback
This procedure does not include rollback steps.

If an issue is detected after tagging, a corrective release must be prepared
following the same procedure.
