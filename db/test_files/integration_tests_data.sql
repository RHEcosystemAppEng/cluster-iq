BEGIN;

-- Cleaning
TRUNCATE expenses, tags, instances, clusters, accounts RESTART IDENTITY CASCADE;

-- ## Accounts ##
INSERT INTO accounts (account_id, account_name, provider, last_scan_ts, created_at) VALUES
  ('111111111111', 'aws-account-demo',   'AWS',   '2025-08-01 10:00:00+00', '2025-01-01 00:00:00+00'),
  ('gcp-project-1', 'gcp-project-demo',  'GCP',   '2025-08-02 10:00:00+00', '2025-01-01 00:00:00+00'),
  ('subs-00000001', 'azure-sub-demo',    'Azure', '2025-08-03 10:00:00+00', '2025-01-01 00:00:00+00');

-- ## Clusters ##
INSERT INTO clusters (cluster_id, cluster_name, infra_id, provider, status, region, account_id, console_link, last_scan_ts, created_at, age, owner) VALUES
  ('aws-cluster-1-aws-infra-1',   'aws-cluster-1',   'aws-infra-1',   'AWS',   'Running',   'us-east-1', 1, 'https://console.aws/1', '2025-08-01 12:00:00+00', '2025-02-01 00:00:00+00', 180, 'team-aws'),
  ('aws-cluster-2-aws-infra-2',   'aws-cluster-2',   'aws-infra-2',   'AWS',   'Stopped',   'us-east-2', 1, 'https://console.aws/2', '2025-08-01 12:00:00+00', '2025-03-01 00:00:00+00', 150, 'team-aws'),
  ('gcp-cluster-1-gcp-infra-1',   'gcp-cluster-1',   'gcp-infra-1',   'GCP',   'Running',   'europe-west1', 2, 'https://console.gcp/1', '2025-08-02 12:00:00+00', '2025-02-01 00:00:00+00', 160, 'team-gcp'),
  ('gcp-cluster-2-gcp-infra-2',   'gcp-cluster-2',   'gcp-infra-2',   'GCP',   'Unknown',   'europe-west2', 2, 'https://console.gcp/2', '2025-08-02 12:00:00+00', '2025-03-01 00:00:00+00', 140, 'team-gcp'),
  ('azure-cluster-1-az-infra-1',  'azure-cluster-1', 'az-infra-1',    'Azure', 'Running',   'westeurope', 3, 'https://portal.azure/1', '2025-08-03 12:00:00+00', '2025-02-01 00:00:00+00', 170, 'team-az'),
  ('azure-cluster-2-az-infra-2',  'azure-cluster-2', 'az-infra-2',    'Azure', 'Stopped',   'westeurope', 3, 'https://portal.azure/2', '2025-08-03 12:00:00+00', '2025-03-01 00:00:00+00', 130, 'team-az');

-- ## Instances ##
INSERT INTO instances (instance_id, instance_name, cluster_id, provider, instance_type, availability_zone, status, last_scan_ts, created_at, age) VALUES
  ('id-0123456789X', 'aws-instance-1a', 1, 'AWS',   't3.micro',   'us-east-1a', 'Running', '2025-08-01 12:00:00+00', '2025-02-10 00:00:00+00', 170),
  ('id-1123456789X', 'aws-instance-1b', 1, 'AWS',   't3.medium',  'us-east-1b', 'Stopped', '2025-08-01 12:00:00+00', '2025-02-11 00:00:00+00', 169),
  ('id-2123456789X', 'aws-instance-2a', 2, 'AWS',   'm6g.large',  'us-east-2a', 'Running', '2025-08-01 12:00:00+00', '2025-03-15 00:00:00+00', 140),
  ('id-3123456789X', 'aws-instance-2b', 2, 'AWS',   'c6i.large',  'us-east-2b', 'Unknown', '2025-08-01 12:00:00+00', '2025-03-20 00:00:00+00', 135),

  ('id-0123456789Y', 'gcp-instance-1a', 3, 'GCP',   'e2-small',   'europe-west1-b', 'Running', '2025-08-02 12:00:00+00', '2025-02-05 00:00:00+00', 175),
  ('id-1123456789Y', 'gcp-instance-1b', 3, 'GCP',   'e2-medium',  'europe-west1-c', 'Stopped', '2025-08-02 12:00:00+00', '2025-02-06 00:00:00+00', 174),
  ('id-2123456789Y', 'gcp-instance-2a', 4, 'GCP',   'n2-standard', 'europe-west2-a','Running', '2025-08-02 12:00:00+00', '2025-03-10 00:00:00+00', 150),
  ('id-3123456789Y', 'gcp-instance-2b', 4, 'GCP',   'n2-standard', 'europe-west2-b','Unknown', '2025-08-02 12:00:00+00', '2025-03-11 00:00:00+00', 149),

  ('id-0123456789Z', 'az-instance-1a',  5, 'Azure', 'B1s',        'westeurope-1',   'Running', '2025-08-03 12:00:00+00', '2025-02-01 00:00:00+00', 180),
  ('id-1123456789Z', 'az-instance-1b',  5, 'Azure', 'B2s',        'westeurope-2',   'Stopped', '2025-08-03 12:00:00+00', '2025-02-02 00:00:00+00', 179),
  ('id-2123456789Z', 'az-instance-2a',  6, 'Azure', 'D2s_v3',     'westeurope-1',   'Running', '2025-08-03 12:00:00+00', '2025-03-05 00:00:00+00', 160),
  ('id-3123456789Z', 'az-instance-2b',  6, 'Azure', 'D2s_v3',     'westeurope-2',   'Unknown', '2025-08-03 12:00:00+00', '2025-03-06 00:00:00+00', 159);

