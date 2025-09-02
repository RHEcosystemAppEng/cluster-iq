# How to migrate ClusterIQ data from version 0.4.X to 0.5

The migration will consist in a ETL process where we are going to export the
needed data from the current DB instance, transform it, re-create the new DB
with the new data model, and restore

## 1. Pre-checks
```sh
# Access to the PGSQL source DB 
mkdir -p /tmp/backups

# Count of elements (save the output for compare them later)
psql -d clusteriq -c "SELECT count(*) as account_count FROM accounts;"
psql -d clusteriq -c "SELECT count(*) as cluster_count FROM clusters;"
psql -d clusteriq -c "SELECT count(*) as instance_count FROM instances;"
psql -d clusteriq -c "SELECT count(*) as tag_count FROM tags;"
psql -d clusteriq -c "SELECT count(*) as expense_count FROM expenses;"
psql -d clusteriq -c "SELECT count(*) as schedule_count FROM schedule;"
psql -d clusteriq -c "SELECT count(*) as logs_count FROM audit_logs;"
```

## 2. Extract
```sh
# Connect to "production"
oc project cluster-iq

# Run remote shell
oc rsh pgsql-0

# Prepare backup destination folder
mkdir /tmp/backups

# Connect to DB (PRO)
psql -d clusteriq
```

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
```

## 3. Save data to local
```sh
# Now, copy the files '/tmp/backups/*.csv' into your local
mkdir -p ./migration/backups
oc cp pgsql-0:/tmp/backups/ ./migration/backups/

# Copy data to the new DB
podman exec pgsql mkdir -p /tmp/backups
podman cp ./migration/backups pgsql:/tmp
podman exec -it pgsql psql -d clusteriq
```

## 4. Restore Accounts
```sql
DELETE FROM accounts;
ALTER SEQUENCE accounts_id_seq RESTART WITH 1;
\COPY accounts (account_id, account_name, provider, last_scan_ts, created_at) FROM '/tmp/backups/accounts.csv' CSV HEADER;
```

## 5. Restore Clusters
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

DELETE FROM clusters WHERE cluster_name = 'NO_CLUSTER';
DELETE FROM clusters WHERE cluster_name = 'UNKNOWN-CLUSTER' AND infra_id = '';
```

## 6. Restore Instances
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
     SELECT * FROM stage_instances si WHERE NOT EXISTS (SELECT 1 FROM instances i WHERE i.instance_id = si.instance_id);
-- If the COPY count == INSERT.count + QUERY.count, you're ok

DROP TABLE stage_instances;
```

## 7. Restore Tags
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
## 8. Restore Expenses
```sql
-- Temporal table for loading backup
CREATE TEMP TABLE stage_expenses (
  instance_id             TEXT,
  date                    DATE,
  amount                  NUMERIC(12,2)
);

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
## 9. Restore Audit Logs
```sql
\COPY audit_logs (id, event_timestamp, triggered_by, action_name, resource_id, resource_type, result, description, severity) FROM '/tmp/backups/auditlogs.csv' CSV HEADER;
```
## 10. Restore Schedule
```sql
CREATE TEMP TABLE IF NOT EXISTS stage_schedule (
  id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  type                    ACTION_TYPE,
  time                    TIMESTAMP WITH TIME ZONE,
  cron_exp                TEXT,
  operation               ACTION_OPERATION NOT NULL,
  target                  TEXT,
  status                  ACTION_STATUS,
  enabled                 BOOLEAN
);

\COPY stage_schedule (id, type, time, cron_exp, operation, target, status, enabled) FROM '/tmp/backups/schedule.csv' CSV HEADER;

UPDATE stage_schedule SET time = NULL WHERE cron_exp != '';

-- Loading data to final table respecting new ID
DELETE FROM schedule;
ALTER SEQUENCE schedule_id_seq RESTART WITH 1;
WITH src AS (
  SELECT
    ss.type,
    ss.time,
    ss.cron_exp,
    ss.operation,
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
## 11. Final Cleaning data
### Clusters
```sql
-- Processing no-clustered clusters
UPDATE clusters SET infra_id = '', cluster_name = 'NO_CLUSTER', cluster_id = 'NO_CLUSTER' WHERE infra_id = '' OR infra_id = 'UNKNOWN-CLUSTER';
-- Processing cluster_id column for removing embeeded account_id
UPDATE clusters SET cluster_id = regexp_replace(cluster_id, '-' || account_id || '$', '') WHERE cluster_id LIKE '%' || account_id AND infra_id != '';
-- Processing console_link
UPDATE clusters SET console_link = '' WHERE console_link = 'UNKNOWN-CONSOLE' OR console_link = 'Unknown Console Link'; 
-- Updating 'created_at' and 'age' column
UPDATE clusters SET age = 0, created_at = now() WHERE created_at = '0001-01-01 00:00:00+00';
```

### Instances
```sql
-- Updating 'created_at' and 'age' column
UPDATE instances SET age = 0, created_at = now() WHERE created_at = '0001-01-01 00:00:00+00';
```
