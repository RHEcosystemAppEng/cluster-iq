DO $$
DECLARE
    next_month DATE := date_trunc('month', current_date) + INTERVAL '1 month';
    partition_name TEXT := format('expenses_%s', to_char(next_month, 'YYYY_MM'));
    start_date DATE := next_month;
    end_date DATE := next_month + INTERVAL '1 month';
BEGIN
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF expenses
         FOR VALUES FROM (%L) TO (%L);',
        partition_name, start_date, end_date
    );
END$$;

