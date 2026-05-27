# ClusterIQ Data Migration Guide (v0.4.X → v0.5)

This document defines a deterministic, manual procedure to migrate ClusterIQ
data from version `0.4.X` to `0.5`.

The migration is performed as an ETL process:
- Extract data from the source database
- Transform it to match the new data model
- Recreate and populate the destination database

## Scope
This procedure covers the data migration from ClusterIQ `0.4.X` to `0.5`.

This procedure does NOT:
- Deploy ClusterIQ
- Perform application upgrades
- Modify ClusterIQ configuration outside the database
- Execute automated CI/CD tasks
- Origin Database is running in Openshift
- Destination Database is running in the Development setup (podman-compose)

## Preconditions
- Source ClusterIQ version is `0.4.X`
- Target database schema for `0.5` is already created
- Operator has access to OpenShift and PostgreSQL
- `oc`, `psql`, and `podman` CLIs are installed and configured
- Sufficient disk space available for CSV backups
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

* [ ] **M1** — Prepare local backup directory and access the source database.
  ```sh
  # Count of elements (save the output for compare them later)
  oc project $CIQ_NAMESPACE
  oc rsh pgsql-0

  # Export data
  psql -d clusteriq -c "SELECT count(*) as account_count FROM accounts;"
  psql -d clusteriq -c "SELECT count(*) as cluster_count FROM clusters;"
  psql -d clusteriq -c "SELECT count(*) as instance_count FROM instances;"
  psql -d clusteriq -c "SELECT count(*) as tag_count FROM tags;"
  psql -d clusteriq -c "SELECT count(*) as expense_count FROM expenses;"
  psql -d clusteriq -c "SELECT count(*) as schedule_count FROM schedule;"
  psql -d clusteriq -c "SELECT count(*) as logs_count FROM audit_logs;"
  ```
  **Expected result:** row counts are recorded for later comparison

---

### Extract

* [ ] **M2** — Connect to the source PostgreSQL instance and prepare backup directory.
  ```sh
  # Connect to source DB
  oc project $CIQ_NAMESPACE

  # Run remote shell
  oc rsh pgsql-0

  # Prepare backup destination folder
  mkdir /tmp/backups

  # Connect to DB (PRO)
  psql -d clusteriq
  ```
  **Expected result:** PostgreSQL prompt available on source DB

* [ ] **M3** — Export all required tables to CSV files.
  ```sql
  -- Export Accounts
  \COPY (SELECT id as account_id, name as account_name, provider, last_scan_timestamp as last_scan_ts, now() as created_at FROM accounts) TO '/tmp/backups/accounts.csv' WITH(FORMAT csv, HEADER);

  -- Export Clusters
  \COPY (SELECT id as cluster_id, name as cluster_name, infra_id, provider, status, region, account_name as account_id, console_link, last_scan_timestamp as last_scan_ts, creation_timestamp as created_at, age, owner FROM clusters) TO '/tmp/backups/clusters.csv' WITH (FORMAT csv, HEADER);

  -- Export Instances
  \COPY (SELECT id as instance_id, name as instance_name, instance_type, provider, availability_zone, status, cluster_id, last_scan_timestamp as last_scan_ts, creation_timestamp as created_at, age FROM instances) TO '/tmp/backups/instances.csv' WITH (FORMAT csv, HEADER);

  -- Export Tags
  \COPY (SELECT * FROM tags) TO '/tmp/backups/tags.csv' WITH (FORMAT csv, HEADER);

  -- Export Expenses
  \COPY (SELECT * FROM expenses) TO '/tmp/backups/expenses.csv' WITH (FORMAT csv, HEADER);

  -- Export audit_logs
  \COPY (SELECT * FROM audit_logs) TO '/tmp/backups/auditlogs.csv' WITH (FORMAT csv, HEADER);

  -- Export schedule
  \COPY (SELECT * FROM schedule) TO '/tmp/backups/schedule.csv' WITH (FORMAT csv, HEADER);

  # Exit
  \q
  ```
  ```sh
  # Check exported files
  ls -la /tmp/backups

  exit
  ```
  ```sh
  # Resume the scanner
  oc patch cronjob scanner -p '{"spec" : {"suspend" : false }}' --type=merge -n $CIQ_NAMESPACE
  ```
  **Expected result:** CSV files created under `/tmp/backups`

---

### Transfer Data to Destination

