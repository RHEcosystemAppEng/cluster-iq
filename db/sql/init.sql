-- ## Tables definitions ##
-- Cloud Providers
CREATE TABLE IF NOT EXISTS providers (
  name TEXT PRIMARY KEY
);

-- Default values for Cloud Providers table
INSERT INTO
  providers(name)
VALUES
  ('AWS'),
  ('GCP'),
  ('Azure'),
  ('UNKNOWN')
;


-- Action Operations
CREATE TABLE IF NOT EXISTS action_operations (
  name TEXT PRIMARY KEY
);

-- Default values for Cloud Providers table
INSERT INTO
  action_operations(name)
VALUES
  ('PowerOnCluster'),
  ('PowerOffCluster')
;

-- Status
CREATE TABLE IF NOT EXISTS status (
  value TEXT PRIMARY KEY
);

-- Default values for Status table
INSERT INTO
  status(value)
VALUES
  ('Running'),
  ('Stopped'),
  ('Terminated')
;


-- Accounts
CREATE TABLE IF NOT EXISTS accounts (
  id TEXT,
  name TEXT PRIMARY KEY,
  provider TEXT REFERENCES providers(name),
  cluster_count INTEGER,
  last_scan_timestamp TIMESTAMP WITH TIME ZONE,
  total_cost REAL,
  last_15_days_cost REAL,
  last_month_cost REAL,
  current_month_so_far_cost REAL
);


-- Clusters
CREATE TABLE IF NOT EXISTS clusters (
  -- id is the result of joining: "name+infra_id+account"
  id TEXT PRIMARY KEY,
  name TEXT,
  infra_id TEXT,
  provider TEXT REFERENCES providers(name),
  status TEXT REFERENCES status(value),
  region TEXT,
  account_name TEXT REFERENCES accounts(name),
  console_link TEXT,
  instance_count INTEGER,
  last_scan_timestamp TIMESTAMP WITH TIME ZONE,
  creation_timestamp TIMESTAMP WITH TIME ZONE,
  age INT,
  owner TEXT,
  total_cost REAL,
  last_15_days_cost REAL,
  last_month_cost REAL,
  current_month_so_far_cost REAL
);


-- Instances
CREATE TABLE IF NOT EXISTS instances (
  id TEXT PRIMARY KEY,
  name TEXT,
  provider TEXT REFERENCES providers(name),
  instance_type TEXT,
  availability_zone TEXT,
  status TEXT REFERENCES status(value),
  cluster_id TEXT REFERENCES clusters(id),
  last_scan_timestamp TIMESTAMP WITH TIME ZONE,
  creation_timestamp TIMESTAMP WITH TIME ZONE,
  age INT,
  daily_cost REAL,
  total_cost REAL
);


-- Instances Tags
CREATE TABLE IF NOT EXISTS tags (
  key TEXT,
  value TEXT,
  instance_id TEXT REFERENCES instances(id),
  PRIMARY KEY (key, instance_id)
);


-- Instances expenses
CREATE TABLE IF NOT EXISTS expenses (
  instance_id TEXT REFERENCES instances(id),
  date DATE,
  amount REAL,
  PRIMARY KEY (instance_id, date)
);

-- Action types table
CREATE TABLE IF NOT EXISTS action_types (
  name TEXT PRIMARY KEY
);

-- Default values for Action Types
INSERT INTO
  action_types(name)
VALUES
  ('cron_action'),
  ('scheduled_action')
;

-- Action Status table
CREATE TABLE IF NOT EXISTS action_status (
  name TEXT PRIMARY KEY
);

-- Default values for Action Types
INSERT INTO
  action_status(name)
VALUES
  ('Success'),
  ('Failed'),
  ('Pending'),
  ('Unknown')
;

-- Scheduled actions
CREATE TABLE IF NOT EXISTS schedule (
  id SERIAL PRIMARY KEY,
  type TEXT REFERENCES action_types(name),
  time TIMESTAMP WITH TIME ZONE,
  cron_exp TEXT,
  operation TEXT REFERENCES action_operations(name),
  target TEXT REFERENCES clusters(id),
  status TEXT REFERENCES action_status(name),
  enabled BOOLEAN
);


