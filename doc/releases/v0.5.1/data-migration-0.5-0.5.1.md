# ClusterIQ Data Migration Guide (v0.5 → v0.5.1)

This document defines the procedure to migrate a ClusterIQ database from version `0.5` to `0.5.1`.

## Changes Summary

| Object | Change | Risk |
|--------|--------|------|
| `instances.cluster_id` | `INTEGER` → `BIGINT` | Low — type widening, no data loss |
| `events.resource_id` | `INTEGER` → `BIGINT` | Low — type widening, no data loss |
| `events.resource_type` | `TEXT` → `RESOURCE_TYPE` enum | Medium — requires value conversion |
| `events_resource_type_check` | Removed (replaced by enum) | None |
| Cascade delete triggers | New (clusters, instances) | None — additive |
| `cluster_events` view | Updated join conditions | None — recreated in place |
| `system_events` view | Updated join conditions | None — recreated in place |
| `clusters_tags` view | Added WHERE filter | None — recreated in place |

## Preconditions

- Source ClusterIQ version is `v0.5`
- Target images are tagged `v0.5.1`
- The `RESOURCE_TYPE` enum already exists (created in v0.5)
- Operator has `oc` and `psql` access
- **No write operations running against the database during migration**

## Procedure

### M0 — Prepare environment

```sh
oc login ...

export CIQ_NAMESPACE="cluster-iq"
oc project $CIQ_NAMESPACE

# Stop scanner to prevent writes during migration
oc patch cronjob scanner -p '{"spec":{"suspend":true}}' --type=merge -n $CIQ_NAMESPACE
```

### M1 — Connect to database

```sh
oc rsh pgsql-0 -n $CIQ_NAMESPACE
psql -d clusteriq
```

### M2 — Pre-migration snapshot

Save row counts before any change. These must match after each step.

```sql
SELECT 'instances' AS tbl, COUNT(*) AS rows FROM instances
UNION ALL
SELECT 'events',          COUNT(*)           FROM events;
```

Record the output. Every verification step below must return these same numbers.

### M3 — Verify RESOURCE_TYPE enum exists

```sql
SELECT enumlabel FROM pg_enum
WHERE enumtypid = 'public.resource_type'::regtype
ORDER BY enumsortorder;
```

Expected output:

```
 enumlabel
-----------
 Account
 Cluster
 Instance
```

If the enum does not exist, **stop here** — the database is not on v0.5.

### M4 — Verify event data is safe to convert

Check that all existing `resource_type` values can map to the enum:

```sql
SELECT DISTINCT resource_type FROM events
WHERE resource_type NOT IN ('cluster', 'Cluster', 'instance', 'Instance', 'account', 'Account');
```

Expected output: **0 rows**. If any rows appear, those values need manual mapping before continuing.

### M5 — Alter `instances` table

```sql
BEGIN;

ALTER TABLE instances DROP CONSTRAINT instances_cluster_id_fkey;

ALTER TABLE instances
  ALTER COLUMN cluster_id TYPE BIGINT
  USING cluster_id::bigint;

ALTER TABLE instances
  ADD CONSTRAINT instances_cluster_id_fkey
  FOREIGN KEY (cluster_id) REFERENCES clusters(id)
  ON DELETE CASCADE;

COMMIT;
```

Verify:

```sql
SELECT COUNT(*) FROM instances;
-- Must match M2
```

### M6 — Alter `events` table

```sql
BEGIN;

ALTER TABLE events
  DROP CONSTRAINT IF EXISTS events_resource_type_check;

ALTER TABLE events
  ALTER COLUMN resource_id TYPE BIGINT
  USING resource_id::bigint;

ALTER TABLE events
  ALTER COLUMN resource_type TYPE public.resource_type
  USING (
    CASE resource_type
      WHEN 'cluster'  THEN 'Cluster'::public.resource_type
      WHEN 'Cluster'  THEN 'Cluster'::public.resource_type
      WHEN 'instance' THEN 'Instance'::public.resource_type
      WHEN 'Instance' THEN 'Instance'::public.resource_type
      WHEN 'account'  THEN 'Account'::public.resource_type
      WHEN 'Account'  THEN 'Account'::public.resource_type
    END
  );

ALTER TABLE events
  ALTER COLUMN resource_type SET NOT NULL;

COMMIT;
```

Verify:

