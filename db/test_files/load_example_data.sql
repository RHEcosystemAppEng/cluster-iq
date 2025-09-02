-- Demo seed for ClusterIQ schema
-- psql postgresql://user:password@pgsql:5432/clusteriq < load_example_data.sql
BEGIN;

-- Limpia datos previos (si los hubiera)
TRUNCATE expenses, tags, instances, clusters, accounts RESTART IDENTITY CASCADE;
-- Inserta 3 cuentas (una por proveedor) y guarda sus IDs
WITH ins AS (
  INSERT INTO accounts (account_id, account_name, provider, last_scan_ts)
  VALUES
    ('111111111111', 'aws-account-demo',   'AWS',    now() - INTERVAL '1 day'),
    ('gcp-project-1', 'gcp-project-demo',  'GCP',    now() - INTERVAL '2 days'),
    ('subs-00000001', 'azure-sub-demo',    'Azure',  now() - INTERVAL '3 days')
  RETURNING id, provider
)
SELECT * FROM ins;

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


DO $$
DECLARE
  cur_start  DATE := date_trunc('month', current_date)::date;
  cur_end    DATE := (cur_start + INTERVAL '1 month')::date;

  prev_start DATE := date_trunc('month', current_date - INTERVAL '1 month')::date;
  prev_end   DATE := cur_start;

  prev_prev_start DATE := date_trunc('month', current_date - INTERVAL '2 month')::date;
  prev_prev_end   DATE := prev_start;

  cur_suffix  TEXT := to_char(cur_start,  'YYYY_MM');
  prev_suffix TEXT := to_char(prev_start, 'YYYY_MM');
  prev_prev_suffix TEXT := to_char(prev_prev_start, 'YYYY_MM');

  part_name   TEXT;
  sql         TEXT;
BEGIN
  -- Partición del mes anterior anterior: expenses_YYYY_MM
  part_name := format('expenses_%s', prev_prev_suffix);
  IF to_regclass(part_name) IS NULL THEN
    sql := format(
      'CREATE TABLE %I PARTITION OF expenses
       FOR VALUES FROM (%L) TO (%L);',
      part_name, prev_prev_start, prev_prev_end
    );
    EXECUTE sql;
  END IF;

  -- Partición del mes anterior: expenses_YYYY_MM
  part_name := format('expenses_%s', prev_suffix);
  IF to_regclass(part_name) IS NULL THEN
    sql := format(
      'CREATE TABLE %I PARTITION OF expenses
       FOR VALUES FROM (%L) TO (%L);',
      part_name, prev_start, prev_end
    );
    EXECUTE sql;
  END IF;

  -- Partición del mes en curso: expenses_YYYY_MM
  part_name := format('expenses_%s', cur_suffix);
  IF to_regclass(part_name) IS NULL THEN
    sql := format(
      'CREATE TABLE %I PARTITION OF expenses
       FOR VALUES FROM (%L) TO (%L);',
      part_name, cur_start, cur_end
    );
    EXECUTE sql;
  END IF;
END
$$;



-- Función auxiliar para obtener un STATUS aleatorio
-- (usamos un SELECT al vuelo dentro del DO; no se define nada permanente)
-- Generación principal
DO $$
DECLARE
  -- cuentas
  r_acc RECORD;

  -- clusters
  n_clusters  INT;
  r_clu_id    BIGINT;
  r_clu_name  TEXT;
  r_clu_infra  TEXT;
  r_region    TEXT;

  -- instancias
  n_insts     INT;
  r_ins_id    BIGINT;
  r_ins_name  TEXT;
  r_az        TEXT;
  st          STATUS;

  -- fechas/importe expenses
  d DATE;
  amt NUMERIC(12,2);
  base_cost NUMERIC(12,2);
  day_offset INT;
