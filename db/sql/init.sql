-- ## Cluster IQ Database Definition ##
-- #################################################################################################

\! echo '## Initializing ClusterIQ Database definition'

-- ## Data Types
-- ############################################################

\! echo '## Creating custom Data Types'

-- Supported values of Status (Valid for Clusters and Instances)
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

-- Supported values for Cloud Providers
CREATE TYPE CLOUD_PROVIDER AS ENUM (
  'AWS',
  'GCP',
  'Azure',
  'UNKNOWN'
);



-- ## Tables, Indexes and Partitions
-- ############################################################

-- #############################################################################
-- Accounts
-- #############################################################################
\! echo '## Creating Accounts table'

CREATE TABLE IF NOT EXISTS accounts (
  id                      INT GENERATED ALWAYS AS IDENTITY NOT NULL,
  account_id              TEXT NOT NULL,
  account_name            TEXT NOT NULL,
  provider                CLOUD_PROVIDER NOT NULL,
  last_scan_ts            TIMESTAMP WITH TIME ZONE,
  created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  UNIQUE (provider, account_id),
  UNIQUE (provider, account_name)
);

CREATE INDEX IF NOT EXISTS ix_accounts_provider          ON accounts (provider);
CREATE INDEX IF NOT EXISTS ix_accounts_id                ON accounts (account_id);
CREATE INDEX IF NOT EXISTS ix_accounts_last_scan_ts_desc ON accounts (last_scan_ts DESC);
CREATE INDEX IF NOT EXISTS ix_accounts_name              ON accounts (lower(account_name));



-- #############################################################################
-- Clusters
-- #############################################################################
\! echo '## Creating Clusters table'

CREATE TABLE IF NOT EXISTS clusters (
  id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  cluster_id              TEXT NOT NULL,  -- cluster_id is the result of joining: "name+infra_id"
  cluster_name            TEXT NOT NULL,
  infra_id                TEXT NOT NULL,
  provider                CLOUD_PROVIDER NOT NULL,
  status                  STATUS DEFAULT 'Unknown' NOT NULL,
  region                  TEXT,
  account_id              INTEGER REFERENCES accounts(id) ON DELETE CASCADE NOT NULL,
  console_link            TEXT,
  last_scan_ts            TIMESTAMP WITH TIME ZONE,
  created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  age                     INTEGER DEFAULT 0,
  owner                   TEXT,
  PRIMARY KEY (id),
  CONSTRAINT uq_clusters_accountid_clusterid UNIQUE (account_id, cluster_id)
);

CREATE INDEX IF NOT EXISTS ix_clusters_account            ON clusters (account_id);
CREATE INDEX IF NOT EXISTS ix_clusters_acct_status        ON clusters (account_id, status);
CREATE INDEX IF NOT EXISTS ix_clusters_acct_region_status ON clusters (account_id, region, status);
CREATE INDEX IF NOT EXISTS ix_clusters_last_scan_ts_desc  ON clusters (last_scan_ts DESC);
CREATE INDEX IF NOT EXISTS ix_clusters_acct_name          ON clusters (account_id, lower(cluster_name));
CREATE INDEX IF NOT EXISTS ix_clusters_last_scan_active   ON clusters (last_scan_ts) WHERE status <> 'Terminated';



-- #############################################################################
-- Instances
-- #############################################################################
\! echo '## Creating Instances table'

CREATE TABLE IF NOT EXISTS instances (
  id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  instance_id             TEXT NOT NULL,
  instance_name           TEXT,
  instance_type           TEXT,
  provider                CLOUD_PROVIDER NOT NULL,
  availability_zone       TEXT,
  status                  STATUS DEFAULT 'Unknown' NOT NULL,
  cluster_id              INTEGER REFERENCES clusters(id) ON DELETE CASCADE NOT NULL,
  last_scan_ts            TIMESTAMP WITH TIME ZONE,
  created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  age                     INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  CONSTRAINT uq_instances_instance_id_cluster_id UNIQUE (instance_id, cluster_id),
  CONSTRAINT uq_instances_instance_id_provider   UNIQUE (instance_id, provider)
);

CREATE INDEX IF NOT EXISTS ix_instances_cluster               ON instances (cluster_id);
CREATE INDEX IF NOT EXISTS ix_instances_cluster_status        ON instances (cluster_id, status);
CREATE INDEX IF NOT EXISTS ix_instances_cluster_lastscan_desc ON instances (cluster_id, last_scan_ts DESC);
CREATE INDEX IF NOT EXISTS ix_instances_cluster_name          ON instances (cluster_id, lower(instance_name));
CREATE INDEX IF NOT EXISTS ix_instances_last_scan_active      ON instances (last_scan_ts) WHERE status <> 'Terminated';



-- ############################################################
-- Instances Tags
-- ############################################################
-- TODO Check if is more efficient to move this table to JSONB column on the Instances table
\! echo '## Creating Tags table'

CREATE TABLE IF NOT EXISTS tags (
  key                     TEXT NOT NULL,
  value                   TEXT,
  instance_id             BIGINT REFERENCES instances(id) ON DELETE CASCADE NOT NULL,
  PRIMARY KEY (key, instance_id)
);

