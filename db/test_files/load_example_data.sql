-- Demo seed for ClusterIQ schema
-- psql postgresql://user:password@pgsql:5432/clusteriq < load_example_data.sql
BEGIN;

-- Clean previous data
TRUNCATE action_runs, schedule, targets, expenses, tags, instances, clusters, accounts RESTART IDENTITY CASCADE;

-- Drop existing expense partitions to avoid conflicts
DO $$
DECLARE
  r RECORD;
BEGIN
  FOR r IN
    SELECT inhrelid::regclass AS part_name
    FROM pg_inherits
    WHERE inhparent = 'expenses'::regclass
  LOOP
    EXECUTE format('DROP TABLE IF EXISTS %I CASCADE', r.part_name);
  END LOOP;
END
$$;

-- Create expense partitions for 3 months
DO $$
DECLARE
  cur_start  DATE := date_trunc('month', current_date)::date;
  cur_end    DATE := (cur_start + INTERVAL '1 month')::date;
  prev_start DATE := date_trunc('month', current_date - INTERVAL '1 month')::date;
  prev_end   DATE := cur_start;
  prev2_start DATE := date_trunc('month', current_date - INTERVAL '2 month')::date;
  prev2_end   DATE := prev_start;
  part_name  TEXT;
BEGIN
  part_name := format('expenses_%s', to_char(prev2_start, 'YYYY_MM'));
  IF to_regclass(part_name) IS NULL THEN
    EXECUTE format('CREATE TABLE %I PARTITION OF expenses FOR VALUES FROM (%L) TO (%L)', part_name, prev2_start, prev2_end);
  END IF;

  part_name := format('expenses_%s', to_char(prev_start, 'YYYY_MM'));
  IF to_regclass(part_name) IS NULL THEN
    EXECUTE format('CREATE TABLE %I PARTITION OF expenses FOR VALUES FROM (%L) TO (%L)', part_name, prev_start, prev_end);
  END IF;

  part_name := format('expenses_%s', to_char(cur_start, 'YYYY_MM'));
  IF to_regclass(part_name) IS NULL THEN
    EXECUTE format('CREATE TABLE %I PARTITION OF expenses FOR VALUES FROM (%L) TO (%L)', part_name, cur_start, cur_end);
  END IF;
END
$$;

-- ============================================================================
-- Accounts: 5 accounts across 3 providers
-- ============================================================================
INSERT INTO accounts (account_id, account_name, provider, last_scan_ts) VALUES
  ('111111111111', 'rh-engineering-prod',  'AWS',   now() - INTERVAL '2 hours'),
  ('222222222222', 'rh-engineering-dev',   'AWS',   now() - INTERVAL '2 hours'),
  ('333333333333', 'rh-qe-staging',       'AWS',   now() - INTERVAL '2 hours'),
  ('gcp-proj-001', 'rh-platform-gcp',     'GCP',   now() - INTERVAL '2 hours'),
  ('azure-sub-01', 'rh-services-azure',   'Azure', now() - INTERVAL '2 hours');

