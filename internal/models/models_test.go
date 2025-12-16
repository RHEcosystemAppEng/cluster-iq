package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// TestFromDBScheduledActionToActions verifies FromDBScheduledActionToActions transforms DB actions into typed actions.
func TestFromDBScheduledActionToActions(t *testing.T) {
	t.Run("Empty input", func(t *testing.T) { testFromDBScheduledActionToActions_Empty(t) })
	t.Run("ScheduledAction only", func(t *testing.T) { testFromDBScheduledActionToActions_ScheduledOnly(t) })
	t.Run("CronAction only", func(t *testing.T) { testFromDBScheduledActionToActions_CronOnly(t) })
	t.Run("Mixed types", func(t *testing.T) { testFromDBScheduledActionToActions_Mixed(t) })
	t.Run("Unknown type ignored", func(t *testing.T) { testFromDBScheduledActionToActions_UnknownTypeIgnored(t) })
	t.Run("Scheduled invalid timestamp ignored", func(t *testing.T) { testFromDBScheduledActionToActions_ScheduledInvalidTimestampIgnored(t) })
	t.Run("Cron empty expression ignored", func(t *testing.T) { testFromDBScheduledActionToActions_CronEmptyExpressionIgnored(t) })
}

func testFromDBScheduledActionToActions_Empty(t *testing.T) {
	res := FromDBScheduledActionToActions(nil)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)
}

func testFromDBScheduledActionToActions_ScheduledOnly(t *testing.T) {
	now := time.Now().UTC()
	dbActions := []DBScheduledAction{
		{
			ID:          "id-1",
			Type:        "scheduled_action",
			Timestamp:   sql.NullTime{Time: now, Valid: true},
			Operation:   actions.ActionOperation("START"),
			ClusterID:   "cluster-1",
			Region:      "eu-west-1",
			AccountName: "acc-1",
			Instances:   pq.StringArray{"i-1", "i-2"},
			Status:      "Pending",
			Enable:      true,
		},
	}

	res := FromDBScheduledActionToActions(dbActions)
	assert.Len(t, res, 1)

	a, ok := res[0].(*actions.ScheduledAction)
	assert.True(t, ok)
	assert.NotNil(t, a)
	assert.Equal(t, "id-1", a.ID)
	assert.Equal(t, actions.ActionOperation("START"), a.Operation)
	assert.Equal(t, "scheduled_action", a.Type)
	assert.Equal(t, now, a.When)
	assert.Equal(t, "Pending", a.Status)
	assert.True(t, a.Enabled)
	assert.Equal(t, "cluster-1", a.Target.ClusterID)
	assert.Equal(t, "eu-west-1", a.Target.Region)
	assert.Equal(t, "acc-1", a.Target.AccountID)
	assert.Equal(t, []string{"i-1", "i-2"}, a.Target.Instances)
}

func testFromDBScheduledActionToActions_CronOnly(t *testing.T) {
	dbActions := []DBScheduledAction{
		{
			ID:             "id-2",
			Type:           "cron_action",
			CronExpression: sql.NullString{String: "0 0 * * *", Valid: true},
			Operation:      actions.ActionOperation("STOP"),
			ClusterID:      "cluster-2",
			Region:         "us-east-1",
			AccountName:    "acc-2",
			Instances:      pq.StringArray{"i-9"},
			Status:         "Pending",
			Enable:         false,
		},
	}

	res := FromDBScheduledActionToActions(dbActions)
	assert.Len(t, res, 1)

	a, ok := res[0].(*actions.CronAction)
	assert.True(t, ok)
	assert.NotNil(t, a)
	assert.Equal(t, "id-2", a.ID)
	assert.Equal(t, actions.ActionOperation("STOP"), a.Operation)
	assert.Equal(t, "cron_action", a.Type)
	assert.Equal(t, "0 0 * * *", a.Expression)
	assert.Equal(t, "Pending", a.Status)
	assert.False(t, a.Enabled)
	assert.Equal(t, "cluster-2", a.Target.ClusterID)
	assert.Equal(t, "us-east-1", a.Target.Region)
	assert.Equal(t, "acc-2", a.Target.AccountID)
	assert.Equal(t, []string{"i-9"}, a.Target.Instances)
}

func testFromDBScheduledActionToActions_Mixed(t *testing.T) {
	now := time.Now().UTC()
	dbActions := []DBScheduledAction{
		{
			ID:          "id-1",
			Type:        "scheduled_action",
			Timestamp:   sql.NullTime{Time: now, Valid: true},
			Operation:   actions.ActionOperation("START"),
			ClusterID:   "cluster-1",
			Region:      "eu-west-1",
			AccountName: "acc-1",
			Instances:   pq.StringArray{"i-1"},
			Status:      "Pending",
			Enable:      true,
		},
		{
			ID:             "id-2",
			Type:           "cron_action",
			CronExpression: sql.NullString{String: "*/5 * * * *", Valid: true},
			Operation:      actions.ActionOperation("STOP"),
			ClusterID:      "cluster-2",
			Region:         "us-east-1",
			AccountName:    "acc-2",
			Instances:      pq.StringArray{"i-2"},
			Status:         "Running",
			Enable:         false,
		},
	}

	res := FromDBScheduledActionToActions(dbActions)
	assert.Len(t, res, 2)

	_, ok0 := res[0].(*actions.ScheduledAction)
	_, ok1 := res[1].(*actions.CronAction)
	assert.True(t, ok0)
	assert.True(t, ok1)
}