CREATE INDEX IF NOT EXISTS ix_itags_instance      ON tags (instance_id);
CREATE INDEX IF NOT EXISTS ix_itags_key_val       ON tags (key, value);
CREATE INDEX IF NOT EXISTS ix_itags_key           ON tags (key);
CREATE INDEX IF NOT EXISTS ix_itags_key_lower_val ON tags (key, lower(value));



-- ############################################################
-- Instances expenses
-- ############################################################
\! echo '## Creating Expenses table'

CREATE TABLE IF NOT EXISTS expenses (
  instance_id              BIGINT REFERENCES instances(id) ON DELETE CASCADE NOT NULL,
  date                     DATE NOT NULL,
  amount                   NUMERIC(12,2) DEFAULT 0.0,
  PRIMARY KEY (instance_id, date),
  CONSTRAINT chk_expenses_amount_nonneg CHECK (amount >=0)
) PARTITION BY RANGE (date);

CREATE INDEX IF NOT EXISTS ix_expenses_instance           ON expenses (instance_id);
CREATE INDEX IF NOT EXISTS ix_expenses_date               ON expenses (date);
CREATE INDEX IF NOT EXISTS ix_expenses_instance_date_desc ON expenses (instance_id, date DESC);

-- Default expenses partition. The rest of expenses will be created by pg_cron
CREATE TABLE expenses_default PARTITION OF expenses DEFAULT;



-- ############################################################
-- Events
-- ############################################################
CREATE TABLE IF NOT EXISTS events (
  id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  event_timestamp         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  triggered_by            TEXT NOT NULL,
  action                  TEXT NOT NULL,
  resource_id             INTEGER,
  resource_type           TEXT NOT NULL,
  result                  TEXT NOT NULL,
  description             TEXT NULL,
  severity                TEXT DEFAULT 'info'::TEXT NOT NULL,
  CONSTRAINT events_resource_type_check CHECK ((resource_type = ANY (ARRAY['cluster'::TEXT, 'instance'::TEXT]))),
  PRIMARY KEY (id, event_timestamp)
) PARTITION BY RANGE (event_timestamp);

-- Events Indexes
CREATE INDEX IF NOT EXISTS ix_events_type_id_time ON events (resource_type, resource_id, event_timestamp DESC);

-- Default expenses partition. The rest of expenses will be created by pg_cron
CREATE TABLE events_default PARTITION OF events DEFAULT;



-- #############################################################################
-- ## Actions and Scheduling definition ##
-- #############################################################################
\! echo '## Creating Schedule table'

CREATE TABLE IF NOT EXISTS schedule (
  id                      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  type                    ACTION_TYPE NOT NULL,
  time                    TIMESTAMP WITH TIME ZONE,
  cron_exp                TEXT,
  operation               ACTION_OPERATION NOT NULL,
  target                  INTEGER REFERENCES clusters(id) ON DELETE CASCADE NOT NULL,
  status                  ACTION_STATUS DEFAULT 'Unknown' NOT NULL,
  enabled                 BOOLEAN DEFAULT false,
	PRIMARY KEY (id),
  CONSTRAINT chk_schedule_time_or_cron CHECK ((time IS NOT NULL) <> (cron_exp IS NOT NULL))
);

CREATE INDEX IF NOT EXISTS ix_schedule_target_enabled ON schedule (target, enabled);
CREATE INDEX IF NOT EXISTS ix_schedule_status         ON schedule (status);



-- ## Advanced Inventory definition (Views & Functions) ##
-- #################################################################################################

\! echo '## Creating Inventory Functions'

-- ############################################################
-- ## Accounts
-- ############################################################

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

CREATE MATERIALIZED VIEW m_accounts_full_view AS SELECT * FROM accounts_full_view;



-- ############################################################
-- ## Clusters
-- ############################################################

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
  c.console_link,
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

CREATE MATERIALIZED VIEW m_clusters_full_view AS SELECT * FROM clusters_full_view;

-- cluster tags view. Returns the cluster_id + every tag omitting repeated tags keys
CREATE view clusters_tags AS
SELECT
    c.cluster_id,
    t.key,
    MIN(t.value) AS value
FROM clusters   c
JOIN instances  i ON i.cluster_id = c.id
JOIN tags       t ON t.instance_id = i.id
GROUP BY c.cluster_id, t.key
HAVING COUNT(*) > 1;



-- #############################################################################
-- Instances
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

CREATE MATERIALIZED VIEW m_instances_full_view AS SELECT * FROM instances_full_view;

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

CREATE MATERIALIZED VIEW m_instances_full_view_with_tags AS SELECT * FROM instances_full_view_with_tags;

-- Instances pending for expense update
CREATE VIEW instances_pending_expense_update AS
SELECT
  a.account_id,
	i.instance_id
FROM
	instances i
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
	i.id = last_expenses.instance_id
JOIN clusters c ON c.id = i.cluster_id
JOIN accounts a ON a.id = c.account_id
WHERE
  i.status <> 'Terminated'
	AND (
    last_expenses.last_expense_date IS NULL
	  OR last_expenses.last_expense_date < CURRENT_DATE - INTERVAL '1 day'
  );




