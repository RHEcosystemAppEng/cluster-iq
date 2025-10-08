package repositories

const (
	EventsTable = "events"

	SelectClusterEventsView = "cluster_events"

	SelectSystemEventsView = "system_events"

	// InsertEventQuery insert a new audit event
	InsertEventQuery = `
		INSERT INTO events(
			event_timestamp,
			triggered_by,
			action,
			resource_id,
			resource_type,
			result,
			description,
			severity
		) VALUES (
			CURRENT_TIMESTAMP,
			:triggered_by,
			:action,
			(
				CASE
					WHEN :resource_type = 'cluster'
					THEN (SELECT id FROM clusters c WHERE c.cluster_id = :resource_id)
					WHEN :resource_type = 'instance'
					THEN (SELECT id FROM instances i WHERE i.instance_id = :resource_id)
				END
			),
			:resource_type,
			:result,
			:description,
			:severity
		) RETURNING id
	`

	// UpdateEventStatusQuery updates the result status of an audit log entry based on its ID.
	UpdateEventStatusQuery = `UPDATE events SET result=:result WHERE id=:id`
)
