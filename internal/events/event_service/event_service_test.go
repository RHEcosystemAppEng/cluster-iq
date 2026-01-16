package eventservice

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var errTest = errors.New("test error")

type mockEventRepo struct {
	listSystemEventsFn  func(ctx context.Context, opts models.ListOptions) ([]db.SystemEventDBResponse, int, error)
	listClusterEventsFn func(ctx context.Context, opts models.ListOptions) ([]db.ClusterEventDBResponse, int, error)
	createEventFn       func(ctx context.Context, event events.Event) (int64, error)
	updateEventStatusFn func(ctx context.Context, eventID int64, result string) error

	createEventCalls       int
	updateEventStatusCalls int

	lastCreateEventCtx   context.Context
	lastCreateEventEvent events.Event

	lastUpdateStatusCtx    context.Context
	lastUpdateStatusID     int64
	lastUpdateStatusResult string
}

func (m *mockEventRepo) ListSystemEvents(ctx context.Context, opts models.ListOptions) ([]db.SystemEventDBResponse, int, error) {
	if m.listSystemEventsFn == nil {
		return nil, 0, nil
	}
	return m.listSystemEventsFn(ctx, opts)
}

func (m *mockEventRepo) ListClusterEvents(ctx context.Context, opts models.ListOptions) ([]db.ClusterEventDBResponse, int, error) {
	if m.listClusterEventsFn == nil {
		return nil, 0, nil
	}
	return m.listClusterEventsFn(ctx, opts)
}

func (m *mockEventRepo) CreateEvent(ctx context.Context, event events.Event) (int64, error) {
	m.createEventCalls++
	m.lastCreateEventCtx = ctx
	m.lastCreateEventEvent = event

	if m.createEventFn == nil {
		return 0, nil
	}
	return m.createEventFn(ctx, event)
}

func (m *mockEventRepo) UpdateEventStatus(ctx context.Context, eventID int64, result string) error {
	m.updateEventStatusCalls++
	m.lastUpdateStatusCtx = ctx
	m.lastUpdateStatusID = eventID
	m.lastUpdateStatusResult = result

	if m.updateEventStatusFn == nil {
		return nil
	}
	return m.updateEventStatusFn(ctx, eventID, result)
}

// TestNewEventService verifies NewEventService creates a valid service.
func TestNewEventService(t *testing.T) {
	t.Run("New EventService", func(t *testing.T) { testNewEventService_Correct(t) })
}

func testNewEventService_Correct(t *testing.T) {
	svc := NewEventService(nil, zap.NewNop())

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.logger)
	assert.NotNil(t, svc.repo)
}

// TestLogEvent verifies LogEvent creates an event and forwards it to the repository.
func TestLogEvent(t *testing.T) {
	t.Run("LogEvent Success", func(t *testing.T) { testLogEvent_Success(t) })
	t.Run("LogEvent Repo error", func(t *testing.T) { testLogEvent_RepoError(t) })
}