-- ============================================================================
-- Main data generation
-- ============================================================================
DO $$
DECLARE
  v_acc_id     INT;
  v_clu_pk     BIGINT;
  v_ins_pk     BIGINT;

  -- Cluster definition arrays
  v_name       TEXT;
  v_infra      TEXT;
  v_region     TEXT;
  v_status     STATUS;
  v_owner      TEXT;
  v_provider   CLOUD_PROVIDER;
  v_age        INT;
  v_partner    TEXT;

  -- Instance vars
  v_ins_name   TEXT;
  v_ins_type   TEXT;
  v_az         TEXT;
  v_ins_status STATUS;

  -- Expense vars
  d            DATE;
  base_cost    NUMERIC(12,2);
  amt          NUMERIC(12,2);

  -- Tag key pools
  partners     TEXT[] := ARRAY['Red Hat', 'Accenture', 'IBM', 'Deloitte', 'Wipro', 'Infosys', 'TCS'];
  owners       TEXT[] := ARRAY['jsmith@redhat.com', 'agarcia@redhat.com', 'mchen@redhat.com', 'pjones@redhat.com', 'lbrown@redhat.com', 'klee@redhat.com', 'ssingh@redhat.com', 'twilson@redhat.com'];
  teams        TEXT[] := ARRAY['Platform', 'SRE', 'QE', 'Performance', 'Security', 'DevOps', 'Middleware'];
  envs         TEXT[] := ARRAY['production', 'staging', 'development', 'qa', 'perf-test', 'sandbox'];

  -- AWS regions
  aws_regions  TEXT[] := ARRAY['us-east-1', 'us-east-2', 'us-west-2', 'eu-west-1', 'eu-central-1', 'ap-southeast-1'];
  -- GCP regions
  gcp_regions  TEXT[] := ARRAY['us-central1', 'europe-west1', 'europe-west3', 'asia-east1'];
  -- Azure regions
  az_regions   TEXT[] := ARRAY['eastus', 'westeurope', 'northeurope', 'southeastasia'];

  -- Instance type pools
  aws_types    TEXT[] := ARRAY['m5.xlarge', 'm5.2xlarge', 'r5.xlarge', 'r5.2xlarge', 'c5.2xlarge', 'c5.4xlarge', 'm6i.xlarge', 'm6i.2xlarge'];
  gcp_types    TEXT[] := ARRAY['e2-standard-4', 'e2-standard-8', 'n2-standard-4', 'n2-standard-8', 'n2-highmem-4'];
  az_types     TEXT[] := ARRAY['Standard_D4s_v3', 'Standard_D8s_v3', 'Standard_E4s_v3', 'Standard_E8s_v3'];

  -- Cluster definitions: (account_index, name, status, region_index, owner_index, partner_index, base_daily_cost)
  -- We'll generate these programmatically per account
  n_instances  INT;
  cost_base    NUMERIC(12,2);

