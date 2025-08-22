-- ## Inventory definition ##
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


-- Accounts
CREATE TABLE IF NOT EXISTS accounts (
  id                  BIGSERIAL PRIMARY KEY,        -- Internal Account Name
  account_id          TEXT NOT NULL,
  name                TEXT NOT NULL,
  provider            CLOUD_PROVIDER,
  cluster_count       INTEGER,
  last_scan_timestamp TIMESTAMP WITH TIME ZONE,
  total_cost NUMERIC(12,2) DEFAULT 0.0,
  last_15_days_cost NUMERIC(12,2) DEFAULT 0.0,
  last_month_cost NUMERIC(12,2) DEFAULT 0.0,
  current_month_so_far_cost NUMERIC(12,2) DEFAULT 0.0,
  PRIMARY KEY (name)
);


-- Clusters
CREATE TABLE IF NOT EXISTS clusters (
  -- id is the result of joining: "name+infra_id+account"
  id TEXT,
  name TEXT,
  infra_id TEXT,
  provider CLOUD_PROVIDER,
  status STATUS,
  region TEXT,
  account_name TEXT REFERENCES accounts(name) ON DELETE CASCADE,
  console_link TEXT,
  instance_count INTEGER,
  last_scan_timestamp TIMESTAMP WITH TIME ZONE,
  creation_timestamp TIMESTAMP WITH TIME ZONE,
  age INT,
  owner TEXT,
  total_cost NUMERIC(12,2) DEFAULT 0.0,
  last_15_days_cost NUMERIC(12,2) DEFAULT 0.0,
  last_month_cost NUMERIC(12,2) DEFAULT 0.0,
  current_month_so_far_cost NUMERIC(12,2) DEFAULT 0.0,
  PRIMARY KEY (id, account_name)
) PARTITION BY HASH (account_name);

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
  id TEXT,
  name TEXT,
  provider CLOUD_PROVIDER,
  instance_type TEXT,
  availability_zone TEXT,
  status STATUS,
  cluster_id TEXT REFERENCES clusters(id) ON DELETE CASCADE,
  last_scan_timestamp TIMESTAMP WITH TIME ZONE,
  creation_timestamp TIMESTAMP WITH TIME ZONE,
  age INT,
  daily_cost NUMERIC(12,2) DEFAULT 0.0,
  total_cost NUMERIC(12,2) DEFAULT 0.0,
  PRIMARY KEY (id, cluster_id)
) PARTITION BY HASH(cluster_id);

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
CREATE TABLE IF NOT EXISTS tags (
  key TEXT,
  value TEXT,
  instance_id TEXT REFERENCES instances(id) ON DELETE CASCADE,
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
CREATE TABLE IF NOT EXISTS expenses (
  instance_id TEXT REFERENCES instances(id) ON DELETE CASCADE,
  date DATE,
  amount NUMERIC(12,2) DEFAULT 0.0,
  PRIMARY KEY (instance_id, date)
) PARTITION BY RANGE (date);


-- ## Actions Definitions ##
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
  target TEXT REFERENCES clusters(id) ON DELETE CASCADE,
  status ACTION_STATUS NOT NULL DEFAULT 'Unknown',
  enabled BOOLEAN DEFAULT false
);


-- Audit logs
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  triggered_by TEXT NOT NULL,
  action_name TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  resource_type RESOURCE_TYPE NOT NULL,
  result TEXT NOT NULL,
  description TEXT NULL,
  severity TEXT DEFAULT 'info'::TEXT NOT NULL,
  PRIMARY KEY (id)
);


-- ## Functions ##
-- Updates the total cost of an instance after a new expense record is inserted
CREATE OR REPLACE FUNCTION update_instance_total_costs_after_insert()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE instances
  SET
    total_cost = (
      SELECT SUM(amount)
      FROM expenses
      WHERE instance_id = NEW.instance_id
    )
    WHERE id = NEW.instance_id;
  RETURN NEW;
END;
$$;


-- Updates the total cost of an instance after an expense record is deleted
CREATE OR REPLACE FUNCTION update_instance_total_costs_after_delete()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE instances
  SET
    total_cost = (
      SELECT SUM(amount)
      FROM expenses
      WHERE instance_id = OLD.instance_id
    )
    WHERE id = OLD.instance_id;
  RETURN OLD;
END;
$$;


-- Updates the daily cost of an instance after a new expense record is inserted
CREATE OR REPLACE FUNCTION update_instance_daily_costs_after_insert()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE instances
  SET
    daily_cost = (
      SELECT COALESCE(SUM(amount)/NULLIF(COUNT(*), 0), 0)
      FROM expenses
      WHERE instance_id = NEW.instance_id
    )
    WHERE id = NEW.instance_id;
  RETURN NEW;
END;
$$;