```sql
SELECT COUNT(*) FROM events;
-- Must match M2
```

### M7 — Create cascade delete triggers

```sql
CREATE OR REPLACE FUNCTION delete_cluster_events()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM events WHERE resource_type = 'Cluster'::RESOURCE_TYPE AND resource_id = OLD.id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_delete_cluster_events
    BEFORE DELETE ON clusters
    FOR EACH ROW
    EXECUTE FUNCTION delete_cluster_events();

CREATE OR REPLACE FUNCTION delete_instance_events()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM events WHERE resource_type = 'Instance'::RESOURCE_TYPE AND resource_id = OLD.id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_delete_instance_events
    BEFORE DELETE ON instances
    FOR EACH ROW
    EXECUTE FUNCTION delete_instance_events();
```

Verify:

```sql
SELECT tgname FROM pg_trigger
WHERE tgname IN ('trg_delete_cluster_events', 'trg_delete_instance_events');
```

Expected: 2 rows.

### M8 — Recreate views

The event views must be updated to use the enum type instead of plain text comparisons.

```sql
CREATE OR REPLACE VIEW cluster_events AS
SELECT
  ev.id,
  ev.event_timestamp,
  ev.triggered_by,
  ev.action,
  COALESCE(c.cluster_id, i.instance_id) AS resource_id,
  ev.resource_type,
  ev.result,
  ev.description,
  ev.severity
FROM events ev
LEFT JOIN clusters  c ON ev.resource_type = 'Cluster'::RESOURCE_TYPE  AND c.id = ev.resource_id
LEFT JOIN instances i ON ev.resource_type = 'Instance'::RESOURCE_TYPE AND i.id = ev.resource_id
ORDER BY event_timestamp DESC;

CREATE OR REPLACE VIEW system_events AS
SELECT
  ev.id,
  ev.event_timestamp,
  ev.triggered_by,
  ev.action,
  COALESCE(c.cluster_id, i.instance_id) AS resource_id,
  ev.resource_type,
  ev.result,
  ev.description,
  ev.severity,
  acc.account_id,
  acc.provider
FROM events ev
LEFT JOIN clusters  c ON ev.resource_type = 'Cluster'::RESOURCE_TYPE  AND c.id = ev.resource_id
LEFT JOIN instances i ON ev.resource_type = 'Instance'::RESOURCE_TYPE AND i.id = ev.resource_id
LEFT JOIN accounts acc ON acc.id = (
  CASE
    WHEN ev.resource_type = 'Cluster'::RESOURCE_TYPE
    THEN (SELECT c.account_id FROM clusters c WHERE c.id = ev.resource_id)
    WHEN ev.resource_type = 'Instance'::RESOURCE_TYPE
    THEN (SELECT c.account_id FROM clusters c WHERE c.id = (SELECT i.cluster_id FROM instances i WHERE i.id = ev.resource_id))
  END
)
ORDER BY ev.event_timestamp DESC;

CREATE OR REPLACE VIEW clusters_tags AS
SELECT
    c.cluster_id,
    t.key,
    MIN(t.value) AS value
FROM clusters   c
JOIN instances  i ON i.cluster_id = c.id
JOIN tags       t ON t.instance_id = i.id
WHERE t.key != 'Name' AND t.key != 'MachineName'
GROUP BY c.cluster_id, t.key
HAVING COUNT(*) > 1;
```

Verify:

```sql
SELECT viewname FROM pg_views
WHERE schemaname = 'public'
  AND viewname IN ('cluster_events', 'system_events', 'clusters_tags');
```

Expected: 3 rows.

### M9 — Refresh materialized views

```sql
SELECT refresh_materialized_views();
```

### M10 — Final verification

```sql
SELECT 'instances' AS tbl, COUNT(*) AS rows FROM instances
UNION ALL
SELECT 'events',          COUNT(*)           FROM events;
```

Must match the counts from M2 exactly.

### M11 — Resume scanner

Exit psql and the pod, then:

```sh
oc patch cronjob scanner -p '{"spec":{"suspend":false}}' --type=merge -n $CIQ_NAMESPACE
```

## Rollback

If any step fails, the transaction (BEGIN/COMMIT) ensures the table is unchanged.
Roll back by fixing the issue and re-running the failed step. No partial state is possible for M5 and M6 since they are wrapped in transactions.

For M7 and M8 (triggers and views), these are idempotent — re-running them is safe.
