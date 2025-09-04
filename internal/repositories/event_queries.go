package repositories

const (
	// InsertEventQuery insert a new audit event
	InsertEventQuery = `
		INSERT INTO audit_logs(
			event_timestamp,
			triggered_by,
			action_name,
			resource_id,
			resource_type,
			result,
			description,
			severity
		) VALUES (
			CURRENT_TIMESTAMP,
			:triggered_by,
			:action_name,
			:resource_id,
			:resource_type,
			:result,
			:description,
			:severity
		) RETURNING id
	`

	// SelectClusterEventsQuery returns audit log events related to a specific cluster.
	SelectClusterEventsQuery = `
		SELECT 
			al.id, 
			al.event_timestamp, 
			al.triggered_by, 
			al.action_name, 
			al.resource_id, 
			al.resource_type, 
			al.result, 
			al.description, 
			al.severity
		FROM audit_logs al
		WHERE al.resource_id = $1
		ORDER BY al.event_timestamp DESC
	`
	// SelectSystemEventsQuery returns system-wide audit logs.
	SelectSystemEventsQuery = `
		SELECT 
			al.id, 
			al.event_timestamp, 
			al.triggered_by, 
			al.action_name, 
			al.resource_id, 
			al.resource_type, 
			al.result, 
			al.description, 
			al.severity,
			acc.id AS account_id,
			acc.provider
		FROM audit_logs al
		LEFT JOIN accounts acc ON acc.name = (
			CASE 
				WHEN al.resource_type = 'cluster' 
				THEN (SELECT c.account_name FROM clusters c WHERE c.id = al.resource_id)
				WHEN al.resource_type = 'instance' 
				THEN (SELECT c.account_name FROM clusters c WHERE c.id = (SELECT i.cluster_id FROM instances i WHERE i.id = al.resource_id))
			END
		)
	`
	// UpdateEventStatusQuery updates the result status of an audit log entry based on its ID.
	UpdateEventStatusQuery = `UPDATE audit_logs SET result=$1 WHERE id=$2`
)