* [ ] **M4** — Copy exported CSV files to local system and into the destination database container.
  ```sh
  # Now, copy the files '/tmp/backups/*.csv' into your local
  mkdir -p ./migration/backups
  oc cp pgsql-0:/tmp/backups/ ./migration/backups/

  # Copy data to the new DB
  podman exec pgsql mkdir -p /tmp/backups
  podman cp ./migration/backups pgsql:/tmp
  podman exec -it pgsql psql -d clusteriq
  ```
  **Expected result:** CSV files available under `/tmp/backups` in destination DB

---

### Restore Data

* [ ] **M5** — Restore Accounts data.
  ```sql
  DELETE FROM accounts;
  ALTER SEQUENCE accounts_id_seq RESTART WITH 1;
  \COPY accounts (account_id, account_name, provider, last_scan_ts, created_at) FROM '/tmp/backups/accounts.csv' CSV HEADER;
  ```
  **Expected result:** accounts table populated from CSV

* [ ] **M6** — Restore Clusters data.
  ```sql
  -- Temporal table for loading backup
  CREATE TEMP TABLE stage_clusters (
    id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    cluster_id              TEXT NOT NULL,
    cluster_name            TEXT NOT NULL,
    infra_id                TEXT NOT NULL,
    provider                CLOUD_PROVIDER NOT NULL,
    status                  STATUS NOT NULL,
    region                  TEXT,
    account_id              TEXT,
    console_link            TEXT,
    last_scan_ts            TIMESTAMP WITH TIME ZONE,
    created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    age                     INTEGER,
    owner                   TEXT
  ); 

  -- Loading clusters backup
  \COPY stage_clusters (cluster_id, cluster_name, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner) FROM '/tmp/backups/clusters.csv' CSV HEADER;

  -- Loading data to final clusters table respecting new ID
  DELETE FROM clusters;
  ALTER SEQUENCE clusters_id_seq RESTART WITH 1;
  WITH src AS (
    SELECT
      sc.cluster_id,
      sc.cluster_name,
      sc.infra_id,
      sc.provider,
      sc.status,
      sc.region,
      ac.id AS account_id,
      sc.console_link,
      sc.last_scan_ts,
      sc.created_at,
      sc.age,
      sc.owner
    FROM stage_clusters AS sc 
    JOIN accounts AS ac
    ON sc.account_id = ac.account_name
  )
  INSERT INTO clusters (
    cluster_id,
    cluster_name,
    infra_id,
    provider,
    status,
    region,
    account_id,
    console_link,
    last_scan_ts,
    created_at,
    age,
    owner
  )
  SELECT * FROM src;

  DROP TABLE stage_clusters;

  -- Remove "UNKNOWN" clusters
  DELETE FROM clusters WHERE cluster_name = 'NO_CLUSTER';
  DELETE FROM clusters WHERE cluster_name = 'UNKNOWN-CLUSTER' AND infra_id = '';
  ```
  **Expected result:** clusters table populated and cleaned

* [ ] **M7** — Restore Instances data.
  ```sql
  -- Temporal table for loading backup
  CREATE TEMP TABLE stage_instances (
    id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    instance_id             TEXT NOT NULL,
    instance_name           TEXT,
    instance_type           TEXT,
    provider                CLOUD_PROVIDER NOT NULL,
    availability_zone       TEXT,
    status                  STATUS NOT NULL,
    cluster_id              TEXT,
    last_scan_ts            TIMESTAMP WITH TIME ZONE,
    created_at              TIMESTAMP WITH TIME ZONE,
    age                     INTEGER
  ); 

  -- Loading instances backup
  \COPY stage_instances (instance_id, instance_name, instance_type, provider, availability_zone, status, cluster_id, last_scan_ts, created_at, age) FROM '/tmp/backups/instances.csv' CSV HEADER;

  -- Loading data to final clusters table respecting new ID
  DELETE FROM instances;
  ALTER SEQUENCE instances_id_seq RESTART WITH 1;
  WITH src AS (
    SELECT
      si.instance_id,
      si.instance_name,
      si.instance_type,
      si.provider,
      si.availability_zone,
      si.status,
      cl.id AS cluster_id,
      si.last_scan_ts,
      si.created_at,
      si.age
    FROM stage_instances AS si
    JOIN clusters AS cl
    ON si.cluster_id = cl.cluster_id
  )
  INSERT INTO instances (
    instance_id,
    instance_name,
    instance_type,
    provider,
    availability_zone,
    status,
    cluster_id,
    last_scan_ts,
    created_at,
    age
  )
  SELECT * FROM src;

  -- WARNING! If you see less instance count being inserted into 'instances' table
  -- is because we removed the duplicated NO_CLUSTERS! You can check it with this query
  -- QUERY:
  --   SELECT * FROM stage_instances si WHERE NOT EXISTS (SELECT 1 FROM instances i WHERE i.instance_id = si.instance_id);
  -- If the COPY count == INSERT.count + QUERY.count, you're ok

  -- Run this query, if it returns '0', you can continue :)
  --   SELECT (SELECT count(*) FROM instances) + (SELECT count(*) FROM stage_instances si WHERE NOT EXISTS (SELECT 1 FROM instances i WHERE i.instance_id = si.instance_id)) - (SELECT count(*) FROM stage_instances);

  DROP TABLE stage_instances;
  ```
  **Expected result:** instances table populated with valid cluster references

