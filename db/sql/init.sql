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
  ('Terminated'),
  ('Unknown')
;


-- Accounts
CREATE TABLE IF NOT EXISTS accounts (
  id TEXT,
  name TEXT PRIMARY KEY,
  provider TEXT REFERENCES providers(name),
  total_cost REAL,
  cluster_count INTEGER,
  last_scan_timestamp TIMESTAMP WITH TIME ZONE
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
  total_cost REAL
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
CREATE OR REPLACE FUNCTION update_cluster_total_costs()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE clusters
  SET total_cost = (SELECT SUM(total_cost) FROM instances WHERE cluster_id = NEW.cluster_id)
  WHERE id = NEW.cluster_id;
  RETURN NEW;
END;
$$;

-- Updates the total cost of an account based on its associated clusters
CREATE OR REPLACE FUNCTION update_account_total_costs()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  UPDATE accounts
  SET total_cost = (SELECT SUM(total_cost) FROM clusters WHERE account_name = NEW.account_name)
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

-- Trigger to update cluster total cost after an instance is updated
CREATE TRIGGER update_cluster_total_cost
AFTER UPDATE
ON instances
FOR EACH ROW
  EXECUTE PROCEDURE update_cluster_total_costs();

-- Trigger to update account total cost after a cluster is updated
CREATE TRIGGER update_account_total_cost
AFTER UPDATE
ON clusters
FOR EACH ROW
  EXECUTE PROCEDURE update_account_total_costs();

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