BEGIN
  -- Itera por cada cuenta creada
  FOR r_acc IN
    SELECT id, provider FROM accounts ORDER BY id
  LOOP
    -- nº de clusters aleatorio por cuenta: 4–10
    n_clusters := 4 + floor(random() * 7)::INT;

    FOR i IN 1..n_clusters LOOP
      -- Región “plausible” por proveedor
      IF r_acc.provider = 'AWS'::cloud_provider THEN
        r_region := ('us-east-' || (1 + floor(random()*2))::INT);
      ELSIF r_acc.provider = 'GCP'::cloud_provider THEN
        r_region := ('europe-west' || (1 + floor(random()*2))::INT);
      ELSE
        r_region := ('westeurope');
      END IF;

      r_clu_name := format('%s-cluster-%s', lower(r_acc.provider::TEXT), i);
      r_clu_infra := format('%s-infra-%s', lower(r_acc.provider::TEXT), i);

      INSERT INTO clusters (
        cluster_name, cluster_id, infra_id, provider, status, region, account_id,
        console_link, last_scan_ts, created_at, age, owner
      )
      VALUES (
        r_clu_name,
        r_clu_name || r_clu_infra,
        r_clu_infra,
        r_acc.provider,
        (ARRAY['Running','Stopped','Unknown']::status[])[1 + floor(random()*3)::INT],
        r_region,
        r_acc.id,
        'https://console.example.local',
        now() - make_interval(days => (1 + floor(random()*5))::INT),
        now() - make_interval(days => (10 + floor(random()*90))::INT),
        10 + floor(random()*90)::INT,
        'team@example.com'
      )
      RETURNING id INTO r_clu_id;

      -- nº de instancias por cluster: 6–12
      n_insts := 6 + floor(random() * 7)::INT;

      FOR j IN 1..n_insts LOOP
        -- Status aleatorio (más prob. de Running)
        st := (ARRAY['Running','Running','Running','Stopped','Unknown']::status[])[1 + floor(random()*5)::INT];

        -- AZ derivada de la región
        r_az := r_region || chr(97 + floor(random()*3)::INT); -- a/b/c

        r_ins_name := format('%s-%s-%s', r_clu_name, r_az, j);

        INSERT INTO instances (
          instance_id, instance_name, cluster_id, provider, instance_type, availability_zone,
          status, last_scan_ts, created_at, age
        )
        VALUES (
          r_ins_name,
          'id-' || r_ins_name,
          r_clu_id,
          r_acc.provider,
          (ARRAY['t3.micro','t3.medium','m6g.large','c6i.large']::TEXT[])[1 + floor(random()*4)::INT],
          r_az,
          st,
          now() - make_interval(days => (0 + floor(random()*3))::INT),
          now() - make_interval(days => (20 + floor(random()*200))::INT),
          20 + floor(random()*200)::INT
        )
        RETURNING id INTO r_ins_id;

        -- Tags fijas por instancia
        INSERT INTO tags(key, value, instance_id)
        VALUES
          ('name',  r_ins_name, r_ins_id),
          ('owner', 'john.doe@example.com', r_ins_id);

        -- Expenses: 60 días hacia atrás (>= 40)
        base_cost := (ARRAY[0.75, 1.25, 1.80, 2.40]::NUMERIC[])[1 + floor(random()*4)::INT];
        FOR day_offset IN 0..59 LOOP
          d := (current_date - day_offset);
          -- Variación diaria suave
          amt := round( GREATEST(0.10, base_cost * (0.8 + random()*0.6))::NUMERIC, 2);
          INSERT INTO expenses(instance_id, date, amount)
          VALUES (r_ins_id, d, amt);
        END LOOP;
      END LOOP;
    END LOOP;
  END LOOP;
END
$$;

COMMIT;

-- Comprobaciones rápidas (opcionales)
-- SELECT provider, COUNT(*) clusters FROM clusters JOIN accounts a ON a.id = clusters.account_id GROUP BY provider;
-- SELECT COUNT(*) AS total_instances FROM instances;
-- SELECT COUNT(*) AS total_expenses FROM expenses;

-- Vistas de métricas
-- SELECT * FROM account_cluster_count ORDER BY account_id;
-- SELECT * FROM account_costs ORDER BY id;
-- SELECT * FROM account_full_view ORDER BY id;

