-- Drop Triggers
DROP TRIGGER update_instance_total_cost_after_insert ON expenses;
DROP TRIGGER update_instance_total_cost_after_delete ON expenses;
DROP TRIGGER update_instance_daily_cost_after_insert ON expenses;
DROP TRIGGER update_instance_daily_cost_after_delete ON expenses;
DROP TRIGGER update_cluster_total_cost ON instances;

-- Drop Functins
DROP FUNCTION update_instance_total_costs_after_insert;
DROP FUNCTION update_instance_total_costs_after_delete;
DROP FUNCTION update_instance_daily_costs_after_insert;
DROP FUNCTION update_instance_daily_costs_after_delete;
DROP FUNCTION update_cluster_total_costs;

-- Drop tables
DROP TABLE tags;
DROP TABLE expenses;
DROP TABLE instances;
DROP TABLE clusters;
DROP TABLE accounts;
DROP TABLE providers;
DROP TABLE status;
