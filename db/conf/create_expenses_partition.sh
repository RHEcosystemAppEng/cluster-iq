#!/bin/bash
PGUSER=postgres
PGDATABASE=clusteriq
PGHOST=localhost

psql -U "$PGUSER" -d "$PGDATABASE" -h "$PGHOST" -f /opt/scripts/create_expenses_partition.sql