BEGIN
  -- ========================================================================
  -- Account 1: rh-engineering-prod (AWS) — 6 clusters, heavy usage
  -- ========================================================================
  SELECT id INTO v_acc_id FROM accounts WHERE account_id = '111111111111';
  v_provider := 'AWS';

  -- Cluster 1: Large production cluster, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-prod-east', 'ocp-prod-east-abc12', 'abc12', v_provider, 'Running', 'us-east-1', v_acc_id, 'https://console-openshift-console.apps.ocp-prod-east.example.com', now() - INTERVAL '1 hour', now() - INTERVAL '120 days', 120, 'jsmith@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..8 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-prod-east', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'us-east-1' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '120 days', 120)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'jsmith@redhat.com', v_ins_pk), ('Partner', 'Red Hat', v_ins_pk), ('Team', 'Platform', v_ins_pk), ('Environment', 'production', v_ins_pk);
    cost_base := 3.50 + random() * 2.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.50, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 2: Staging cluster, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-staging-east', 'ocp-staging-east-def34', 'def34', v_provider, 'Running', 'us-east-1', v_acc_id, 'https://console-openshift-console.apps.ocp-staging-east.example.com', now() - INTERVAL '1 hour', now() - INTERVAL '90 days', 90, 'agarcia@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..5 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-staging-east', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'us-east-1' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '90 days', 90)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'agarcia@redhat.com', v_ins_pk), ('Partner', 'Accenture', v_ins_pk), ('Team', 'SRE', v_ins_pk), ('Environment', 'staging', v_ins_pk);
    cost_base := 2.20 + random() * 1.5;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.30, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 3: Stopped weekend cluster
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-weekend-west', 'ocp-weekend-west-ghi56', 'ghi56', v_provider, 'Stopped', 'us-west-2', v_acc_id, 'https://console-openshift-console.apps.ocp-weekend-west.example.com', now() - INTERVAL '3 hours', now() - INTERVAL '60 days', 60, 'mchen@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..4 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-weekend-west', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'us-west-2' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Stopped', now() - INTERVAL '3 hours', now() - INTERVAL '60 days', 60)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'mchen@redhat.com', v_ins_pk), ('Partner', 'IBM', v_ins_pk), ('Team', 'QE', v_ins_pk), ('Environment', 'development', v_ins_pk);
    cost_base := 1.80 + random() * 1.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.10, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 4: EU production cluster, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-prod-eu', 'ocp-prod-eu-jkl78', 'jkl78', v_provider, 'Running', 'eu-west-1', v_acc_id, 'https://console-openshift-console.apps.ocp-prod-eu.example.com', now() - INTERVAL '1 hour', now() - INTERVAL '200 days', 200, 'pjones@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..6 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-prod-eu', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'eu-west-1' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '200 days', 200)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'pjones@redhat.com', v_ins_pk), ('Partner', 'Deloitte', v_ins_pk), ('Team', 'Security', v_ins_pk), ('Environment', 'production', v_ins_pk);
    cost_base := 4.00 + random() * 2.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.80, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 5: Terminated old cluster
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-legacy-east', 'ocp-legacy-east-mno90', 'mno90', v_provider, 'Terminated', 'us-east-2', v_acc_id, '', now() - INTERVAL '30 days', now() - INTERVAL '365 days', 365, 'lbrown@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..3 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-legacy-east', j);
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, 'm5.xlarge', 'us-east-2a', 'Terminated', now() - INTERVAL '30 days', now() - INTERVAL '365 days', 365)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'lbrown@redhat.com', v_ins_pk), ('Partner', 'Red Hat', v_ins_pk), ('Team', 'Platform', v_ins_pk), ('Environment', 'production', v_ins_pk);
  END LOOP;

  -- Cluster 6: Frankfurt perf-test cluster, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-perf-fra', 'ocp-perf-fra-pqr12', 'pqr12', v_provider, 'Running', 'eu-central-1', v_acc_id, 'https://console-openshift-console.apps.ocp-perf-fra.example.com', now() - INTERVAL '1 hour', now() - INTERVAL '45 days', 45, 'ssingh@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..6 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-perf-fra', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'eu-central-1' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '45 days', 45)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'ssingh@redhat.com', v_ins_pk), ('Partner', 'Wipro', v_ins_pk), ('Team', 'Performance', v_ins_pk), ('Environment', 'perf-test', v_ins_pk);
    cost_base := 5.00 + random() * 3.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(1.00, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- ========================================================================
  -- Account 2: rh-engineering-dev (AWS) — 4 clusters, moderate usage
  -- ========================================================================
  SELECT id INTO v_acc_id FROM accounts WHERE account_id = '222222222222';

  -- Cluster 7: Dev cluster, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-dev-east', 'ocp-dev-east-stu34', 'stu34', v_provider, 'Running', 'us-east-1', v_acc_id, 'https://console-openshift-console.apps.ocp-dev-east.example.com', now() - INTERVAL '1 hour', now() - INTERVAL '30 days', 30, 'klee@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..4 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-dev-east', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'us-east-1' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '30 days', 30)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'klee@redhat.com', v_ins_pk), ('Partner', 'Infosys', v_ins_pk), ('Team', 'DevOps', v_ins_pk), ('Environment', 'development', v_ins_pk);
    cost_base := 1.50 + random() * 1.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.20, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 8: Sandbox, stopped
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-sandbox', 'ocp-sandbox-vwx56', 'vwx56', v_provider, 'Stopped', 'us-east-2', v_acc_id, 'https://console-openshift-console.apps.ocp-sandbox.example.com', now() - INTERVAL '5 hours', now() - INTERVAL '15 days', 15, 'twilson@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..3 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-sandbox', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'us-east-2' || chr(97 + (j % 2));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Stopped', now() - INTERVAL '5 hours', now() - INTERVAL '15 days', 15)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'twilson@redhat.com', v_ins_pk), ('Partner', 'TCS', v_ins_pk), ('Team', 'QE', v_ins_pk), ('Environment', 'sandbox', v_ins_pk);
    cost_base := 0.80 + random() * 0.5;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.05, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 9: CI cluster, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-ci-west', 'ocp-ci-west-yza78', 'yza78', v_provider, 'Running', 'us-west-2', v_acc_id, 'https://console-openshift-console.apps.ocp-ci-west.example.com', now() - INTERVAL '1 hour', now() - INTERVAL '75 days', 75, 'jsmith@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..5 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-ci-west', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'us-west-2' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '75 days', 75)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'jsmith@redhat.com', v_ins_pk), ('Partner', 'Red Hat', v_ins_pk), ('Team', 'SRE', v_ins_pk), ('Environment', 'qa', v_ins_pk);
    cost_base := 2.50 + random() * 1.5;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.30, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 10: Terminated dev cluster
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-dev-old', 'ocp-dev-old-bcd90', 'bcd90', v_provider, 'Terminated', 'us-east-1', v_acc_id, '', now() - INTERVAL '45 days', now() - INTERVAL '180 days', 180, 'agarcia@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..3 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-dev-old', j);
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, 'm5.xlarge', 'us-east-1a', 'Terminated', now() - INTERVAL '45 days', now() - INTERVAL '180 days', 180)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'agarcia@redhat.com', v_ins_pk), ('Partner', 'Accenture', v_ins_pk), ('Team', 'DevOps', v_ins_pk), ('Environment', 'development', v_ins_pk);
  END LOOP;

  -- ========================================================================
  -- Account 3: rh-qe-staging (AWS) — 3 clusters
  -- ========================================================================
  SELECT id INTO v_acc_id FROM accounts WHERE account_id = '333333333333';

  -- Cluster 11: QE main cluster, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-qe-main', 'ocp-qe-main-efg12', 'efg12', v_provider, 'Running', 'us-east-1', v_acc_id, 'https://console-openshift-console.apps.ocp-qe-main.example.com', now() - INTERVAL '1 hour', now() - INTERVAL '100 days', 100, 'mchen@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..6 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-qe-main', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'us-east-1' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '100 days', 100)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'mchen@redhat.com', v_ins_pk), ('Partner', 'IBM', v_ins_pk), ('Team', 'QE', v_ins_pk), ('Environment', 'qa', v_ins_pk);
    cost_base := 2.80 + random() * 1.5;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.40, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 12: QE nightly, stopped
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-qe-nightly', 'ocp-qe-nightly-hij34', 'hij34', v_provider, 'Stopped', 'us-west-2', v_acc_id, 'https://console-openshift-console.apps.ocp-qe-nightly.example.com', now() - INTERVAL '6 hours', now() - INTERVAL '50 days', 50, 'lbrown@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..4 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-qe-nightly', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'us-west-2' || chr(97 + (j % 2));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Stopped', now() - INTERVAL '6 hours', now() - INTERVAL '50 days', 50)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'lbrown@redhat.com', v_ins_pk), ('Partner', 'Deloitte', v_ins_pk), ('Team', 'QE', v_ins_pk), ('Environment', 'qa', v_ins_pk);
    cost_base := 1.20 + random() * 0.8;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.10, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 13: APAC QE cluster, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('ocp-qe-apac', 'ocp-qe-apac-klm56', 'klm56', v_provider, 'Running', 'ap-southeast-1', v_acc_id, 'https://console-openshift-console.apps.ocp-qe-apac.example.com', now() - INTERVAL '2 hours', now() - INTERVAL '25 days', 25, 'ssingh@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..4 LOOP
    v_ins_name := format('i-%s-node-%s', 'ocp-qe-apac', j);
    v_ins_type := aws_types[1 + floor(random()*array_length(aws_types,1))::INT];
    v_az := 'ap-southeast-1' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '2 hours', now() - INTERVAL '25 days', 25)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'ssingh@redhat.com', v_ins_pk), ('Partner', 'Wipro', v_ins_pk), ('Team', 'Performance', v_ins_pk), ('Environment', 'perf-test', v_ins_pk);
    cost_base := 2.00 + random() * 1.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.25, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- ========================================================================
  -- Account 4: rh-platform-gcp (GCP) — 3 clusters
  -- ========================================================================
  SELECT id INTO v_acc_id FROM accounts WHERE account_id = 'gcp-proj-001';
  v_provider := 'GCP';

  -- Cluster 14: GCP production, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('gke-prod-eu', 'gke-prod-eu-nop78', 'nop78', v_provider, 'Running', 'europe-west1', v_acc_id, 'https://console.cloud.google.com/kubernetes/clusters/gke-prod-eu', now() - INTERVAL '1 hour', now() - INTERVAL '150 days', 150, 'pjones@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..5 LOOP
    v_ins_name := format('gke-%s-node-%s', 'prod-eu', j);
    v_ins_type := gcp_types[1 + floor(random()*array_length(gcp_types,1))::INT];
    v_az := 'europe-west1-' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '150 days', 150)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'pjones@redhat.com', v_ins_pk), ('Partner', 'Red Hat', v_ins_pk), ('Team', 'Middleware', v_ins_pk), ('Environment', 'production', v_ins_pk);
    cost_base := 3.00 + random() * 2.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.50, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 15: GCP staging, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('gke-staging-us', 'gke-staging-us-qrs90', 'qrs90', v_provider, 'Running', 'us-central1', v_acc_id, 'https://console.cloud.google.com/kubernetes/clusters/gke-staging-us', now() - INTERVAL '1 hour', now() - INTERVAL '80 days', 80, 'klee@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..4 LOOP
    v_ins_name := format('gke-%s-node-%s', 'staging-us', j);
    v_ins_type := gcp_types[1 + floor(random()*array_length(gcp_types,1))::INT];
    v_az := 'us-central1-' || chr(97 + (j % 3));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '80 days', 80)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'klee@redhat.com', v_ins_pk), ('Partner', 'TCS', v_ins_pk), ('Team', 'Platform', v_ins_pk), ('Environment', 'staging', v_ins_pk);
    cost_base := 2.00 + random() * 1.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.30, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 16: GCP dev, stopped
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('gke-dev-asia', 'gke-dev-asia-tuv12', 'tuv12', v_provider, 'Stopped', 'asia-east1', v_acc_id, 'https://console.cloud.google.com/kubernetes/clusters/gke-dev-asia', now() - INTERVAL '8 hours', now() - INTERVAL '20 days', 20, 'agarcia@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..3 LOOP
    v_ins_name := format('gke-%s-node-%s', 'dev-asia', j);
    v_ins_type := gcp_types[1 + floor(random()*array_length(gcp_types,1))::INT];
    v_az := 'asia-east1-' || chr(97 + (j % 2));
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Stopped', now() - INTERVAL '8 hours', now() - INTERVAL '20 days', 20)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'agarcia@redhat.com', v_ins_pk), ('Partner', 'Infosys', v_ins_pk), ('Team', 'DevOps', v_ins_pk), ('Environment', 'development', v_ins_pk);
    cost_base := 1.00 + random() * 0.5;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.10, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- ========================================================================
  -- Account 5: rh-services-azure (Azure) — 3 clusters
  -- ========================================================================
  SELECT id INTO v_acc_id FROM accounts WHERE account_id = 'azure-sub-01';
  v_provider := 'Azure';

  -- Cluster 17: Azure production, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('aro-prod-westeu', 'aro-prod-westeu-wxy34', 'wxy34', v_provider, 'Running', 'westeurope', v_acc_id, 'https://portal.azure.com/aro-prod-westeu', now() - INTERVAL '1 hour', now() - INTERVAL '110 days', 110, 'twilson@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..5 LOOP
    v_ins_name := format('aro-%s-node-%s', 'prod-westeu', j);
    v_ins_type := az_types[1 + floor(random()*array_length(az_types,1))::INT];
    v_az := 'westeurope-' || (j % 3 + 1);
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '1 hour', now() - INTERVAL '110 days', 110)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'twilson@redhat.com', v_ins_pk), ('Partner', 'Accenture', v_ins_pk), ('Team', 'Security', v_ins_pk), ('Environment', 'production', v_ins_pk);
    cost_base := 3.80 + random() * 2.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.60, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 18: Azure staging, running
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('aro-staging-north', 'aro-staging-north-zab56', 'zab56', v_provider, 'Running', 'northeurope', v_acc_id, 'https://portal.azure.com/aro-staging-north', now() - INTERVAL '2 hours', now() - INTERVAL '40 days', 40, 'mchen@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..3 LOOP
    v_ins_name := format('aro-%s-node-%s', 'staging-north', j);
    v_ins_type := az_types[1 + floor(random()*array_length(az_types,1))::INT];
    v_az := 'northeurope-' || (j % 3 + 1);
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Running', now() - INTERVAL '2 hours', now() - INTERVAL '40 days', 40)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'mchen@redhat.com', v_ins_pk), ('Partner', 'IBM', v_ins_pk), ('Team', 'Middleware', v_ins_pk), ('Environment', 'staging', v_ins_pk);
    cost_base := 2.50 + random() * 1.0;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.30, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

  -- Cluster 19: Azure APAC dev, stopped
  INSERT INTO clusters (cluster_name, cluster_id, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner)
  VALUES ('aro-dev-seasia', 'aro-dev-seasia-cde78', 'cde78', v_provider, 'Stopped', 'southeastasia', v_acc_id, 'https://portal.azure.com/aro-dev-seasia', now() - INTERVAL '10 hours', now() - INTERVAL '10 days', 10, 'jsmith@redhat.com')
  RETURNING id INTO v_clu_pk;
  FOR j IN 1..3 LOOP
    v_ins_name := format('aro-%s-node-%s', 'dev-seasia', j);
    v_ins_type := az_types[1 + floor(random()*array_length(az_types,1))::INT];
    v_az := 'southeastasia-' || (j % 2 + 1);
    INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age)
    VALUES (v_ins_name, v_ins_name, v_clu_pk, v_provider, v_ins_type, v_az, 'Stopped', now() - INTERVAL '10 hours', now() - INTERVAL '10 days', 10)
    RETURNING id INTO v_ins_pk;
    INSERT INTO tags(key, value, instance_id) VALUES ('Name', v_ins_name, v_ins_pk), ('Owner', 'jsmith@redhat.com', v_ins_pk), ('Partner', 'Deloitte', v_ins_pk), ('Team', 'SRE', v_ins_pk), ('Environment', 'development', v_ins_pk);
    cost_base := 1.20 + random() * 0.6;
    FOR day_offset IN 0..59 LOOP d := current_date - day_offset; amt := round(GREATEST(0.10, cost_base * (0.85 + random()*0.30))::NUMERIC, 2); INSERT INTO expenses(instance_id, date, amount) VALUES (v_ins_pk, d, amt); END LOOP;
  END LOOP;

