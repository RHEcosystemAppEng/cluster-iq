# ClusterIQ Data Migration Guide (v0.4.X → v0.5)
This document defines a deterministic, manual procedure to migrate ClusterIQ
data from version `0.5` to `0.5.1`.

The migration is performed directly in the target database:
- Apply DB schema transformations

## Scope
This procedure covers the data migration from ClusterIQ `0.5` to `0.5.1`.

This procedure does NOT:
- Deploy ClusterIQ
- Perform application upgrades
- Modify ClusterIQ configuration outside the database
- Execute automated CI/CD tasks
- Origin Database is running in Openshift
- Destination Database is running in the Development setup (podman-compose)

## Preconditions
- Source ClusterIQ version is `0.5`
- pgsql deployment using `v0.5.1` image tag
- Operator has access to OpenShift and PostgreSQL
- `oc`, `psql`, and `podman` CLIs are installed and configured
- No write operations are running against the source database

## Procedure Checklist

### Pre-checks

* [ ] **M0** — Prepare CLI and Namespace
  ```sh
  # Login into Openshift
  oc login …

  # Define Namespace
  export CIQ_NAMESPACE="<NAMESPACE>" 
  oc project $CIQ_NAMESPACE

  # Stop Scanner to prevent updates during the data migration
  oc patch cronjob scanner -p '{"spec" : {"suspend" : true }}' --type=merge -n $CIQ_NAMESPACE
  ```

* [ ] **M1** — Connect to the target Database (assuming it runs in OCP).
  ```sh
  oc rsh pgsql-0 -n $CIQ_NAMESPACE

  psql -d clusteriq
  ```

* [ ] **M2** — Check `RESOURCE_TYPE` enum exists.
  ```sql
  clusteriq=# \dT+ RESOURCE_TYPE
  ```
  **Expected result:**
  ```sql
                                           List of data types
   Schema |     Name      | Internal name | Size | Elements | Owner | Access privileges | Description
  --------+---------------+---------------+------+----------+-------+-------------------+-------------
   public | resource_type | resource_type | 4    | Account +| user  |                   |
          |               |               |      | Cluster +|       |                   |
          |               |               |      | Instance |       |                   |
  (1 row)

  ```

* [ ] **M3** — Modify instances table.
  ```sql
  -- Count before modifying
  SELECT COUNT(*) FROM instances;

  -- Tx
  BEGIN;

  -- 1. Drop FK
  ALTER TABLE instances DROP CONSTRAINT instances_cluster_id_fkey;

  -- 2. Change column type
  ALTER TABLE instances
    ALTER COLUMN cluster_id TYPE BIGINT
    USING cluster_id::bigint;

  -- 3. Recreate FK
  ALTER TABLE instances
    ADD CONSTRAINT instances_cluster_id_fkey
    FOREIGN KEY (cluster_id) REFERENCES clusters(id)
    ON DELETE CASCADE;

  COMMIT;

  -- Count after modifying must be same result as before
  SELECT COUNT(*) FROM instances;
  ```

  **Expected result:** Alter table tx correct

* [ ] **M4** — Modify events table.
  ```sql
  -- Count before modifying
  SELECT COUNT(*) FROM instances;

  -- Tx
  BEGIN;

  -- 1) Drop legacy CHECK constraint (if present)
  ALTER TABLE public.events
    DROP CONSTRAINT IF EXISTS events_resource_type_check;

  -- 2) resource_id: INTEGER -> BIGINT
  ALTER TABLE public.events
    ALTER COLUMN resource_id TYPE BIGINT
    USING resource_id::bigint;

  -- 3) resource_type: TEXT -> RESOURCE_TYPE (align with enum: Account/Cluster/Instance)
  ALTER TABLE public.events
    ALTER COLUMN resource_type TYPE public.resource_type
    USING (
      CASE resource_type
        WHEN 'cluster'  THEN 'Cluster'::public.resource_type
        WHEN 'Cluster'  THEN 'Cluster'::public.resource_type
        WHEN 'instance' THEN 'Instance'::public.resource_type
        WHEN 'Instance' THEN 'Instance'::public.resource_type
        WHEN 'account'  THEN 'Account'::public.resource_type
        WHEN 'Account'  THEN 'Account'::public.resource_type
        ELSE resource_type::public.resource_type
      END
    );

  -- 4) Keep NOT NULL (it was NOT NULL already, but explicit is fine)
  ALTER TABLE public.events
    ALTER COLUMN resource_type SET NOT NULL;

  COMMIT;

  -- Count after modifying must be same result as before
  SELECT COUNT(*) FROM instances;
  ```

  **Expected result:** Alter table tx correct

## Postconditions
- Tables modified correctly.
- Views re-created.