-- #############################################################################
-- ## Tags
-- #############################################################################




-- #############################################################################
-- ## Schedule
-- #############################################################################

-- Schedule with cluster and instances list view
CREATE VIEW schedule_full_view AS
SELECT
	s.id,
	s.type,
	s.time,
	s.cron_exp,
	s.operation,
	s.status,
	s.enabled,
	c.id AS cluster_id,
	c.region,
	c.account_id,
	COALESCE(
		array_agg(DISTINCT i.instance_id ORDER BY i.instance_id),
		'{}'
	) AS instances
FROM
	schedule s
JOIN clusters c ON c.id = s.target
LEFT JOIN instances i ON i.cluster_id = c.id
GROUP BY s.id, c.id
ORDER BY s.id;



-- #############################################################################
-- ## Expenses
-- #############################################################################

-- Function for creating expenses partitions
CREATE OR REPLACE FUNCTION create_next_month_expenses_partition()
RETURNS text
LANGUAGE plpgsql
AS $$
DECLARE
  next_month  date := (date_trunc('month', current_date)::date + interval '1 month')::date;
  start_date  date := next_month;
  end_date    date := (next_month + interval '1 month')::date;
  part_name   text := format('expenses_%s', to_char(next_month, 'YYYY_MM'));
BEGIN
  EXECUTE format(
    'CREATE TABLE IF NOT EXISTS %I PARTITION OF public.expenses
       FOR VALUES FROM (%L) TO (%L);',
    part_name, start_date, end_date
  );
  RETURN part_name;
END;
$$;



-- #############################################################################
-- ## Events
-- #############################################################################

-- View for Cluster Events
CREATE VIEW cluster_events AS
SELECT
  id,
  event_timestamp,
  triggered_by,
  action,
  resource_id,
  resource_type,
  result,
  description,
  severity
FROM events
ORDER BY event_timestamp DESC;

-- View for System Events
CREATE VIEW system_events AS
SELECT
  ev.id,
  ev.event_timestamp,
  ev.triggered_by,
  ev.action,
  ev.resource_id,
  ev.resource_type,
  ev.result,
  ev.description,
  ev.severity,
  acc.account_id,
  acc.provider
FROM events ev
LEFT JOIN accounts acc ON acc.id = (
  CASE
    WHEN ev.resource_type = 'cluster'
    THEN (SELECT c.account_id FROM clusters c WHERE c.id = ev.resource_id)
    WHEN ev.resource_type = 'instance'
    THEN (SELECT c.account_id FROM clusters c WHERE c.id = (SELECT i.cluster_id FROM instances i WHERE i.id = ev.resource_id))
  END
);



-- Function for creating events partitions
CREATE OR REPLACE FUNCTION create_next_month_events_partition()
RETURNS text
LANGUAGE plpgsql
AS $$
DECLARE
  next_month  date := (date_trunc('month', current_date)::date + interval '1 month')::date;
  start_date  date := next_month;
  end_date    date := (next_month + interval '1 month')::date;
  part_name   text := format('events_%s', to_char(next_month, 'YYYY_MM'));
BEGIN
  EXECUTE format(
    'CREATE TABLE IF NOT EXISTS %I PARTITION OF public.events
       FOR VALUES FROM (%L) TO (%L);',
    part_name, start_date, end_date
  );
  RETURN part_name;
END;
$$;





-- ## Maintenance Functions ##
-- #############################################################################

\! echo '## Creating Maintenance Functions'

-- Marks instances as 'Terminated' if they haven't been scanned in the last 24 hours
CREATE OR REPLACE FUNCTION check_terminated_instances()
RETURNS void AS $$
BEGIN
  UPDATE instances
  SET status = 'Terminated'
  WHERE status <> 'Terminated'
	  AND last_scan_ts < NOW() - INTERVAL '1 day';
END;
$$ LANGUAGE plpgsql;


-- Marks clusters as 'Terminated' if they haven't been scanned in the last 24 hours
CREATE OR REPLACE FUNCTION check_terminated_clusters()
RETURNS void AS $$
BEGIN
  UPDATE clusters
  SET status = 'Terminated'
  WHERE status <> 'Terminated'
	  AND last_scan_ts < NOW() - INTERVAL '1 day';
END;
$$ LANGUAGE plpgsql;

-- Runs check_terminated_clusters/instances functions
CREATE OR REPLACE FUNCTION check_terminated_inventory()
RETURNS void AS $$
BEGIN
  PERFORM check_terminated_clusters();
  PERFORM check_terminated_instances();
END;
$$ LANGUAGE plpgsql;

-- Updates every materialized view
CREATE OR REPLACE FUNCTION refresh_materialized_views()
RETURNS void AS $$
BEGIN
  REFRESH MATERIALIZED VIEW m_accounts_full_view;
  REFRESH MATERIALIZED VIEW m_clusters_full_view;
  REFRESH MATERIALIZED VIEW m_instances_full_view;
  REFRESH MATERIALIZED VIEW m_instances_full_view_with_tags;
END;
$$ LANGUAGE plpgsql;
