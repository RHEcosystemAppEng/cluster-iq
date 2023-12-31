-- Providers
CREATE TABLE IF NOT EXISTS providers (
  name TEXT PRIMARY KEY
);

INSERT INTO
  providers(name)
VALUES
  ('AWS'),
  ('GCP'),
  ('Azure'),
  ('UNKNOWN')
;


-- Accounts
CREATE TABLE IF NOT EXISTS accounts (
  name TEXT PRIMARY KEY,
  provider TEXT REFERENCES providers(name),
  cluster_count INTEGER
);


-- Clusters
CREATE TABLE IF NOT EXISTS clusters (
  name TEXT PRIMARY KEY,
  provider TEXT REFERENCES providers(name),
  state TEXT,
  region TEXT,
  account_name TEXT REFERENCES accounts(name),
  console_link TEXT,
  instance_count INTEGER
);


-- Instances
CREATE TABLE IF NOT EXISTS instances (
  id TEXT PRIMARY KEY,
  name TEXT,
  provider TEXT REFERENCES providers(name),
  instance_type TEXT,
  region TEXT,
  state TEXT,
  cluster_name TEXT REFERENCES clusters(name)
  -- TODO: ADD tags
);


CREATE TABLE IF NOT EXISTS tags (
  key TEXT,
  value TEXT,
  instance_id TEXT REFERENCES instances(id),
  PRIMARY KEY (key, instance_id)
);