-- Audit logs
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
  event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  triggered_by text NOT NULL,
  action_name text NOT NULL,
  resource_id text NOT NULL,
  resource_type text NOT NULL,
  result text NOT NULL,
  description text NULL,
  severity text DEFAULT 'info'::text NOT NULL,
  CONSTRAINT audit_logs_pkey PRIMARY KEY (id),
  CONSTRAINT audit_logs_resource_type_check CHECK ((resource_type = ANY (ARRAY['cluster'::text, 'instance'::text])))
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
  SET total_cost = (SELECT SUM(amount) FROM expenses WHERE instance_id = NEW.instance_id)
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
  SET total_cost = (SELECT SUM(amount) FROM expenses WHERE instance_id = OLD.instance_id)
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
  SET daily_cost = (SELECT SUM(amount)/count(*) FROM expenses WHERE instance_id = NEW.instance_id)
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
  SET daily_cost = (SELECT SUM(amount)/count(*) FROM expenses WHERE instance_id = OLD.instance_id)
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
    total_cost = (SELECT SUM(total_cost) FROM instances WHERE cluster_id = NEW.cluster_id),
    last_15_days_cost = (SELECT SUM(expenses.amount) FROM instances JOIN expenses ON instances.id = expenses.instance_id WHERE instances.cluster_id = NEW.cluster_id AND expenses.date >= NOW()::date - interval '15 day'),
    last_month_cost = (SELECT SUM(expenses.amount) FROM instances JOIN expenses ON instances.id = expenses.instance_id WHERE instances.cluster_id = NEW.cluster_id AND (EXTRACT(MONTH FROM NOW()::date - interval '1 month') = EXTRACT(MONTH FROM expenses.date))),
    current_month_so_far_cost = (SELECT SUM(expenses.amount) FROM instances JOIN expenses ON instances.id = expenses.instance_id WHERE instances.cluster_id = NEW.cluster_id AND (EXTRACT(MONTH FROM NOW()::date) = EXTRACT(MONTH FROM expenses.date)))
  WHERE id = NEW.cluster_id;
  RETURN NEW;
END;
$$;

-- Updates the total cost of a cluster based on its associated instances
-- SELECT SUM(expenses.amount) FROM instances JOIN expenses ON instances.id = expenses.instance_id WHERE instances.cluster_id='***********' AND expenses.date >= NOW()::date - interval '15 day';
CREATE OR REPLACE FUNCTION update_cluster_cost_info()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE clusters
  SET
    last_15_days_cost = (SELECT SUM(expenses.amount) FROM instances JOIN expenses ON instances.id = expenses.instance_id WHERE instances.cluster_id = NEW.cluster_id AND expenses.date >= NOW()::date - interval '15 day'),
    last_month_cost = (SELECT SUM(expenses.amount) FROM instances JOIN expenses ON instances.id = expenses.instance_id WHERE instances.cluster_id = NEW.cluster_id AND (EXTRACT(MONTH FROM NOW()::date - interval '1 month') = EXTRACT(MONTH FROM expenses.date))),
    current_month_so_far_cost = (SELECT SUM(expenses.amount) FROM instances JOIN expenses ON instances.id = expenses.instance_id WHERE instances.cluster_id = NEW.cluster_id AND (EXTRACT(MONTH FROM NOW()::date) = EXTRACT(MONTH FROM expenses.date)))
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
    total_cost = (SELECT SUM(clusters.total_cost) FROM clusters WHERE account_name = NEW.account_name),
    last_15_days_cost = (SELECT SUM(clusters.last_15_days_cost) FROM clusters WHERE account_name = NEW.account_name),
    last_month_cost = (SELECT SUM(clusters.last_month_cost) FROM clusters WHERE account_name = NEW.account_name),
    current_month_so_far_cost = (SELECT SUM(clusters.current_month_so_far_cost) FROM clusters WHERE account_name = NEW.account_name)
  WHERE name = NEW.account_name;
  RETURN NEW;
END;
$$;

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

-- Trigger to update cluster costs info
CREATE TRIGGER update_cluster_cost_info
AFTER UPDATE
ON instances
FOR EACH ROW
  EXECUTE PROCEDURE update_cluster_cost_info();


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
--