func testLogEvent_Success(t *testing.T) {
	desc := "something happened"
	opts := EventOptions{
		TriggeredBy:  "scanner",
		Action:       actions.ActionOperation("START"),
		ResourceID:   "cluster-1",
		ResourceType: "cluster",
		Result:       ResultPending,
		Description:  &desc,
		Severity:     SeverityInfo,
	}

	repo := &mockEventRepo{
		createEventFn: func(ctx context.Context, event events.Event) (int64, error) {
			return 42, nil
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}

	before := time.Now().UTC()
	id, err := svc.LogEvent(opts)
	after := time.Now().UTC()

	assert.NoError(t, err)
	assert.Equal(t, int64(42), id)

	assert.Equal(t, 1, repo.createEventCalls)
	assert.NotNil(t, repo.lastCreateEventCtx)

	ev := repo.lastCreateEventEvent
	assert.Equal(t, opts.TriggeredBy, ev.TriggeredBy)
	assert.Equal(t, opts.Action, ev.Action)
	assert.Equal(t, opts.ResourceID, ev.ResourceID)
	assert.Equal(t, opts.ResourceType, ev.ResourceType)
	assert.Equal(t, opts.Result, ev.Result)
	assert.Equal(t, opts.Description, ev.Description)
	assert.Equal(t, opts.Severity, ev.Severity)

	assert.False(t, ev.EventTimestamp.IsZero())
	assert.True(t, (ev.EventTimestamp.Equal(before) || ev.EventTimestamp.After(before)) &&
		(ev.EventTimestamp.Equal(after) || ev.EventTimestamp.Before(after)))
}

func testLogEvent_RepoError(t *testing.T) {
	opts := EventOptions{
		TriggeredBy:  "api",
		Action:       actions.ActionOperation("STOP"),
		ResourceID:   "cluster-1",
		ResourceType: "cluster",
		Result:       ResultPending,
		Description:  nil,
		Severity:     SeverityError,
	}

	repo := &mockEventRepo{
		createEventFn: func(ctx context.Context, event events.Event) (int64, error) {
			return 0, errTest
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}

	id, err := svc.LogEvent(opts)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, 1, repo.createEventCalls)
}

// TestUpdateEventStatus verifies UpdateEventStatus forwards the update to the repository.
func TestUpdateEventStatus(t *testing.T) {
	t.Run("UpdateEventStatus Success", func(t *testing.T) { testUpdateEventStatus_Success(t) })
	t.Run("UpdateEventStatus Repo error", func(t *testing.T) { testUpdateEventStatus_RepoError(t) })
}

func testUpdateEventStatus_Success(t *testing.T) {
	repo := &mockEventRepo{
		updateEventStatusFn: func(ctx context.Context, eventID int64, result string) error {
			return nil
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}

	err := svc.UpdateEventStatus(7, ResultSuccess)
	assert.NoError(t, err)

	assert.Equal(t, 1, repo.updateEventStatusCalls)
	assert.Equal(t, int64(7), repo.lastUpdateStatusID)
	assert.Equal(t, ResultSuccess, repo.lastUpdateStatusResult)
	assert.NotNil(t, repo.lastUpdateStatusCtx)
}

func testUpdateEventStatus_RepoError(t *testing.T) {
	repo := &mockEventRepo{
		updateEventStatusFn: func(ctx context.Context, eventID int64, result string) error {
			return errTest
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}

	err := svc.UpdateEventStatus(7, ResultFailed)
	assert.Error(t, err)

	assert.Equal(t, 1, repo.updateEventStatusCalls)
	assert.Equal(t, int64(7), repo.lastUpdateStatusID)
	assert.Equal(t, ResultFailed, repo.lastUpdateStatusResult)
}

// TestStartTracking verifies StartTracking creates a tracker when initial log succeeds.
func TestStartTracking(t *testing.T) {
	t.Run("StartTracking Success", func(t *testing.T) { testStartTracking_Success(t) })
	t.Run("StartTracking LogEvent error", func(t *testing.T) { testStartTracking_LogEventError(t) })
}

func testStartTracking_Success(t *testing.T) {
	opts := &EventOptions{
		TriggeredBy:  "agent",
		Action:       actions.ActionOperation("START"),
		ResourceID:   "cluster-1",
		ResourceType: "cluster",
		Result:       ResultPending,
		Description:  nil,
		Severity:     SeverityInfo,
	}

	repo := &mockEventRepo{
		createEventFn: func(ctx context.Context, event events.Event) (int64, error) {
			return 99, nil
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}

	tracker := svc.StartTracking(opts)
	assert.NotNil(t, tracker)
	assert.Equal(t, int64(99), tracker.eventID)
	assert.NotNil(t, tracker.service)
	assert.NotNil(t, tracker.logger)
}

func testStartTracking_LogEventError(t *testing.T) {
	opts := &EventOptions{
		TriggeredBy:  "agent",
		Action:       actions.ActionOperation("STOP"),
		ResourceID:   "cluster-1",
		ResourceType: "cluster",
		Result:       ResultPending,
		Description:  nil,
		Severity:     SeverityError,
	}

	repo := &mockEventRepo{
		createEventFn: func(ctx context.Context, event events.Event) (int64, error) {
			return 0, errTest
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}

	tracker := svc.StartTracking(opts)
	assert.Nil(t, tracker)
}

// TestEventTracker verifies EventTracker updates status to Success/Failed and handles repo errors.
func TestEventTracker(t *testing.T) {
	t.Run("Tracker Success update ok", func(t *testing.T) { testEventTracker_Success_OK(t) })
	t.Run("Tracker Success update error", func(t *testing.T) { testEventTracker_Success_Error(t) })
	t.Run("Tracker Failed update ok", func(t *testing.T) { testEventTracker_Failed_OK(t) })
	t.Run("Tracker Failed update error", func(t *testing.T) { testEventTracker_Failed_Error(t) })
}

func testEventTracker_Success_OK(t *testing.T) {
	repo := &mockEventRepo{
		updateEventStatusFn: func(ctx context.Context, eventID int64, result string) error {
			return nil
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}
	tracker := &EventTracker{eventID: 1, service: svc, logger: zap.NewNop()}

	assert.NotPanics(t, func() { tracker.Success() })
	assert.Equal(t, 1, repo.updateEventStatusCalls)
	assert.Equal(t, ResultSuccess, repo.lastUpdateStatusResult)
}

func testEventTracker_Success_Error(t *testing.T) {
	repo := &mockEventRepo{
		updateEventStatusFn: func(ctx context.Context, eventID int64, result string) error {
			return errTest
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}
	tracker := &EventTracker{eventID: 1, service: svc, logger: zap.NewNop()}

	assert.NotPanics(t, func() { tracker.Success() })
	assert.Equal(t, 1, repo.updateEventStatusCalls)
	assert.Equal(t, ResultSuccess, repo.lastUpdateStatusResult)
}

func testEventTracker_Failed_OK(t *testing.T) {
	repo := &mockEventRepo{
		updateEventStatusFn: func(ctx context.Context, eventID int64, result string) error {
			return nil
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}
	tracker := &EventTracker{eventID: 2, service: svc, logger: zap.NewNop()}

	assert.NotPanics(t, func() { tracker.Failed() })
	assert.Equal(t, 1, repo.updateEventStatusCalls)
	assert.Equal(t, ResultFailed, repo.lastUpdateStatusResult)
}

func testEventTracker_Failed_Error(t *testing.T) {
	repo := &mockEventRepo{
		updateEventStatusFn: func(ctx context.Context, eventID int64, result string) error {
			return errTest
		},
	}

	svc := &EventService{repo: repo, logger: zap.NewNop()}
	tracker := &EventTracker{eventID: 2, service: svc, logger: zap.NewNop()}

	assert.NotPanics(t, func() { tracker.Failed() })
	assert.Equal(t, 1, repo.updateEventStatusCalls)
	assert.Equal(t, ResultFailed, repo.lastUpdateStatusResult)
}
