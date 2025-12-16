package dto

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/stretchr/testify/assert"
)

// TestActionDTORequest_ToModelAction verifies DTO to actions.Action conversion by type.
func TestActionDTORequest_ToModelAction(t *testing.T) {
	t.Run("ToModelAction ScheduledAction", func(t *testing.T) { testActionDTORequest_ToModelAction_Scheduled(t) })
	t.Run("ToModelAction CronAction", func(t *testing.T) { testActionDTORequest_ToModelAction_Cron(t) })
	t.Run("ToModelAction Unknown type", func(t *testing.T) { testActionDTORequest_ToModelAction_UnknownType(t) })
}

func testActionDTORequest_ToModelAction_Scheduled(t *testing.T) {
	now := time.Now().UTC()

	dto := ActionDTORequest{
		ID:        "id-1",
		Type:      string(actions.ScheduledActionType),
		Time:      now,
		CronExp:   "",
		Operation: "START",
		Status:    "Pending",
		Enabled:   true,
		ClusterID: "cluster-1",
		Region:    "eu-west-1",
		AccountID: "acc-1",
		Instances: []string{"i-1", "i-2"},
	}

	act := dto.ToModelAction()
	assert.NotNil(t, act)

	sa, ok := act.(*actions.ScheduledAction)
	assert.True(t, ok)
	assert.NotNil(t, sa)

	assert.Equal(t, dto.ID, sa.ID)
	assert.Equal(t, actions.ActionOperation(dto.Operation), sa.Operation)
	assert.Equal(t, dto.Status, sa.Status)
	assert.Equal(t, dto.Enabled, sa.Enabled)

	assert.Equal(t, dto.AccountID, sa.Target.AccountID)
	assert.Equal(t, dto.Region, sa.Target.Region)
	assert.Equal(t, dto.ClusterID, sa.Target.ClusterID)
	assert.Equal(t, dto.Instances, sa.Target.Instances)

	assert.Equal(t, dto.Time, sa.When)
	assert.Equal(t, dto.Type, sa.Type)
}

func testActionDTORequest_ToModelAction_Cron(t *testing.T) {
	dto := ActionDTORequest{
		ID:        "id-2",
		Type:      string(actions.CronActionType),
		Time:      time.Time{},
		CronExp:   "0 0 * * *",
		Operation: "STOP",
		Status:    "Pending",
		Enabled:   false,
		ClusterID: "cluster-2",
		Region:    "us-east-1",
		AccountID: "acc-2",
		Instances: []string{"i-9"},
	}

	act := dto.ToModelAction()
	assert.NotNil(t, act)

	ca, ok := act.(*actions.CronAction)
	assert.True(t, ok)
	assert.NotNil(t, ca)

	assert.Equal(t, dto.ID, ca.ID)
	assert.Equal(t, actions.ActionOperation(dto.Operation), ca.Operation)
	assert.Equal(t, dto.Status, ca.Status)
	assert.Equal(t, dto.Enabled, ca.Enabled)

	assert.Equal(t, dto.AccountID, ca.Target.AccountID)
	assert.Equal(t, dto.Region, ca.Target.Region)
	assert.Equal(t, dto.ClusterID, ca.Target.ClusterID)
	assert.Equal(t, dto.Instances, ca.Target.Instances)

	assert.Equal(t, dto.CronExp, ca.Expression)
	assert.Equal(t, dto.Type, ca.Type)
}

func testActionDTORequest_ToModelAction_UnknownType(t *testing.T) {
	dto := ActionDTORequest{
		ID:   "id-x",
		Type: "unknown_action",
	}

	act := dto.ToModelAction()
	assert.Nil(t, act)
}

// TestToModelActionList verifies slice conversion from DTOs to actions.Action.
func TestToModelActionList(t *testing.T) {
	t.Run("Action list success", func(t *testing.T) { testToModelActionList_Correct(t) })
	t.Run("Action list with invalid action", func(t *testing.T) { testToModelActionList_Invalid(t) })
}

func testToModelActionList_Correct(t *testing.T) {
	now := time.Now().UTC()

	dtos := []ActionDTORequest{
		{
			ID:        "id-1",
			Type:      string(actions.ScheduledActionType),
			Time:      now,
			Operation: "START",
			Status:    "Pending",
			Enabled:   true,
			ClusterID: "cluster-1",
			Region:    "eu-west-1",
			AccountID: "acc-1",
			Instances: []string{"i-1"},
		},
		{
			ID:        "id-2",
			Type:      string(actions.CronActionType),
			CronExp:   "*/5 * * * *",
			Operation: "STOP",
			Status:    "Pending",
			Enabled:   false,
			ClusterID: "cluster-2",
			Region:    "us-east-1",
			AccountID: "acc-2",
			Instances: []string{"i-2"},
		},
	}

	list := ToModelActionList(dtos)
	assert.NotNil(t, list)
	assert.Len(t, *list, 2)

	_, ok0 := (*list)[0].(*actions.ScheduledAction)
	_, ok1 := (*list)[1].(*actions.CronAction)
	assert.True(t, ok0)
	assert.True(t, ok1)
}

func testToModelActionList_Invalid(t *testing.T) {
	dtos := []ActionDTORequest{
		{
			ID:   "id-1",
			Type: "unknown_action",
		},
	}

	list := ToModelActionList(dtos)
	assert.Nil(t, list)
}