END
$$;

-- ============================================================================
-- Scheduled actions (with targets)
-- ============================================================================
DO $$
DECLARE
  v_target_id BIGINT;
  v_cluster_id BIGINT;
  v_operation ACTION_OPERATION;
BEGIN
  -- 4 scheduled actions on random clusters
  FOR g IN 1..4 LOOP
    SELECT id INTO v_cluster_id FROM clusters WHERE status != 'Terminated' ORDER BY random() LIMIT 1;
    v_operation := (ARRAY['PowerOn','PowerOff'])[1 + (random()*1)::int]::ACTION_OPERATION;

    INSERT INTO targets (target_type, select_all) VALUES ('Cluster', false) RETURNING id INTO v_target_id;
    INSERT INTO target_clusters (target_id, cluster_id) VALUES (v_target_id, v_cluster_id);
    INSERT INTO schedule (type, time, cron_exp, operation, target, status, enabled, requester, description)
    VALUES ('scheduled_action', now() + (g * interval '1 day'), NULL, v_operation, v_target_id, 'Pending', true, 'scheduler@clusteriq', format('Scheduled %s for cluster', v_operation));
  END LOOP;
END
$$;

-- ============================================================================
-- Cron-based actions (with targets)
-- ============================================================================
DO $$
DECLARE
  v_target_id BIGINT;
  v_cluster_id BIGINT;
  v_operation ACTION_OPERATION;
  v_cron TEXT;
