DELETE FROM accounts;
DELETE FROM clusters;


INSERT INTO
  accounts (name, provider)
VALUES
  ('engineering', 'AWS'),
  ('partners', 'GCP'),
  ('business', 'Azure')
;

INSERT INTO
  clusters (name, provider, state, region, account_name, console_link)
VALUES
  ('cluster-01', 'AWS', 'Running', 'eu-west-1', 'engineering', 'http://console.cluster-01'),
  ('cluster-02', 'AWS', 'Stopped', 'eu-east-1', 'engineering', 'http://console.cluster-02'),
  ('cluster-03', 'GCP', 'Running', 'eu-north-1', 'partners', 'http://console.cluster-03'),
  ('cluster-04', 'GCP', 'Stopped', 'eu-south-1', 'partners', 'http://console.cluster-04'),
  ('cluster-05', 'Azure', 'Running', 'us-west-1', 'business', 'http://console.cluster-05'), ('cluster-06', 'Azure', 'Stopped', 'us-east-1', 'business', 'http://console.cluster-06')
;

INSERT INTO
  instances (id, name, provider, instance_type, state, region, cluster_name)
VALUES
  ('id-001', 'instance-01-master', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-01'),
  ('id-002', 'instance-02-master', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-01'),
  ('id-003', 'instance-03-master', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-01'),
  ('id-004', 'instance-04-worker', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-01'),
  ('id-005', 'instance-05-worker', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-01'),
  ('id-006', 'instance-06-worker', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-01'),
                                                       
  ('id-007', 'instance-01-master', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-02'),
  ('id-008', 'instance-02-master', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-02'),
  ('id-009', 'instance-03-master', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-02'),
  ('id-010', 'instance-04-worker', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-02'),
  ('id-011', 'instance-05-worker', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-02'),
  ('id-012', 'instance-06-worker', 'AWS', 't2.medium', 'Running', 'eu-west-1', 'cluster-02'),
                                                       
  ('id-013', 'instance-01-master', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-03'),
  ('id-014', 'instance-02-master', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-03'),
  ('id-015', 'instance-03-master', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-03'),
  ('id-016', 'instance-04-worker', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-03'),
  ('id-017', 'instance-05-worker', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-03'),
  ('id-018', 'instance-06-worker', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-03'),
                                                       
  ('id-019', 'instance-01-master', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-04'),
  ('id-020', 'instance-02-master', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-04'),
  ('id-021', 'instance-03-master', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-04'),
  ('id-022', 'instance-04-worker', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-04'),
  ('id-023', 'instance-05-worker', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-04'),
  ('id-024', 'instance-06-worker', 'GCP', 't2.medium', 'Running', 'eu-west-1', 'cluster-04'),

  ('id-025', 'instance-01-master', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-05'),
  ('id-026', 'instance-02-master', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-05'),
  ('id-027', 'instance-03-master', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-05'),
  ('id-028', 'instance-04-worker', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-05'),
  ('id-029', 'instance-05-worker', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-05'),
  ('id-030', 'instance-06-worker', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-05'),
                                                         
  ('id-031', 'instance-01-master', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-06'),
  ('id-032', 'instance-02-master', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-06'),
  ('id-033', 'instance-03-master', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-06'),
  ('id-034', 'instance-04-worker', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-06'),
  ('id-035', 'instance-05-worker', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-06'),
  ('id-036', 'instance-06-worker', 'Azure', 't2.medium', 'Running', 'eu-west-1', 'cluster-06')
;

SELECT count(*) FROM accounts;
SELECT count(*) FROM clusters;
SELECT count(*) FROM instances;
