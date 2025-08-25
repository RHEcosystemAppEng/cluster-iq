-- ## Basic Inventory definition (Types, Tables and partitions) ##
-- #################################################################################################
-- Supported values for Cloud Providers
CREATE TYPE CLOUD_PROVIDER AS ENUM (
  'AWS',
  'GCP',
  'Azure',
  'UNKNOWN'
);


-- Supported values of Status (Valid for Clusters and and Instances)
CREATE TYPE STATUS AS ENUM (
  'Running',
  'Stopped',
  'Terminated',
  'Unknown'
);


-- Supported values of resources_types
CREATE TYPE RESOURCE_TYPE AS ENUM (
  'Account',
  'Cluster',
  'Instance'
);


-- Accounts Table
CREATE TABLE IF NOT EXISTS accounts (
  id                      INT GENERATED ALWAYS AS IDENTITY NOT NULL,
  account_id              TEXT NOT NULL,
  account_name            TEXT NOT NULL,
  provider                CLOUD_PROVIDER NOT NULL,
  last_scan_ts            TIMESTAMP WITH TIME ZONE,
  created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  UNIQUE (account_id)
);


-- Clusters Table
CREATE TABLE IF NOT EXISTS clusters (
  id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  cluster_id              TEXT NOT NULL,                                   -- cluster_id is the result of joining: "name+infra_id"
  cluster_name            TEXT NOT NULL,
  infra_id                TEXT NOT NULL,
  provider                CLOUD_PROVIDER NOT NULL,
  status                  STATUS,
  region                  TEXT,
  account_id              INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
  console_link            TEXT,
  last_scan_ts            TIMESTAMP WITH TIME ZONE,
  created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  age                     INTEGER DEFAULT 0,
  owner                   TEXT,
  PRIMARY KEY (id),
  UNIQUE (id, cluster_id, account_id)
) PARTITION BY HASH (id);

CREATE TABLE clusters_p0 PARTITION OF clusters FOR VALUES WITH (MODULUS 8, REMAINDER 0);
CREATE TABLE clusters_p1 PARTITION OF clusters FOR VALUES WITH (MODULUS 8, REMAINDER 1);
CREATE TABLE clusters_p2 PARTITION OF clusters FOR VALUES WITH (MODULUS 8, REMAINDER 2);
CREATE TABLE clusters_p3 PARTITION OF clusters FOR VALUES WITH (MODULUS 8, REMAINDER 3);
CREATE TABLE clusters_p4 PARTITION OF clusters FOR VALUES WITH (MODULUS 8, REMAINDER 4);
CREATE TABLE clusters_p5 PARTITION OF clusters FOR VALUES WITH (MODULUS 8, REMAINDER 5);
CREATE TABLE clusters_p6 PARTITION OF clusters FOR VALUES WITH (MODULUS 8, REMAINDER 6);
CREATE TABLE clusters_p7 PARTITION OF clusters FOR VALUES WITH (MODULUS 8, REMAINDER 7);


-- Instances
CREATE TABLE IF NOT EXISTS instances (
  id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  instance_id             TEXT NOT NULL,
  instance_name           TEXT,
  instance_type           TEXT,
  provider                CLOUD_PROVIDER NOT NULL,
  availability_zone       TEXT,
  status                  STATUS,
  cluster_id              INTEGER REFERENCES clusters(id) ON DELETE CASCADE,
  last_scan_ts            TIMESTAMP WITH TIME ZONE,
  created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  age                     INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE (id, instance_name, cluster_id)
) PARTITION BY HASH(id);

CREATE TABLE instances_p00 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 00);
CREATE TABLE instances_p01 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 01);
CREATE TABLE instances_p02 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 02);
CREATE TABLE instances_p03 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 03);
CREATE TABLE instances_p04 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 04);
CREATE TABLE instances_p05 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 05);
CREATE TABLE instances_p06 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 06);
CREATE TABLE instances_p07 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 07);
CREATE TABLE instances_p08 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 08);
CREATE TABLE instances_p09 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 09);
CREATE TABLE instances_p10 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 10);
CREATE TABLE instances_p11 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 11);
CREATE TABLE instances_p12 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 12);
CREATE TABLE instances_p13 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 13);
CREATE TABLE instances_p14 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 14);
CREATE TABLE instances_p15 PARTITION OF instances FOR VALUES WITH (MODULUS 16, REMAINDER 15);


-- Instances Tags
-- TODO Check if is more efficient to move this table to JSONB column on the Instances table
CREATE TABLE IF NOT EXISTS tags (
  key                     TEXT,
  value                   TEXT,
  instance_id             BIGINT REFERENCES instances(id) ON DELETE CASCADE,
  PRIMARY KEY (key, instance_id)
) PARTITION BY HASH(instance_id);