* [ ] **M8** — Restore Tags data.
  ```sql
  -- Temporal table for loading backup
  CREATE TEMP TABLE stage_tags (
    key                     TEXT NOT NULL,
    value                   TEXT,
    instance_id             TEXT
  );

  -- Loading tags backup
  \COPY stage_tags (key, value, instance_id) FROM '/tmp/backups/tags.csv' CSV HEADER;

  -- Loading data to final table respecting new ID
  DELETE FROM tags;
  WITH src AS (
    SELECT
      st.key,
      st.value,
      i.id
    FROM stage_tags AS st
    JOIN instances AS i
    ON st.instance_id = i.instance_id
  )
  INSERT INTO tags (
    key,
    value,
    instance_id
  )
  SELECT * FROM src;

  DROP TABLE stage_tags;
  ```
  **Expected result:** tags table populated

* [ ] **M9** — Restore Expenses data.
  ```sql
  -- Temporal table for loading backup
  CREATE TEMP TABLE stage_expenses (
    instance_id             TEXT,
    date                    DATE,
    amount                  NUMERIC(12,2)
  );

  CREATE TABLE IF NOT EXISTS expenses_default PARTITION OF expenses DEFAULT;


  -- Loading tags backup
  \COPY stage_expenses (instance_id, date, amount) FROM '/tmp/backups/expenses.csv' CSV HEADER;

  -- Loading data to final table respecting new ID
  DELETE FROM expenses;
  WITH src AS (
    SELECT
      i.id,
      se.date,
      se.amount
    FROM stage_expenses AS se
    JOIN instances AS i
    ON se.instance_id = i.instance_id
  )
  INSERT INTO expenses (
    instance_id,
    date,
    amount
  )
  SELECT * FROM src;

  DROP TABLE stage_expenses;
  ```
  **Expected result:** expenses table populated

* [ ] **M10** — Restore Audit Logs data.
  ```sql
  -- Temporal table for loading backup
  CREATE TEMP TABLE IF NOT EXISTS stage_events (
    id BIGINT           GENERATED ALWAYS AS IDENTITY NOT NULL,
    event_timestamp     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    triggered_by        TEXT NOT NULL,
    action_name         TEXT NOT NULL,
    resource_id         TEXT NOT NULL,
    resource_type       TEXT NOT NULL,
    result              TEXT NOT NULL,
    description         TEXT NULL,
    severity            TEXT DEFAULT 'info'::TEXT NOT NULL
  );

  -- Loading events backup
  \COPY stage_events (id, event_timestamp, triggered_by, action_name, resource_id, resource_type, result, description, severity) FROM '/tmp/backups/auditlogs.csv' CSV HEADER;

  -- Loading data to final table
  DELETE FROM events;
  ALTER SEQUENCE events_id_seq RESTART WITH 1;
  WITH src AS (
    SELECT DISTINCT ON (se.id)
      se.id,
      se.event_timestamp,
      se.triggered_by,
      se.action_name,
      se.resource_type,
      se.result,
      se.description,
      se.severity,
      CASE
        WHEN se.resource_type = 'cluster'  THEN c.id
        WHEN se.resource_type = 'instance' THEN i.id
        ELSE NULL
      END AS resource_fk_id
    FROM stage_events se
    LEFT JOIN clusters  c ON se.resource_type='cluster'  AND c.cluster_id  = se.resource_id
    LEFT JOIN instances i ON se.resource_type='instance' AND i.instance_id = se.resource_id
    ORDER BY se.id, se.event_timestamp DESC
  )
  INSERT INTO events (
    event_timestamp, triggered_by, action, resource_id, resource_type, result, description, severity
  )
  SELECT
    event_timestamp,
    triggered_by,
    action_name,
    resource_fk_id,
    resource_type,
    result::ACTION_STATUS,
    description,
    severity
  FROM src;

  DELETE FROM events WHERE triggered_by='ClusterIQ Agent';

  DROP TABLE stage_events;
  ```
  **Expected result:** audit_logs table populated