func testFromDBScheduledActionToActions_UnknownTypeIgnored(t *testing.T) {
	dbActions := []DBScheduledAction{
		{ID: "id-x", Type: "instant_action"},
	}

	res := FromDBScheduledActionToActions(dbActions)
	assert.Len(t, res, 0)
}

func testFromDBScheduledActionToActions_ScheduledInvalidTimestampIgnored(t *testing.T) {
	dbActions := []DBScheduledAction{
		{
			ID:        "id-1",
			Type:      "scheduled_action",
			Timestamp: sql.NullTime{Valid: false},
		},
	}

	res := FromDBScheduledActionToActions(dbActions)
	assert.Len(t, res, 1)
	assert.Nil(t, res[0])
}

func testFromDBScheduledActionToActions_CronEmptyExpressionIgnored(t *testing.T) {
	dbActions := []DBScheduledAction{
		{
			ID:             "id-2",
			Type:           "cron_action",
			CronExpression: sql.NullString{String: "", Valid: true},
		},
	}

	res := FromDBScheduledActionToActions(dbActions)
	assert.Len(t, res, 1)
	assert.Nil(t, res[0])
}

// TestFromDBScheduledActionToScheduledAction verifies FromDBScheduledActionToScheduledAction mapping.
func TestFromDBScheduledActionToScheduledAction(t *testing.T) {
	t.Run("ScheduledAction valid timestamp", func(t *testing.T) { testFromDBScheduledActionToScheduledAction_Valid(t) })
	t.Run("ScheduledAction invalid timestamp", func(t *testing.T) { testFromDBScheduledActionToScheduledAction_InvalidTimestamp(t) })
}

func testFromDBScheduledActionToScheduledAction_Valid(t *testing.T) {
	now := time.Now().UTC()
	dbAction := DBScheduledAction{
		ID:          "id-1",
		Type:        "scheduled_action",
		Timestamp:   sql.NullTime{Time: now, Valid: true},
		Operation:   actions.ActionOperation("START"),
		ClusterID:   "cluster-1",
		Region:      "eu-west-1",
		AccountName: "acc-1",
		Instances:   pq.StringArray{"i-1", "i-2"},
		Status:      "Pending",
		Enable:      true,
	}

	a := FromDBScheduledActionToScheduledAction(dbAction)
	assert.NotNil(t, a)
	assert.Equal(t, "id-1", a.ID)
	assert.Equal(t, actions.ActionOperation("START"), a.Operation)
	assert.Equal(t, "scheduled_action", a.Type)
	assert.Equal(t, now, a.When)
	assert.Equal(t, "Pending", a.Status)
	assert.True(t, a.Enabled)
	assert.Equal(t, "cluster-1", a.Target.ClusterID)
	assert.Equal(t, "eu-west-1", a.Target.Region)
	assert.Equal(t, "acc-1", a.Target.AccountID)
	assert.Equal(t, []string{"i-1", "i-2"}, a.Target.Instances)
}

func testFromDBScheduledActionToScheduledAction_InvalidTimestamp(t *testing.T) {
	dbAction := DBScheduledAction{
		ID:        "id-1",
		Type:      "scheduled_action",
		Timestamp: sql.NullTime{Valid: false},
	}

	a := FromDBScheduledActionToScheduledAction(dbAction)
	assert.Nil(t, a)
}

// TestFromDBScheduledActionToCronAction verifies FromDBScheduledActionToCronAction mapping.
func TestFromDBScheduledActionToCronAction(t *testing.T) {
	t.Run("CronAction valid expression", func(t *testing.T) { testFromDBScheduledActionToCronAction_Valid(t) })
	t.Run("CronAction empty expression", func(t *testing.T) { testFromDBScheduledActionToCronAction_EmptyExpression(t) })
}

func testFromDBScheduledActionToCronAction_Valid(t *testing.T) {
	dbAction := DBScheduledAction{
		ID:             "id-2",
		Type:           "cron_action",
		CronExpression: sql.NullString{String: "0 0 * * *", Valid: true},
		Operation:      actions.ActionOperation("STOP"),
		ClusterID:      "cluster-2",
		Region:         "us-east-1",
		AccountName:    "acc-2",
		Instances:      pq.StringArray{"i-9"},
		Status:         "Pending",
		Enable:         false,
	}

	a := FromDBScheduledActionToCronAction(dbAction)
	assert.NotNil(t, a)
	assert.Equal(t, "id-2", a.ID)
	assert.Equal(t, actions.ActionOperation("STOP"), a.Operation)
	assert.Equal(t, "cron_action", a.Type)
	assert.Equal(t, "0 0 * * *", a.Expression)
	assert.Equal(t, "Pending", a.Status)
	assert.False(t, a.Enabled)
	assert.Equal(t, "cluster-2", a.Target.ClusterID)
	assert.Equal(t, "us-east-1", a.Target.Region)
	assert.Equal(t, "acc-2", a.Target.AccountID)
	assert.Equal(t, []string{"i-9"}, a.Target.Instances)
}

func testFromDBScheduledActionToCronAction_EmptyExpression(t *testing.T) {
	dbAction := DBScheduledAction{
		ID:             "id-2",
		Type:           "cron_action",
		CronExpression: sql.NullString{String: "", Valid: true},
	}

	a := FromDBScheduledActionToCronAction(dbAction)
	assert.Nil(t, a)
}
