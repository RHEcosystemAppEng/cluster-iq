-- Cleaning
DELETE FROM tags;
DELETE FROM expenses;
DELETE FROM instances;
DELETE FROM clusters;
DELETE FROM accounts;
DELETE FROM providers;
DELETE FROM status;

-- Providers
INSERT INTO
  providers(name)
VALUES
  ('AWS'),
  ('GCP'),
  ('Azure'),
  ('UNKNOWN')
;


-- Status
INSERT INTO
  status(value)
VALUES
  ('Running'),
  ('Stopped'),
  ('Terminated'),
  ('Unknown')
;


-- Accounts
INSERT INTO
  accounts (id, name, provider, cluster_count, last_scan_timestamp)
VALUES
  ('ABC123', 'engineering', 'AWS', 2, TO_DATE('20/10/2021', 'DD/MM/YYYY')),
  ('XYZ098', 'partners', 'GCP', 2, TO_DATE('20/10/2021', 'DD/MM/YYYY')),
  ('FGH456', 'business', 'Azure', 2, TO_DATE('20/10/2021', 'DD/MM/YYYY'))
;


-- Clusters
INSERT INTO
  clusters (id, name, infra_id, provider, status, region, account_name, console_link, instance_count, last_scan_timestamp, creation_timestamp, age, owner, total_cost)
VALUES
  ('cluster-A-A01-engineering', 'cluster-A', 'A01', 'AWS',   'Running',    'eu-west-1',  'engineering', 'http://console.cluster-A', 6, TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 'John Doe 1', 0.0),
  ('cluster-B-B02-engineering', 'cluster-B', 'B02', 'AWS',   'Stopped',    'eu-east-1',  'engineering', 'http://console.cluster-B', 6, TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 'John Doe 2', 0.0),
  ('cluster-C-C03-partners',    'cluster-C', 'C03', 'GCP',   'Running',    'eu-north-1', 'partners',    'http://console.cluster-C', 6, TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 'John Doe 3', 0.0),
  ('cluster-D-D04-partners',    'cluster-D', 'D04', 'GCP',   'Unknown',    'eu-south-1', 'partners',    'http://console.cluster-D', 6, TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 'John Doe 4', 0.0),
  ('cluster-E-E05-business',    'cluster-E', 'E05', 'Azure', 'Terminated', 'us-west-1',  'business',    'http://console.cluster-E', 6, TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 'John Doe 5', 0.0),
  ('cluster-F-F06-business',    'cluster-F', 'F06', 'Azure', 'Stopped',    'us-east-1',  'business',    'http://console.cluster-F', 6, TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 'John Doe 6', 0.0)
;


-- Instances (36)
INSERT INTO
  instances (id, name, provider, instance_type, availability_zone, status, cluster_id, last_scan_timestamp, creation_timestamp, age, daily_cost, total_cost)
VALUES
  ('id-001', 'instance-01-master', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-A-A01-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-002', 'instance-02-master', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-A-A01-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-003', 'instance-03-master', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-A-A01-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-004', 'instance-04-worker', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-A-A01-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-005', 'instance-05-worker', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-A-A01-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-006', 'instance-06-worker', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-A-A01-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-007', 'instance-01-master', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-B-B02-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-008', 'instance-02-master', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-B-B02-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-009', 'instance-03-master', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-B-B02-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-010', 'instance-04-worker', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-B-B02-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-011', 'instance-05-worker', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-B-B02-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-012', 'instance-06-worker', 'AWS',   't2.medium', 'eu-west-1', 'Running', 'cluster-B-B02-engineering', TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-013', 'instance-01-master', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-C-C03-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-014', 'instance-02-master', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-C-C03-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-015', 'instance-03-master', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-C-C03-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-016', 'instance-04-worker', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-C-C03-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-017', 'instance-05-worker', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-C-C03-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-018', 'instance-06-worker', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-C-C03-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-019', 'instance-01-master', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-D-D04-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-020', 'instance-02-master', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-D-D04-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-021', 'instance-03-master', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-D-D04-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-022', 'instance-04-worker', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-D-D04-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-023', 'instance-05-worker', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-D-D04-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-024', 'instance-06-worker', 'GCP',   't2.medium', 'eu-west-1', 'Running', 'cluster-D-D04-partners',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-025', 'instance-01-master', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-E-E05-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-026', 'instance-02-master', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-E-E05-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-027', 'instance-03-master', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-E-E05-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-028', 'instance-04-worker', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-E-E05-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-029', 'instance-05-worker', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-E-E05-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-030', 'instance-06-worker', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-E-E05-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-031', 'instance-01-master', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-F-F06-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-032', 'instance-02-master', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-F-F06-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-033', 'instance-03-master', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-F-F06-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-034', 'instance-04-worker', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-F-F06-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-035', 'instance-05-worker', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-F-F06-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0),
  ('id-036', 'instance-06-worker', 'Azure', 't2.medium', 'eu-west-1', 'Running', 'cluster-F-F06-business',    TO_DATE('24/06/2024', 'DD/MM/YYYY'), TO_DATE('20/06/2024', 'DD/MM/YYYY'), 4, 0.0, 0.0)
;


-- Expenses (72)
INSERT INTO
  expenses (instance_id, date, amount)
