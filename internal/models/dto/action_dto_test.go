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
	t.Run("ToModelAction Instant", func(t *testing.T) { testActionDTORequest_ToModelAction_Instant(t) })
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
		Status:    string(actions.StatusPending),
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

	assert.Equal(t, "", sa.ID)
	assert.Equal(t, actions.ActionOperation(dto.Operation), sa.Operation)
	assert.Equal(t, actions.StatusPending, sa.Status)
	assert.Equal(t, dto.Enabled, sa.Enabled)
	assert.Equal(t, "", sa.Requester)
	assert.Nil(t, sa.Description)

	assert.Equal(t, dto.AccountID, sa.Target.AccountID)
	assert.Equal(t, dto.Region, sa.Target.Region)
	assert.Equal(t, dto.ClusterID, sa.Target.ClusterID)
	assert.Equal(t, dto.Instances, sa.Target.Instances)

	assert.Equal(t, dto.Time, sa.When)
	assert.Equal(t, actions.ScheduledActionType, sa.GetType())
}

func testActionDTORequest_ToModelAction_Cron(t *testing.T) {
	dto := ActionDTORequest{
		ID:        "id-2",
		Type:      string(actions.CronActionType),
		Time:      time.Time{},
		CronExp:   "0 0 * * *",
		Operation: "STOP",
		Status:    string(actions.StatusPending),
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

	assert.Equal(t, "", ca.ID)
	assert.Equal(t, actions.ActionOperation(dto.Operation), ca.Operation)
	assert.Equal(t, actions.StatusPending, ca.Status)
	assert.Equal(t, dto.Enabled, ca.Enabled)
	assert.Equal(t, "", ca.Requester)
	assert.Nil(t, ca.Description)

	assert.Equal(t, dto.AccountID, ca.Target.AccountID)
	assert.Equal(t, dto.Region, ca.Target.Region)
	assert.Equal(t, dto.ClusterID, ca.Target.ClusterID)
	assert.Equal(t, dto.Instances, ca.Target.Instances)

	assert.Equal(t, dto.CronExp, ca.Expression)
	assert.Equal(t, actions.CronActionType, ca.GetType())
}

func testActionDTORequest_ToModelAction_Instant(t *testing.T) {
	dto := ActionDTORequest{
		ID:        "id-2",
		Type:      string(actions.InstantActionType),
		Time:      time.Time{},
		CronExp:   "0 0 * * *",
		Operation: "STOP",
		Status:    string(actions.StatusPending),
		Enabled:   false,
		ClusterID: "cluster-2",
		Region:    "us-east-1",
		AccountID: "acc-2",
		Instances: []string{"i-9"},
	}

	act := dto.ToModelAction()
	assert.NotNil(t, act)

	ia, ok := act.(*actions.InstantAction)
	assert.True(t, ok)
	assert.NotNil(t, ia)

	assert.Equal(t, "", ia.ID)
	assert.Equal(t, actions.ActionOperation(dto.Operation), ia.Operation)
	assert.Equal(t, actions.StatusPending, ia.Status)
	assert.Equal(t, dto.Enabled, ia.Enabled)
	assert.Equal(t, "", ia.Requester)
	assert.Nil(t, ia.Description)

	assert.Equal(t, dto.AccountID, ia.Target.AccountID)
	assert.Equal(t, dto.Region, ia.Target.Region)
	assert.Equal(t, dto.ClusterID, ia.Target.ClusterID)
	assert.Equal(t, dto.Instances, ia.Target.Instances)

	assert.Equal(t, actions.InstantActionType, ia.GetType())
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

// TestActionDTOResponse_ToModelAction verifies ActionDTOResponse conversion to actions.Action.
func TestActionDTOResponse_ToModelAction(t *testing.T) {
	t.Run("ToModelAction ScheduledAction", func(t *testing.T) { testActionDTOResponse_ToModelAction_Scheduled(t) })
	t.Run("ToModelAction CronAction", func(t *testing.T) { testActionDTOResponse_ToModelAction_Cron(t) })
	t.Run("ToModelAction Instant", func(t *testing.T) { testActionDTOResponse_ToModelAction_Instant(t) })
	t.Run("ToModelAction unknown type", func(t *testing.T) { testActionDTOResponse_ToModelAction_UnknownType(t) })
}

func testActionDTOResponse_ToModelAction_Scheduled(t *testing.T) {
	now := time.Now().UTC()

	dto := ActionDTOResponse{
		ID:        "id-1",
		Type:      string(actions.ScheduledActionType),
		Time:      now,
		CronExp:   "",
		Operation: "PowerOnCluster",
		Status:    string(actions.StatusPending),
		Enabled:   true,
		ClusterID: "cluster-1",
		Region:    "eu-west-1",
		AccountID: "acc-1",
		Instances: []string{"i-1", "i-2"},
	}

	a := dto.ToModelAction()
	assert.NotNil(t, a)

	sa, ok := a.(*actions.ScheduledAction)
	assert.True(t, ok)

	assert.Equal(t, "", sa.ID)
	assert.Equal(t, actions.ActionOperation(dto.Operation), sa.Operation)
	assert.Equal(t, actions.StatusPending, sa.Status)
	assert.Equal(t, dto.Enabled, sa.Enabled)

	assert.Equal(t, dto.AccountID, sa.Target.AccountID)
	assert.Equal(t, dto.Region, sa.Target.Region)
	assert.Equal(t, dto.ClusterID, sa.Target.ClusterID)
	assert.Equal(t, dto.Instances, sa.Target.Instances)

	assert.Equal(t, dto.Time, sa.When)
	assert.Equal(t, actions.ScheduledActionType, sa.GetType())
}

func testActionDTOResponse_ToModelAction_Cron(t *testing.T) {
	dto := ActionDTOResponse{
		ID:        "id-2",
		Type:      string(actions.CronActionType),
		Time:      time.Time{},
		CronExp:   "*/5 * * * *",
		Operation: "PowerOffCluster",
		Status:    string(actions.StatusPending),
		Enabled:   false,
		ClusterID: "cluster-2",
		Region:    "us-east-1",
		AccountID: "acc-2",
		Instances: []string{"i-9"},
	}

	a := dto.ToModelAction()
	assert.NotNil(t, a)

	ca, ok := a.(*actions.CronAction)
	assert.True(t, ok)

	assert.Equal(t, "", ca.ID)
	assert.Equal(t, actions.ActionOperation(dto.Operation), ca.Operation)
	assert.Equal(t, actions.StatusPending, ca.Status)
	assert.Equal(t, dto.Enabled, ca.Enabled)

	assert.Equal(t, dto.AccountID, ca.Target.AccountID)
	assert.Equal(t, dto.Region, ca.Target.Region)
	assert.Equal(t, dto.ClusterID, ca.Target.ClusterID)
	assert.Equal(t, dto.Instances, ca.Target.Instances)

	assert.Equal(t, dto.CronExp, ca.Expression)
	assert.Equal(t, actions.CronActionType, ca.GetType())
}

func testActionDTOResponse_ToModelAction_Instant(t *testing.T) {
	dto := ActionDTOResponse{
		ID:        "id-2",
		Type:      string(actions.InstantActionType),
		Time:      time.Time{},
		CronExp:   "0 0 * * *",
		Operation: "STOP",
		Status:    string(actions.StatusPending),
		Enabled:   false,
		ClusterID: "cluster-2",
		Region:    "us-east-1",
		AccountID: "acc-2",
		Instances: []string{"i-9"},
	}

	act := dto.ToModelAction()
	assert.NotNil(t, act)

	ia, ok := act.(*actions.InstantAction)
	assert.True(t, ok)
	assert.NotNil(t, ia)

	assert.Equal(t, "", ia.ID)
	assert.Equal(t, actions.ActionOperation(dto.Operation), ia.Operation)
	assert.Equal(t, actions.StatusPending, ia.Status)
	assert.Equal(t, dto.Enabled, ia.Enabled)
	assert.Equal(t, "", ia.Requester)
	assert.Nil(t, ia.Description)

	assert.Equal(t, dto.AccountID, ia.Target.AccountID)
	assert.Equal(t, dto.Region, ia.Target.Region)
	assert.Equal(t, dto.ClusterID, ia.Target.ClusterID)
	assert.Equal(t, dto.Instances, ia.Target.Instances)

	assert.Equal(t, actions.InstantActionType, ia.GetType())
}

func testActionDTOResponse_ToModelAction_UnknownType(t *testing.T) {
	dto := ActionDTOResponse{
		ID:   "id-x",
		Type: "unknown_type",
	}

	a := dto.ToModelAction()
	assert.Nil(t, a)
}

// TestToModelActionListFromResponse verifies list conversion and error handling.
func TestToModelActionListFromResponse(t *testing.T) {
	t.Run("ToModelActionListFromResponse OK", func(t *testing.T) { testToModelActionListFromResponse_OK(t) })
	t.Run("ToModelActionListFromResponse unknown type", func(t *testing.T) { testToModelActionListFromResponse_UnknownType(t) })
}

func testToModelActionListFromResponse_OK(t *testing.T) {
	now := time.Now().UTC()

	dtos := []ActionDTOResponse{
		{
			ID:        "id-1",
			Type:      string(actions.ScheduledActionType),
			Time:      now,
			Operation: "PowerOnCluster",
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
			CronExp:   "0 0 * * *",
			Operation: "PowerOffCluster",
			Status:    "Pending",
			Enabled:   false,
			ClusterID: "cluster-2",
			Region:    "us-east-1",
			AccountID: "acc-2",
			Instances: []string{"i-2"},
		},
	}

	actionsList, err := ToModelActionListFromResponse(dtos)

	assert.NoError(t, err)
	assert.NotNil(t, actionsList)
	assert.Len(t, actionsList, 2)

	_, ok0 := actionsList[0].(*actions.ScheduledAction)
	_, ok1 := actionsList[1].(*actions.CronAction)
	assert.True(t, ok0)
	assert.True(t, ok1)
}

func testToModelActionListFromResponse_UnknownType(t *testing.T) {
	dtos := []ActionDTOResponse{
		{ID: "id-1", Type: "unknown_type"},
	}

	actionsList, err := ToModelActionListFromResponse(dtos)

	assert.Error(t, err)
	assert.Nil(t, actionsList)
	assert.ErrorContains(t, err, "unknown action type: unknown_type")
}