-- Updates the daily cost of an instance after an expense record is deleted
CREATE OR REPLACE FUNCTION update_instance_daily_costs_after_delete()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE instances
  SET
    daily_cost = (
      SELECT COALESCE(SUM(amount)/NULLIF(COUNT(*), 0), 0)
      FROM expenses
      WHERE instance_id = NEW.instance_id
    )
    WHERE id = OLD.instance_id;
  RETURN OLD;
END;
$$;


-- Updates the total cost of a cluster based on its associated instances
CREATE OR REPLACE FUNCTION update_cluster_cost_info()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE clusters
  SET
    total_cost = (
      SELECT COALESCE(SUM(total_cost), 0) as sum
      FROM instances
      WHERE cluster_id = NEW.cluster_id
    ),
    last_15_days_cost = (
      SELECT COALESCE(SUM(expenses.amount), 0)
      FROM instances
      JOIN expenses ON instances.id = expenses.instance_id
      WHERE instances.cluster_id = NEW.cluster_id
        AND expenses.date >= NOW()::date - interval '15 day'
    ),
    last_month_cost = (
      SELECT COALESCE(SUM(expenses.amount), 0)
      FROM instances
      JOIN expenses ON instances.id = expenses.instance_id
      WHERE instances.cluster_id = NEW.cluster_id
        AND EXTRACT(YEAR FROM NOW()::date - interval '1 month') = EXTRACT(YEAR FROM expenses.date)
        AND EXTRACT(MONTH FROM NOW()::date - interval '1 month') = EXTRACT(MONTH FROM expenses.date)
    ),
    current_month_so_far_cost = (
      SELECT COALESCE(SUM(expenses.amount), 0)
      FROM instances
      JOIN expenses ON instances.id = expenses.instance_id
      WHERE instances.cluster_id = NEW.cluster_id
        AND (EXTRACT(MONTH FROM NOW()::date) = EXTRACT(MONTH FROM expenses.date)
      )
    )
    WHERE id = NEW.cluster_id;
  RETURN NEW;
END;
$$;


-- Updates the total cost of an account based on its associated clusters
CREATE OR REPLACE FUNCTION update_account_cost_info()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE accounts
  SET
    total_cost = (
      SELECT COALESCE(SUM(clusters.total_cost), 0)
      FROM clusters
      WHERE account_name = NEW.account_name
    ),
    last_15_days_cost = (
      SELECT COALESCE(SUM(clusters.last_15_days_cost), 0)
      FROM clusters
      WHERE account_name = NEW.account_name
    ),
    last_month_cost = (
      SELECT COALESCE(SUM(clusters.last_month_cost), 0)
      FROM clusters
      WHERE account_name = NEW.account_name
    ),
    current_month_so_far_cost = (
      SELECT COALESCE(SUM(clusters.current_month_so_far_cost), 0)
      FROM clusters
      WHERE account_name = NEW.account_name
    )
    WHERE name = NEW.account_name;
  RETURN NEW;
END;
$$;


-- ## Maintenance Functions ##
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


-- ## Triggers ##
-- Trigger to update instance total cost after an expense is inserted
CREATE TRIGGER update_instance_total_cost_after_insert
AFTER INSERT
ON expenses
FOR EACH ROW
  EXECUTE PROCEDURE update_instance_total_costs_after_insert();


-- Trigger to update instance total cost after an expense is updated
CREATE TRIGGER update_instance_total_cost_after_update
AFTER UPDATE
ON expenses
FOR EACH ROW
  EXECUTE PROCEDURE update_instance_total_costs_after_insert();


-- Trigger to update instance total cost after an expense is deleted
CREATE TRIGGER update_instance_total_cost_after_delete
AFTER DELETE
ON expenses
FOR EACH ROW
  EXECUTE PROCEDURE update_instance_total_costs_after_delete();


-- Trigger to update instance daily cost after an expense is inserted
CREATE TRIGGER update_instance_daily_cost_after_insert
AFTER INSERT
ON expenses
FOR EACH ROW
  EXECUTE PROCEDURE update_instance_daily_costs_after_insert();


-- Trigger to update instance daily cost after an expense is updated
CREATE TRIGGER update_instance_daily_cost_after_update
AFTER UPDATE
ON expenses
FOR EACH ROW
  EXECUTE PROCEDURE update_instance_daily_costs_after_insert();


-- Trigger to update instance daily cost after an expense is deleted
CREATE TRIGGER update_instance_daily_cost_after_delete
AFTER DELETE
ON expenses
FOR EACH ROW
  EXECUTE PROCEDURE update_instance_daily_costs_after_delete();


-- Trigger to update cluster costs info
CREATE TRIGGER update_cluster_cost_info
AFTER UPDATE
ON instances
FOR EACH ROW
  EXECUTE PROCEDURE update_cluster_cost_info();


-- Trigger to update account total cost after a cluster is updated
CREATE TRIGGER update_account_cost_info
AFTER UPDATE
ON clusters
FOR EACH ROW
  EXECUTE PROCEDURE update_account_cost_info();