* [ ] **M11** — Restore Schedule data.
  ```sql
  CREATE TEMP TABLE IF NOT EXISTS stage_schedule (
    id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    type                    ACTION_TYPE,
    time                    TIMESTAMP WITH TIME ZONE,
    cron_exp                TEXT,
    operation               TEXT NOT NULL,
    target                  TEXT,
    status                  ACTION_STATUS,
    enabled                 BOOLEAN
  );

  \COPY stage_schedule (id, type, time, cron_exp, operation, target, status, enabled) FROM '/tmp/backups/schedule.csv' CSV HEADER;

  UPDATE stage_schedule SET time = NULL WHERE cron_exp != '';
  UPDATE stage_schedule SET operation = 'PowerOn' WHERE operation = 'PowerOnCluster';
  UPDATE stage_schedule SET operation = 'PowerOff' WHERE operation = 'PowerOffCluster';

  -- Loading data to final table respecting new ID
  DELETE FROM schedule;
  ALTER SEQUENCE schedule_id_seq RESTART WITH 1;
  WITH src AS (
    SELECT
      ss.type,
      ss.time,
      ss.cron_exp,
      ss.operation::ACTION_OPERATION,
      c.id AS target,
      ss.status,
      ss.enabled
    FROM stage_schedule AS ss
    JOIN clusters AS c
    ON ss.target = c.cluster_id
  )
  INSERT INTO schedule (
    type,
    time,
    cron_exp,
    operation,
    target,
    status,
    enabled
  )
  SELECT * FROM src;

  DROP TABLE stage_schedule;
  ```
  **Expected result:** schedule table populated

---

### Final Data Cleanup

* [ ] **M12** — Clean and normalize clusters data.
  ```sql
  -- Verifying "-<ACCOUNT_NAME>" suffix
  SELECT
    c.id,
    c.cluster_id,
    a.account_name
  FROM clusters c
  JOIN accounts a ON a.id = c.account_id
  WHERE c.cluster_id ~* ('-' || regexp_replace(a.account_name, '([\\W])', '\\\1', 'g') || '$');

  -- Remove "-<ACCOUNT_NAME>" suffix from cluster_id
  UPDATE clusters c
  SET cluster_id = regexp_replace(
    c.cluster_id,
    '-' || regexp_replace(a.account_name, '([\\W])', '\\\1', 'g') || '$',
    '',
    'i'
  )
  FROM accounts a
  WHERE a.id = c.account_id
    AND c.cluster_id ~* ('-' || regexp_replace(a.account_name, '([\\W])', '\\\1', 'g') || '$');

  -- Processing no-clustered clusters
  UPDATE clusters SET infra_id = '', cluster_name = 'NO_CLUSTER', cluster_id = 'NO_CLUSTER' WHERE infra_id = '' OR infra_id = 'UNKNOWN-CLUSTER';
  -- Processing cluster_id column for removing embeeded account_id
  UPDATE clusters SET cluster_id = regexp_replace(cluster_id, '-' || account_id || '$', '') WHERE cluster_id LIKE '%' || account_id AND infra_id != '';
  -- Processing console_link
  UPDATE clusters SET console_link = '' WHERE console_link = 'UNKNOWN-CONSOLE' OR console_link = 'Unknown Console Link'; 
  -- Updating 'created_at' and 'age' column
  UPDATE clusters SET age = 0, created_at = now() WHERE created_at = '0001-01-01 00:00:00+00';
  ```
  **Expected result:** clusters data normalized

* [ ] **M13** — Clean and normalize instances data.
  ```sql
  -- Updating 'created_at' and 'age' column
  UPDATE instances SET age = 0, created_at = now() WHERE created_at = '0001-01-01 00:00:00+00';
  ```
  **Expected result:** instances data normalized

* [ ] **M14** — Refresh Materialized views
  ```sql
  -- Refresh Materialized views
  SELECT refresh_materialized_views();
  ```
  **Expected result:** No return and views refreshed

## Postconditions
- All tables restored and populated
- Row counts match or are explainable against pre-check values
- No temporary tables remain
- Database is ready for ClusterIQ `0.5`
