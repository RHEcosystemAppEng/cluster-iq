#!/bin/bash

# Default values
DB_NAME="clusteriq"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"
BACKUP_DIR="./backups"

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
    --db-host)
      DB_HOST="$2"
      shift 2
      ;;
    --db-port)
      DB_PORT="$2"
      shift 2
      ;;
    --backup-dir)
      BACKUP_DIR="$2"
      shift 2
      ;;
    *)
      echo -e "[\033[33m⚠️ \033[0m] Unknown option: $1"
      exit 1
      ;;
  esac
done

# Warn if defaults are used
[[ "$DB_NAME" == "clusteriq" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_NAME: clusteriq"; }
[[ "$DB_USER" == "postgres" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_USER: postgres"; }
[[ "$DB_HOST" == "localhost" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_HOST: localhost"; }
[[ "$DB_PORT" == "5432" ]] && { echo -e "[\033[33m⚠️ \033[0m] Using default DB_PORT: 5432"; }

# Checking DB password
if [[ -n "$DB_PASSWORD" ]]; then
  export PGPASSWORD="$DB_PASSWORD"
else
  echo -e "[\033[33m⚠️ \033[0m] No password provided, relying on .pgpass or interactive prompt"
fi

# Generate backup file name
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="${BACKUP_DIR}/${DB_NAME}_backup_${DATE}.sql"

# Ensure backup directory exists
[[ ! -d $BACKUP_DIR ]] && { echo -e "[\033[33m⚠️ \033[0m] backups folder missing. Creating it."; mkdir -p "$BACKUP_DIR"; }

# Dump db
echo -e "[\033[34m🔄\033[0m] Backing up DB data..."
pg_dump -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" \
  --format=plain \
  --encoding=UTF8 \
  --no-owner \
  --no-privileges \
  "$DB_NAME" > "${FILENAME}"
[[ $? -eq 0 ]] && { echo -e "[\033[32m✅\033[0m] Data backup successful"; } || { echo -e "[\033[31m❌\033[0m] Data backup error!"; }
echo -e "[\033[32m✅\033[0m] Backup completed and saved to: $FILENAME"


# Row count summary
echo "[📊] Table row counts after backup:"
TABLES=$(psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -Atc \
"SELECT tablename FROM pg_tables WHERE schemaname = 'public'")

for table in $TABLES; do
  count=$(psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -Atc "SELECT COUNT(*) FROM public.\"$table\";")
  printf "  - %-30s: %s\n" "$table" "$count"
done

echo -e "[\033[32m✅\033[0m] Backup Done"
