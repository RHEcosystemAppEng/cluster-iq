

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