-- ## Instance Tags ##
INSERT INTO tags (key, value, instance_id) VALUES
  ('name','aws-instance-1a',1), ('owner','john.doe',1),
  ('name','aws-instance-1b',2), ('owner','john.doe',2),
  ('name','aws-instance-2a',3), ('owner','john.doe',3),
  ('name','aws-instance-2b',4), ('owner','john.doe',4),
  ('name','gcp-instance-1a',5), ('owner','jane.doe',5),
  ('name','gcp-instance-1b',6), ('owner','jane.doe',6),
  ('name','gcp-instance-2a',7), ('owner','jane.doe',7),
  ('name','gcp-instance-2b',8), ('owner','jane.doe',8),
  ('name','az-instance-1a',9), ('owner','alice',9),
  ('name','az-instance-1b',10),('owner','alice',10),
  ('name','az-instance-2a',11),('owner','bob',11),
  ('name','az-instance-2b',12),('owner','bob',12);

-- ## Expenses ##
-- 5 days expenses per instance
CREATE TABLE expenses_2025_07 PARTITION OF expenses FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');
CREATE TABLE expenses_2025_08 PARTITION OF expenses FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');

INSERT INTO expenses (instance_id, date, amount) VALUES
  (1,'2025-07-30',1.00), (1,'2025-07-31',1.10), (1,'2025-08-01',1.20), (1,'2025-08-02',1.30), (1,'2025-08-03',1.40),
  (2,'2025-07-30',0.90), (2,'2025-07-31',0.95), (2,'2025-08-01',1.00), (2,'2025-08-02',1.05), (2,'2025-08-03',1.10),
  (3,'2025-07-30',2.00), (3,'2025-07-31',2.10), (3,'2025-08-01',2.20), (3,'2025-08-02',2.30), (3,'2025-08-03',2.40),
  (4,'2025-07-30',1.50), (4,'2025-07-31',1.55), (4,'2025-08-01',1.60), (4,'2025-08-02',1.65), (4,'2025-08-03',1.70),
  (5,'2025-07-30',1.20), (5,'2025-07-31',1.25), (5,'2025-08-01',1.30), (5,'2025-08-02',1.35), (5,'2025-08-03',1.40),
  (6,'2025-07-30',1.00), (6,'2025-07-31',1.05), (6,'2025-08-01',1.10), (6,'2025-08-02',1.15), (6,'2025-08-03',1.20),
  (7,'2025-07-30',2.50), (7,'2025-07-31',2.55), (7,'2025-08-01',2.60), (7,'2025-08-02',2.65), (7,'2025-08-03',2.70),
  (8,'2025-07-30',2.10), (8,'2025-07-31',2.15), (8,'2025-08-01',2.20), (8,'2025-08-02',2.25), (8,'2025-08-03',2.30),
  (9,'2025-07-30',0.80), (9,'2025-07-31',0.85), (9,'2025-08-01',0.90), (9,'2025-08-02',0.95), (9,'2025-08-03',1.00),
  (10,'2025-07-30',1.40),(10,'2025-07-31',1.45),(10,'2025-08-01',1.50),(10,'2025-08-02',1.55),(10,'2025-08-03',1.60),
  (11,'2025-07-30',1.70),(11,'2025-07-31',1.75),(11,'2025-08-01',1.80),(11,'2025-08-02',1.85),(11,'2025-08-03',1.90),
  (12,'2025-07-30',1.10),(12,'2025-07-31',1.15),(12,'2025-08-01',1.20),(12,'2025-08-02',1.25),(12,'2025-08-03',1.30);


INSERT INTO events ( event_timestamp, triggered_by, action, resource_id, resource_type, result, description, severity) VALUES
  ('2025-08-02 12:00:00+00', 'cluster-iq-tester', 'test', '1', 'cluster', 'OK', 'integration test event', 'info'),
  ('2025-08-02 12:00:00+00', 'cluster-iq-tester', 'test', '10', 'instance', 'Pending', 'integration test event', 'critical');


COMMIT;