CREATE TABLE tags_p00 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 00);
CREATE TABLE tags_p01 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 01);
CREATE TABLE tags_p02 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 02);
CREATE TABLE tags_p03 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 03);
CREATE TABLE tags_p04 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 04);
CREATE TABLE tags_p05 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 05);
CREATE TABLE tags_p06 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 06);
CREATE TABLE tags_p07 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 07);
CREATE TABLE tags_p08 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 08);
CREATE TABLE tags_p09 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 09);
CREATE TABLE tags_p10 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 10);
CREATE TABLE tags_p11 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 11);
CREATE TABLE tags_p12 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 12);
CREATE TABLE tags_p13 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 13);
CREATE TABLE tags_p14 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 14);
CREATE TABLE tags_p15 PARTITION OF tags FOR VALUES WITH (MODULUS 16, REMAINDER 15);


-- Instances expenses
-- TODO pg_cron to create partitions automatically
CREATE TABLE IF NOT EXISTS expenses (
  instance_id              BIGINT REFERENCES instances(id) ON DELETE CASCADE,
  date                     DATE,
  amount                   NUMERIC(12,2) DEFAULT 0.0,
  PRIMARY KEY (instance_id, date)
) PARTITION BY RANGE (date);

-- Default expenses partition. The rest of expenses will be created by pg_cron
CREATE TABLE expenses_default PARTITION OF expenses DEFAULT;





-- ## Advanced Inventory definition (Views, Indexes...) ##
-- #################################################################################################

-- ## Accounts
-- #############################################################################
-- Accounts Cluster Count view
CREATE VIEW accounts_with_cluster_count AS
SELECT c.account_id, COUNT(*)::bigint AS cluster_count
FROM clusters c
GROUP BY c.account_id;


-- Accounts Costs view
CREATE VIEW accounts_with_costs AS
WITH base AS (
  SELECT a.id, e.date, e.amount
  FROM accounts a
  JOIN clusters c ON c.account_id = a.id
  JOIN instances i ON i.cluster_id = c.id
  JOIN expenses e ON e.instance_id = i.id
)
SELECT
  a.id,
  COALESCE(SUM(b.amount), 0.0)                                           AS total_cost,
  COALESCE(SUM(b.amount) FILTER (
            WHERE b.date >= current_date - 14), 0)                       AS last_15_days_cost,
  COALESCE(SUM(b.amount) FILTER (
            WHERE b.date >= (date_trunc('month', current_date)::date - INTERVAL '1 month')
              AND b.date <  date_trunc('month', current_date)::date), 0) AS last_month_cost,
  COALESCE(SUM(b.amount) FILTER (
            WHERE b.date >= date_trunc('month', current_date)::date
              AND b.date <= current_date), 0)                            AS current_month_so_far_cost
FROM accounts a
LEFT JOIN base b ON b.id = a.id
GROUP BY a.id;


-- Accounts Full view
-- TODO: check if we can include MATERIALIZED views (using pg_cron)
CREATE VIEW accounts_full_view AS
SELECT
  a.account_id,
  a.account_name,
  a.provider,
  a.last_scan_ts,
  a.created_at,
  COALESCE(cc.cluster_count, 0)                                        AS cluster_count,
  COALESCE(ac.total_cost, 0)                                           AS total_cost,
  COALESCE(ac.last_15_days_cost, 0)                                    AS last_15_days_cost,
  COALESCE(ac.last_month_cost, 0)                                      AS last_month_cost,
  COALESCE(ac.current_month_so_far_cost, 0)                            AS current_month_so_far_cost
FROM accounts a
LEFT JOIN accounts_with_cluster_count cc ON cc.account_id = a.id
LEFT JOIN accounts_with_costs         ac ON ac.id = a.id;


-- ## Accounts Indexes
CREATE INDEX ix_accounts_provider
  ON accounts (provider);

CREATE INDEX ix_accounts_id
  ON accounts (account_id);

CREATE INDEX ix_accounts_last_scan_ts_desc
  ON accounts (last_scan_ts DESC);

CREATE INDEX ix_accounts_name
  ON accounts (lower(account_name));





-- ## Clusters
-- #############################################################################

-- Cluster Instances Count view
CREATE VIEW clusters_with_instance_count AS
SELECT i.cluster_id, COUNT(*)::bigint AS instance_count
FROM instances i
GROUP BY i.cluster_id;