BEGIN
  FOR g IN 1..3 LOOP
    SELECT id INTO v_cluster_id FROM clusters WHERE status != 'Terminated' ORDER BY random() LIMIT 1;
    v_operation := (ARRAY['PowerOn','PowerOff'])[1 + (random()*1)::int]::ACTION_OPERATION;
    v_cron := (ARRAY['0 8 * * 1-5', '0 20 * * 1-5', '0 6 * * *'])[g];

    INSERT INTO targets (target_type, select_all) VALUES ('Cluster', false) RETURNING id INTO v_target_id;
    INSERT INTO target_clusters (target_id, cluster_id) VALUES (v_target_id, v_cluster_id);
    INSERT INTO schedule (type, time, cron_exp, operation, target, status, enabled, requester, description)
    VALUES ('cron_action', NULL, v_cron, v_operation, v_target_id, 'Pending', true, 'scheduler@clusteriq', format('Recurring %s (%s)', v_operation, v_cron));
  END LOOP;
END
$$;

-- ============================================================================
-- Events: realistic audit trail
-- ============================================================================
DO $$
DECLARE
  v_cluster RECORD;
  v_account RECORD;
BEGIN
  -- Scan events (one per account, recent)
  FOR v_account IN SELECT id, account_name FROM accounts LOOP
    INSERT INTO events (event_timestamp, requester, action, resource_id, resource_type, result, description, severity)
    VALUES (now() - (random() * interval '2 hours'), 'scanner@clusteriq', 'Scan', v_account.id, 'Account', 'Success', format('Inventory scan completed for %s', v_account.account_name), 'info');
  END LOOP;

  -- PowerOn/PowerOff events on clusters
  FOR v_cluster IN SELECT id, cluster_name, status FROM clusters WHERE status != 'Terminated' ORDER BY random() LIMIT 8 LOOP
    INSERT INTO events (event_timestamp, requester, action, resource_id, resource_type, result, description, severity)
    VALUES (
      now() - (random() * interval '5 days'),
      (ARRAY['jsmith@redhat.com', 'agarcia@redhat.com', 'scheduler@clusteriq', 'agent@clusteriq'])[1 + floor(random()*4)::INT],
      CASE WHEN v_cluster.status = 'Running' THEN 'PowerOn' ELSE 'PowerOff' END,
      v_cluster.id,
      'Cluster',
      'Success',
      format('%s cluster %s', CASE WHEN v_cluster.status = 'Running' THEN 'Started' ELSE 'Stopped' END, v_cluster.cluster_name),
      'info'
    );
  END LOOP;

  -- A few failed events
  INSERT INTO events (event_timestamp, requester, action, resource_id, resource_type, result, description, severity)
  SELECT
    now() - (random() * interval '7 days'),
    'agent@clusteriq',
    'PowerOff',
    (SELECT id FROM clusters WHERE status = 'Running' ORDER BY random() LIMIT 1),
    'Cluster',
    'Failed',
    'Timeout waiting for instances to stop',
    'error'
  FROM generate_series(1, 2);

  -- Warning events
  INSERT INTO events (event_timestamp, requester, action, resource_id, resource_type, result, description, severity)
  VALUES
    (now() - interval '1 day', 'scanner@clusteriq', 'Scan', (SELECT id FROM accounts ORDER BY random() LIMIT 1), 'Account', 'Success', 'Scan completed with warnings: 2 instances unreachable', 'warning'),
    (now() - interval '3 days', 'scheduler@clusteriq', 'PowerOn', (SELECT id FROM clusters WHERE status = 'Stopped' ORDER BY random() LIMIT 1), 'Cluster', 'Success', 'Scheduled power-on executed', 'info');
END
$$;

-- ============================================================================
-- Refresh materialized views
-- ============================================================================
REFRESH MATERIALIZED VIEW m_accounts_full_view;
REFRESH MATERIALIZED VIEW m_clusters_full_view;
REFRESH MATERIALIZED VIEW m_instances_full_view;
REFRESH MATERIALIZED VIEW m_instances_full_view_with_tags;

COMMIT;
