package repositories

import (
	"context"
	"fmt"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
)

const (
	// DB Table for events
	EventsTable = "events"
	// View for SELECT operations of cluster related events
	SelectClusterEventsView = "cluster_events"
	// View for SELECT operations of system-wide events
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

var _ EventRepository = (*eventRepositoryImpl)(nil)

// EventRepository defines the interface for data access operations for events.
type EventRepository interface {
	ListSystemEvents(ctx context.Context, opts models.ListOptions) ([]db.SystemEventDBResponse, int, error)
	ListClusterEvents(ctx context.Context, opts models.ListOptions) ([]db.ClusterEventDBResponse, int, error)
	CreateEvent(ctx context.Context, events events.Event) (int64, error)
	UpdateEventStatus(ctx context.Context, eventID int64, result string) error
}

type eventRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewEventRepository(db *dbclient.DBClient) EventRepository {
	return &eventRepositoryImpl{db: db}
}

// ListSystemEvents retrieves system-wide events with extended metadata.
func (r *eventRepositoryImpl) ListSystemEvents(ctx context.Context, opts models.ListOptions) ([]db.SystemEventDBResponse, int, error) {
	var events []db.SystemEventDBResponse

	if err := r.db.SelectWithContext(ctx, &events, SelectSystemEventsView, opts, "event_timestamp", "*"); err != nil {
		return events, 0, fmt.Errorf("failed to list events: %w", err)
	}

	return events, len(events), nil
}

// ListClusterEvents retrieves events for a specific resource (like a cluster).
func (r *eventRepositoryImpl) ListClusterEvents(ctx context.Context, opts models.ListOptions) ([]db.ClusterEventDBResponse, int, error) {
	var events []db.ClusterEventDBResponse
	var id int

	if err := r.db.GetWithContext(ctx, &id, ClustersTable, models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters:  map[string]interface{}{"cluster_id": opts.Filters["resource_id"]},
	}, "id"); err != nil {
		return events, 0, fmt.Errorf("failed to get cluster internal id: %w", err)
	}

	opts.Filters["resource_id"] = id

	if err := r.db.SelectWithContext(ctx, &events, SelectClusterEventsView, opts, "event_timestamp", "*"); err != nil {
		return events, 0, fmt.Errorf("failed to list cluster events: %w", err)
	}

	return events, len(events), nil
}

// AddEvent inserts a new audit event into the database and returns the event ID.
func (r *eventRepositoryImpl) CreateEvent(ctx context.Context, event events.Event) (int64, error) {
	if err := r.db.InsertWithContext(ctx, InsertEventQuery, event); err != nil {
		return -1, err
	}

	return 0, nil
}

// UpdateEventStatus updates the result status of an audit event.
func (r *eventRepositoryImpl) UpdateEventStatus(ctx context.Context, eventID int64, result string) error {
	if err := r.db.NamedUpdateWithContext(ctx, UpdateEventStatusQuery, events.Event{ID: eventID, Result: result}); err != nil {
		return err
	}

	return nil
}