-- clusters Costs view
CREATE VIEW clusters_with_costs AS
WITH base AS (
  SELECT c.id, e.date, e.amount
  FROM clusters c
  JOIN instances i ON i.cluster_id = c.id
  JOIN expenses e ON e.instance_id = i.id
)
SELECT
  c.id,
  COALESCE(SUM(b.amount), 0.0)                                          AS total_cost,
  COALESCE(SUM(b.amount) FILTER (WHERE b.date >= current_date - 14), 0) AS last_15_days_cost,
  COALESCE(SUM(b.amount) FILTER (
            WHERE b.date >= (date_trunc('month', current_date)::date - INTERVAL '1 month')
              AND b.date <  date_trunc('month', current_date)::date), 0) AS last_month_cost,
  COALESCE(SUM(b.amount) FILTER (
            WHERE b.date >= date_trunc('month', current_date)::date
              AND b.date <= current_date), 0)                            AS current_month_so_far_cost
FROM clusters c
LEFT JOIN base b ON b.id = c.id
GROUP BY c.id;


-- clusters Full view
CREATE VIEW clusters_full_view AS
SELECT
  c.cluster_id,
  c.cluster_name,
  c.infra_id,
  c.provider,
  c.status,
  c.region,
  a.account_id,
  a.account_name,
  c.last_scan_ts,
  c.created_at,
  c.age,
  c.owner,
  COALESCE(cc.instance_count, 0)                                       AS instance_count,
  COALESCE(ac.total_cost, 0)                                           AS total_cost,
  COALESCE(ac.last_15_days_cost, 0)                                    AS last_15_days_cost,
  COALESCE(ac.last_month_cost, 0)                                      AS last_month_cost,
  COALESCE(ac.current_month_so_far_cost, 0)                            AS current_month_so_far_cost
FROM clusters c
LEFT JOIN accounts a ON c.account_id = a.id
LEFT JOIN clusters_with_instance_count cc ON cc.cluster_id = c.id
LEFT JOIN clusters_with_costs         ac ON ac.id = c.id;



-- ## Clusters Indexes
CREATE INDEX ix_clusters_account
  ON clusters (account_id);

CREATE INDEX ix_clusters_acct_status
  ON clusters (account_id, status);

CREATE INDEX ix_clusters_acct_region_status
  ON clusters (account_id, region, status);

CREATE INDEX ix_clusters_last_scan_ts_desc
  ON clusters (last_scan_ts DESC);

CREATE INDEX ix_clusters_acct_name
  ON clusters (account_id, lower(cluster_name));

CREATE INDEX ix_clusters_last_scan_active
  ON clusters (last_scan_ts)
  WHERE status <> 'Terminated';





-- ## Instances
-- #############################################################################

-- Instances Costs view
CREATE VIEW instances_with_costs AS
WITH base AS (
  SELECT i.id, e.date, e.amount
  FROM instances i
  JOIN expenses e ON e.instance_id = i.id
)
SELECT
  i.id,
  COALESCE(SUM(b.amount), 0.0)                                          AS total_cost,
  COALESCE(SUM(b.amount) FILTER (WHERE b.date >= current_date - 14), 0) AS last_15_days_cost,
  COALESCE(SUM(b.amount) FILTER (
            WHERE b.date >= (date_trunc('month', current_date)::date - INTERVAL '1 month')
              AND b.date <  date_trunc('month', current_date)::date), 0) AS last_month_cost,
  COALESCE(SUM(b.amount) FILTER (
            WHERE b.date >= date_trunc('month', current_date)::date
              AND b.date <= current_date), 0)                            AS current_month_so_far_cost
FROM instances i
LEFT JOIN base b ON b.id = i.id
GROUP BY i.id;


-- Instances Full view
CREATE VIEW instances_full_view AS
SELECT
  i.instance_id,
  i.instance_name,
  i.instance_type,
  i.provider,
  i.availability_zone,
  i.status,
  c.cluster_id,
  c.cluster_name,
  i.last_scan_ts,
  i.created_at,
  i.age,
  COALESCE(ic.total_cost, 0)                 AS total_cost,
  COALESCE(ic.last_15_days_cost, 0)          AS last_15_days_cost,
  COALESCE(ic.last_month_cost, 0)            AS last_month_cost,
  COALESCE(ic.current_month_so_far_cost, 0)  AS current_month_so_far_cost
FROM instances i
LEFT JOIN clusters        c ON c.id = i.cluster_id
LEFT JOIN instances_with_costs ic ON ic.id = i.id;


-- Instances Full view
CREATE VIEW instances_full_view_with_tags AS
SELECT
  i.instance_id,
  i.instance_name,
  i.instance_type,
  i.provider,
  i.availability_zone,
  i.status,
  c.cluster_id,
  c.cluster_name,
  i.last_scan_ts,
  i.created_at,
  i.age,
  COALESCE(ic.total_cost, 0)                 AS total_cost,
  COALESCE(ic.last_15_days_cost, 0)          AS last_15_days_cost,
  COALESCE(ic.last_month_cost, 0)            AS last_month_cost,
  COALESCE(ic.current_month_so_far_cost, 0)  AS current_month_so_far_cost,
  COALESCE(t.tags, '[]'::jsonb)              AS tags_json