VALUES
  ('id-001', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-001', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-001', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-001', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-002', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-002', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-002', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-002', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-003', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-003', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-003', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-003', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-004', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-004', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-004', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-004', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-005', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-005', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-005', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-005', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-006', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-006', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-006', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-006', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-007', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-007', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-007', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-007', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-008', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-008', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-008', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-008', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-009', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-009', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-009', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-009', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-010', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-010', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-010', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-010', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-011', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-011', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-011', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-011', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-012', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-012', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-012', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-012', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-013', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-013', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-013', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-013', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-014', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-014', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-014', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-014', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-015', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-015', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-015', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-015', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-016', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-016', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-016', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-016', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-017', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-017', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-017', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-017', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-018', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-018', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-018', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-018', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-019', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-019', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-019', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-019', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-020', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-020', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-020', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-020', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-021', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-021', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-021', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-021', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-022', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-022', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-022', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-022', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-023', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-023', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-023', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-023', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-024', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-024', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-024', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-024', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-025', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-025', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-025', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-025', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-026', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-026', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-026', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-026', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-027', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-027', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-027', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-027', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-028', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-028', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-028', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-028', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-029', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-029', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-029', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-029', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-030', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-030', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-030', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-030', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-031', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-031', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-031', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-031', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-032', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-032', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-032', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-032', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-033', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-033', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-033', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-033', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-034', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-034', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-034', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-034', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-035', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-035', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-035', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-035', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12),
  ('id-036', TO_DATE('20/06/2024', 'DD/MM/YYYY'), 8.25),
  ('id-036', TO_DATE('21/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-036', TO_DATE('22/06/2024', 'DD/MM/YYYY'), 12.35),
  ('id-036', TO_DATE('24/06/2024', 'DD/MM/YYYY'), 10.12)
;


INSERT INTO
  tags(key, value, instance_id)
values
  ('key01', 'value', 'id-001'),
  ('key02', 'value', 'id-001'),
  ('key01', 'value', 'id-002'),
  ('key02', 'value', 'id-002'),
  ('key01', 'value', 'id-003'),
  ('key02', 'value', 'id-003'),
  ('key01', 'value', 'id-004'),
  ('key02', 'value', 'id-004'),
  ('key01', 'value', 'id-005'),
  ('key02', 'value', 'id-005'),
  ('key01', 'value', 'id-006'),
  ('key02', 'value', 'id-006'),
  ('key01', 'value', 'id-007'),
  ('key02', 'value', 'id-007'),
  ('key01', 'value', 'id-008'),
  ('key02', 'value', 'id-008'),
  ('key01', 'value', 'id-009'),
  ('key02', 'value', 'id-009'),
  ('key01', 'value', 'id-010'),
  ('key02', 'value', 'id-010'),
  ('key01', 'value', 'id-011'),
  ('key02', 'value', 'id-011'),
  ('key01', 'value', 'id-012'),
  ('key02', 'value', 'id-012'),
  ('key01', 'value', 'id-013'),
  ('key02', 'value', 'id-013'),
  ('key01', 'value', 'id-014'),
  ('key02', 'value', 'id-014'),
  ('key01', 'value', 'id-015'),
  ('key02', 'value', 'id-015'),
  ('key01', 'value', 'id-016'),
  ('key02', 'value', 'id-016'),
  ('key01', 'value', 'id-017'),
  ('key02', 'value', 'id-017'),
  ('key01', 'value', 'id-018'),
  ('key02', 'value', 'id-018'),
  ('key01', 'value', 'id-019'),
  ('key02', 'value', 'id-019'),
  ('key01', 'value', 'id-020'),
  ('key02', 'value', 'id-020'),
  ('key01', 'value', 'id-021'),
  ('key02', 'value', 'id-021'),
  ('key01', 'value', 'id-022'),
  ('key02', 'value', 'id-022'),
  ('key01', 'value', 'id-023'),
  ('key02', 'value', 'id-023'),
  ('key01', 'value', 'id-024'),
  ('key02', 'value', 'id-024'),
  ('key01', 'value', 'id-025'),
  ('key02', 'value', 'id-025'),
  ('key01', 'value', 'id-026'),
  ('key02', 'value', 'id-026'),
  ('key01', 'value', 'id-027'),
  ('key02', 'value', 'id-027'),
  ('key01', 'value', 'id-028'),
  ('key02', 'value', 'id-028'),
  ('key01', 'value', 'id-029'),
  ('key02', 'value', 'id-029'),
  ('key01', 'value', 'id-030'),
  ('key02', 'value', 'id-030'),
  ('key01', 'value', 'id-031'),
  ('key02', 'value', 'id-031'),
  ('key01', 'value', 'id-032'),
  ('key02', 'value', 'id-032'),
  ('key01', 'value', 'id-033'),
  ('key02', 'value', 'id-033'),
  ('key01', 'value', 'id-034'),
  ('key02', 'value', 'id-034'),
  ('key01', 'value', 'id-035'),
  ('key02', 'value', 'id-035'),
  ('key01', 'value', 'id-036'),
  ('key02', 'value', 'id-036')
;


SELECT count(*) FROM accounts;
SELECT count(*) FROM clusters;
SELECT count(*) FROM instances;
SELECT count(*) FROM expenses;
