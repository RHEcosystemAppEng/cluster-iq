#!/bin/bash

# Default values
DB_NAME="clusteriq"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"
BACKUP_DIR="./backups"
RESTORE_LOG_FILE="./ciq-restore.log"

# Parse CLI arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --db-name)
      DB_NAME="$2"
      shift 2
      ;;
    --db-user)
      DB_USER="$2"
      shift 2
      ;;
    --db-passwd)
      DB_PASSWORD="$2"
      shift 2
      ;;
    --db-admin-user)
      DB_ADMIN_USER="$2"
      shift 2
      ;;
    --db-admin-passwd)
      DB_ADMIN_PASSWORD="$2"
      shift 2
      ;;
    --db-host)
      DB_HOST="$2"
      shift 2
      ;;
    --db-port)
      DB_PORT="$2"
      shift 2
      ;;
    --backup)
      BACKUP="$2"
      shift 2
      ;;
    *)
      echo -e "[\033[33m⚠️ \033[0m] Unknown option: $1"
      exit 1
      ;;
  esac
done

# Checking if needed binaries are installed
if ! command -v pg_dump &> /dev/null || ! command -v psql &> /dev/null; then
  echo -e "[\033[31m❌\033[0m] Required PGSQL binaries missing. Check if your system has 'psql' and 'pg_dump' binaries available before continuing"
  exit 1
fi

# Warn if defaults are used
[[ "$DB_NAME" == "clusteriq" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_NAME: clusteriq"; }
[[ "$DB_USER" == "postgres" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_USER: postgres"; }
[[ "$DB_ADMIN_USER" == "postgres" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_ADMIN_USER: postgres"; }
[[ "$DB_HOST" == "localhost" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_HOST: localhost"; }
[[ "$DB_PORT" == "5432" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_PORT: 5432"; }
[[ -z "$BACKUP" ]] && { echo -e "[\033[31m❌\033[0m] No backups have been specified"; exit 1; }

# Checking DB password
if [[ -n "$DB_PASSWORD" ]]; then
  export PGPASSWORD="$DB_PASSWORD"
else
  echo -e "[\033[33m⚠️ \033[0m] No password provided for restore user"
fi

# Check if DB_NAME exists before restoring
echo -e "[\033[31m💢\033[0m] Checking if $DB_NAME exists before restoring. If so, it will be cleaned before restoring"
truncate -s 0 $RESTORE_LOG_FILE
psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1

if [[ $? -eq 0 ]]; then # If DB exists, re-create schema
	# Checking Admin DB password
	if [[ ! -n "$DB_ADMIN_PASSWORD" ]]; then
		echo -e "[\033[33m⚠️ \033[0m] No password provided for superuser"
	fi

  echo -e "[\033[34m⚠️ \033[0m] Database '$DB_NAME' exists, dropping it..."
  PGPASSWORD="$DB_ADMIN_PASSWORD" psql -U "$DB_ADMIN_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" &>> $RESTORE_LOG_FILE

else # If the DB doesn't exist, create it
  echo -e "[\033[34m🚀\033[0m] Creating database '$DB_NAME'..."
  PGPASSWORD="$DB_ADMIN_PASSWORD" psql -U "$DB_ADMIN_USER" -h "$DB_HOST" -p "$DB_PORT" -d postgres -c "CREATE DATABASE \"$DB_NAME\"" &>> $RESTORE_LOG_FILE
fi

# Granting permissions on new schema/DB
echo -e "[\033[34m🔄\033[0m] Granting $DB_USER permissions on $DB_NAME"
PGPASSWORD="$DB_ADMIN_PASSWORD" psql -U "$DB_ADMIN_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -c "GRANT USAGE, CREATE ON SCHEMA public TO \"$DB_USER\";" &>> $RESTORE_LOG_FILE

# Run the restore using psql
echo -e "[\033[34m🔄\033[0m] Restoring database '$DB_NAME' from backup: $BACKUP"
psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -f "$BACKUP" &>> $RESTORE_LOG_FILE
[[ $? -eq 0 ]] && { echo -e "[\033[32m✅\033[0m] Restore completed successfully" ; } || { echo -e "[\033[31m❌\033[0m] Restore failed"; exit 1; }

# Row count summary
echo "[📊] Table row counts after restore:"
TABLES=$(psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -Atc \
"SELECT tablename FROM pg_tables WHERE schemaname = 'public'")

for table in $TABLES; do
  count=$(psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -Atc "SELECT COUNT(*) FROM public.\"$table\";")
  printf "  - %-30s: %s\n" "$table" "$count"
done

echo -e "[\033[32m✅\033[0m] Restore Done"