FROM instances i
LEFT JOIN clusters        c ON c.id = i.cluster_id
LEFT JOIN instances_with_costs ic ON ic.id = i.id
LEFT JOIN LATERAL (
    SELECT
      jsonb_agg(jsonb_build_object('key', t.key, 'value', t.value) ORDER BY t.key) AS tags
    FROM tags t
    WHERE t.instance_id = i.id
) t ON true;


-- Instances pending for expense update
CREATE VIEW instances_pending_expense_update AS
SELECT
	instances.instance_id
FROM
	instances
LEFT JOIN (
	SELECT
		instance_id,
		MAX(date) AS last_expense_date
	FROM
		expenses
	GROUP BY
		instance_id
) AS last_expenses
ON
	instances.id = last_expenses.instance_id
WHERE
	last_expenses.last_expense_date IS NULL
	OR last_expenses.last_expense_date < CURRENT_DATE - INTERVAL '1 day';



-- ## Instances Indexes
CREATE INDEX ix_instances_cluster
  ON instances (cluster_id);

CREATE INDEX ix_instances_cluster_status
  ON instances (cluster_id, status);

CREATE INDEX ix_instances_cluster_lastscan_desc
  ON instances (cluster_id, last_scan_ts DESC);

CREATE INDEX ix_instances_cluster_name
  ON instances (cluster_id, lower(instance_name));

CREATE INDEX ix_instances_last_scan_active
  ON instances (last_scan_ts)
  WHERE status <> 'Terminated';





-- ## Tags
-- #############################################################################

-- ## Tags Indexes
CREATE INDEX ix_itags_instance
  ON tags (instance_id);

CREATE INDEX ix_itags_key_val
  ON tags (key, value);

CREATE INDEX ix_itags_key
  ON tags (key);

CREATE INDEX ix_itags_key_lower_val
  ON tags (key, lower(value));





-- ## Expenses
-- #############################################################################

-- ## Expenses Indexes
CREATE INDEX ix_expenses_instance
  ON expenses (instance_id);





-- ## Actions and Scheduling definition ##
-- #############################################################################

-- Supported values of Action Operations
CREATE TYPE ACTION_OPERATION AS ENUM (
  'PowerOnCluster',
  'PowerOffCluster'
);


-- Supported values of action types
CREATE TYPE ACTION_TYPE AS ENUM (
  'instant_action',
  'cron_action',
  'scheduled_action'
);


-- Supported values of action status
CREATE TYPE ACTION_STATUS AS ENUM (
  'Pending',
  'Running',
  'Success',
  'Failed',
  'Unknown'
);


-- Scheduled actions
CREATE TABLE IF NOT EXISTS schedule (
  id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  type ACTION_TYPE NOT NULL DEFAULT 'scheduled_action',
  time TIMESTAMP WITH TIME ZONE,
  cron_exp TEXT,
  operation ACTION_OPERATION NOT NULL,
  target INTEGER REFERENCES clusters(id) ON DELETE CASCADE,
  status ACTION_STATUS NOT NULL DEFAULT 'Unknown',
  enabled BOOLEAN DEFAULT false
);







-- ## Audit logs
-- #############################################################################

-- Audit logs
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  triggered_by TEXT NOT NULL,
  action_name TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  resource_type TEXT NOT NULL,
  result TEXT NOT NULL,
  description TEXT NULL,
  severity TEXT DEFAULT 'info'::TEXT NOT NULL,
  CONSTRAINT audit_logs_resource_type_check CHECK ((resource_type = ANY (ARRAY['cluster'::TEXT, 'instance'::TEXT]))),
  PRIMARY KEY (id, event_timestamp)
) PARTITION BY RANGE (event_timestamp);







-- ## Maintenance Functions ##
-- #############################################################################

-- Marks instances as 'Terminated' if they haven't been scanned in the last 24 hours
CREATE OR REPLACE FUNCTION check_terminated_instances()
RETURNS void AS $$
BEGIN
  UPDATE instances
  SET status = 'Terminated'
  WHERE last_scan_timestamp < NOW() - INTERVAL '1 day';
END;
$$ LANGUAGE plpgsql;


-- Marks clusters as 'Terminated' if they haven't been scanned in the last 24 hours
CREATE OR REPLACE FUNCTION check_terminated_clusters()
RETURNS void AS $$
BEGIN
  UPDATE clusters
  SET status = 'Terminated'
  WHERE last_scan_timestamp < NOW() - INTERVAL '1 day';
END;
$$ LANGUAGE plpgsql;
