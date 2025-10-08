package repositories

const (
	ScheduleTable = "schedule"

	SelectScheduleFullView = "schedule_full_view"

	// SelectScheduledActionsQuery returns the list of scheduled actions on the inventory with all the parameters needed for action execution
	// ARRAY_AGG is used for joining every instance on the same row
	SelectScheduledActionsQuery = `
		SELECT * FROM (
			SELECT
				schedule.id,
				schedule.type,
				schedule.time,
				schedule.cron_exp,
				schedule.operation,
				schedule.status,
				schedule.enabled,
				clusters.id AS cluster_id,
				clusters.region,
				clusters.account_name,
				ARRAY_AGG(CAST(instances.id AS TEXT)) FILTER (WHERE instances IS NOT NULL) AS instances
			FROM schedule
			JOIN clusters ON schedule.target = clusters.id
			JOIN instances ON clusters.id = instances.cluster_id
			GROUP BY
				schedule.id,
				schedule.type,
				schedule.time,
				schedule.cron_exp,
				schedule.operation,
				schedule.status,
				schedule.enabled,
				clusters.id,
				clusters.region,
				clusters.account_name
		) AS subquery
		`

	// EnableActionQuery enables the action to be re-scheduled on next agent polling
	EnableActionQuery = `
		UPDATE
			schedule
		SET
			enabled = true
		WHERE id = $1
`

	// DisableActionQuery disables the action to don't be re-scheduled on next agent polling
	DisableActionQuery = `
		UPDATE
			schedule
		SET
			enabled = false
		WHERE id = $1
`

	// SelectScheduledActionByIDQuery returns scheduled action on the inventory for a specific ID
	SelectScheduledActionsByIDQuery = `
		SELECT
			schedule.id,
			schedule.type,
		  schedule.time,
		  schedule.cron_exp,
			schedule.operation,
			schedule.status,
			schedule.enabled,
			clusters.id AS cluster_id,
			clusters.region,
			clusters.account_name,
		ARRAY_AGG(CAST(instances.id AS TEXT)) FILTER (WHERE instances IS NOT NULL) AS instances
		FROM schedule
		JOIN clusters ON schedule.target = clusters.id
		JOIN instances ON clusters.id = instances.cluster_id
		WHERE schedule.id = $1
		GROUP BY
			schedule.id,
			schedule.type,
			schedule.time,
			schedule.cron_exp,
			schedule.operation,
			schedule.status,
			schedule.enabled,
			clusters.id,
			clusters.region,
			clusters.account_name
		`

	// InsertScheduledActionQuery inserts new scheduled actions on the DB
	InsertScheduledActionsQuery = `
		INSERT INTO schedule (
			type,
			time,
			operation,
			target,
			status,
			enabled
		) VALUES (
			:type,
			:time,
			:operation,
			:target.cluster_id,
			:status,
			:enabled
		)
	`
	// InsertCronActionQuery inserts new Cron actions on the DB
	InsertCronActionsQuery = `
		INSERT INTO schedule (
			type,
			cron_exp,
			operation,
			target,
			status,
			enabled
		) VALUES (
			:type,
			:cron_exp,
			:operation,
			:target.cluster_id,
			:status,
			:enabled
		)
	`

	// PatchScheduledActionsQuery
	PatchScheduledActionsQuery = `
		UPDATE
			schedule
		SET
			time = :time,
			operation = :operation,
			target = :target.cluster_id,
			enabled = :enabled
		WHERE
			id = :id
	`

	// PatchCronActionsQuery
	PatchCronActionsQuery = `
		UPDATE
			schedule
		SET
			cron_exp = :cron_exp,
			operation = :operation,
			target = :target.cluster_id,
			enabled = :enabled
		WHERE
			id = :id
	`

	// PatchActionStatusQuery
	PatchActionStatusQuery = `
		UPDATE
			schedule
		SET
			status = $2,
			enabled = $3
		WHERE
			id = $1
	`

	// DeleteScheduledActionQuery
	DeleteScheduledActionsQuery = `DELETE FROM schedule WHERE id=$1`
)
