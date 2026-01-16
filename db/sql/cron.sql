\c postgres

-- pg_cron task for updating the 'Terminated' elements in the inventory every 6 hours
SELECT cron.schedule_in_database(
  'check_terminated_inventory',
  '0 */6 * * *',
  $$SELECT check_terminated_inventory();$$,
  'clusteriq'
);

-- pg_cron task for creating a new monthly partition for 'Expenses' table every Sunday
SELECT cron.schedule_in_database(
  'expenses_partitioning',
  '0 0 * * 6',
  $$SELECT create_next_month_expenses_partition();$$,
  'clusteriq'
);

-- pg_cron task for creating a new monthly partition for 'audit_logs' table every Sunday
SELECT cron.schedule_in_database(
  'events_partitioning',
  '0 0 * * 6',
  $$SELECT create_next_month_events_partition();$$,
  'clusteriq'
);

-- pg_cron task for refreshing the materialized views if the process is not triggered
SELECT cron.schedule_in_database(
  'm_views_refresh',
  '0 */2 * * *',
  $$SELECT refresh_materialized_views();$$,
  'clusteriq'
);

-- Function to check easier how the pg_cron tasks went
CREATE OR REPLACE FUNCTION pg_cron_history(p_limit int DEFAULT 20)
RETURNS TABLE(
  jobid          INT,
  jobname        TEXT,
  schedule       TEXT,
  command        TEXT,
  status         TEXT,
  start_time     TIMESTAMPTZ,
  end_time       TIMESTAMPTZ,
  return_message TEXT,
  active         BOOLEAN
)
AS $$
  SELECT
    j.jobid,
    j.jobname,
    j.schedule,
    j.command,
    jrd.status,
    jrd.start_time,
    jrd.end_time,
    jrd.return_message,
    j.active
  FROM cron.job j
  JOIN cron.job_run_details jrd ON j.jobid = jrd.jobid
  ORDER BY jrd.start_time DESC
  LIMIT p_limit;
$$ LANGUAGE sql STABLE;
