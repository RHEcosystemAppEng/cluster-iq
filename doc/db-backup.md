# ClusterIQ DB guide
This document describes how to manage different operations on ClusterIQ database.


## DB Backup
To create a database backup follow this steps:
1. Log in on your Openshift environment with the `oc` CLI, and set the
   `NAMESPACE` env var.
   ```sh
   oc login ...

   export NAMESPACE="<YOUR_CLUSTER_IQ_NAMESPACE>"
   ```

2. Stop the Scanner Cronjob to prevent DB changes
   ```sh
   oc patch cronjob scanner -p '{"spec" : {"suspend" : true }}' --type=merge -n $NAMESPACE
   ```

3. Run a port-forward command to have access to the DB without exposing it
   publicly.
   ```sh
   oc port-forward svc/pgsql 5432:5432 -n $NAMESPACE
   ```

4. Run the backup script following this example:
   ```sh
   bash ./scripts/ciq-db-backup.sh \
     --db-name clusteriq \
     --db-user <PSQL_USER> \
     --db-passwd <PSQL_USER_PASS> \
     --db-host localhost \
     --db-port 5432 \
     --backup-dir ./backups
   ```
	 :warning: If the directory doesn't exist, it will be created

5. Resume Scanner execution
   ```sh
   oc patch cronjob scanner -p '{"spec" : {"suspend" : false }}' --type=merge -n $NAMESPACE
   ```

6. Stop port-forward process and check the backup file looks good.


## DB Restore
To restore a database backup follow this steps:
1. Log in on your Openshift environment with the `oc` CLI, and set the
   `NAMESPACE` env var.
   ```sh
   oc login ...

   export NAMESPACE="<YOUR_CLUSTER_IQ_NAMESPACE>"
   ```

2. Stop the Scanner Cronjob to prevent DB changes
   ```sh
   oc patch cronjob scanner -p '{"spec" : {"suspend" : true }}' --type=merge -n $NAMESPACE
   ```

3. Run a port-forward command to have access to the DB without exposing it
   publicly.
   ```sh
   oc port-forward svc/pgsql 5432:5432 -n $NAMESPACE
   ```

4. Run the backup script following this example:
   ```sh
   bash ./scripts/ciq-db-restore.sh \
     --db-name clusteriq \
     --db-user postgres \
     --db-passwd postgres \
     --db-admin-user postgres \
     --db-admin-passwd admin \
     --db-host localhost \
     --db-port 5432 \
     --backup <PATH_TO_BACKUP_FILE>
   ```

5. Resume scanner CronJob
   ```sh
   oc patch cronjob scanner -p '{"spec" : {"suspend" : false }}' --type=merge -n $NAMESPACE
   ```

6. Stop port-forward process and check the database was correctly restored


### Example for devel DB instance
This examples considers the PGSQL connection properties for the development environment

#### Backing Up 
```sh
bash ./scripts/ciq-db-backup.sh \
  --db-name clusteriq \
  --db-user user \
  --db-passwd postgres \
  --db-host localhost \
  --db-port 5433 \
  --backup-dir ./backups
```

#### Restoring
```sh
bash ./scripts/ciq-db-restore.sh \
  --db-name clusteriq \
  --db-user user \
  --db-passwd password \
  --db-admin-user postgres \
  --db-admin-passwd admin \
  --db-host localhost \
  --db-port 5432 \
  --backup <PATH_TO_BACKUP_FILE>
```
