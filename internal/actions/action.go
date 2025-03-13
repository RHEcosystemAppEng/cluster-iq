package actions

import (
	"encoding/json"
	"fmt"
)

// Action defines the interface for cloud actions that can be executed.
// Implementations of this interface should provide details about the action type,
// target region, target resource, and a unique identifier.
type Action interface {
	// GetActionOperation returns the type of action being performed.
	//
	// Returns:
	// - An ActionOperation indicating the action type (e.g., PowerOnCluster, PowerOffCluster).
	GetActionOperation() ActionOperation

	// GetRegion returns the cloud region where the action is executed.
	//
	// Returns:
	// - A string representing the cloud region.
	GetRegion() string

	// GetTarget returns the target resource of the action.
	//
	// Returns:
	// - An ActionTarget representing the target cluster and instances affected by the action.
	GetTarget() ActionTarget

	// GetID returns a unique identifier for the action.
	//
	// Returns:
	// - A string representing the unique action ID.
	GetID() string

	// GetType returns the action type
	//
	// Returns:
	// - A string representing action type
	GetType() ActionType
}

// DecodeActions received a http response body as a []byte for decoding the
// actions on it and unmarshall them evaulating its specific action type.
// Every decoded action will be parsed as a specific action type based on its
// properties
//
// Parameters:
// - action as a byte array. It's supposed to be used for unmarshalling the body from a HTTP request
//
// Returns:
// - A slice of Action where every element was unmarshalled as it's specific Action Type
// - err if an Error occurs
func DecodeActions(actions []json.RawMessage) (*[]Action, error) {
	var resultActions []Action
	for _, action := range actions {

		// Auxiliar struct for getting the action type
		var r struct {
			Type string `json:"type"`
		}

		// Unmarshalling action type
		if err := json.Unmarshal(action, &r); err != nil {
			return nil, err
		}

		// Unmarshalling based ont Action Type
		switch r.Type {
		case SCHEDULED_ACTION_TYPE: // Unmarshall as ScheduledAction
			var a ScheduledAction
			if err := json.Unmarshal(action, &a); err != nil {
				return nil, err
			}
			resultActions = append(resultActions, a)
		case CRON_ACTION_TYPE: // Unmarshall as CronAction
			var a CronAction
			if err := json.Unmarshal(action, &a); err != nil {
				return nil, err
			}
			resultActions = append(resultActions, a)
		default:
			return nil, fmt.Errorf("Unknown Action Type: %s", r.Type)
		}
	}

	return &resultActions, nil
}

// DecodeActions takes an array of Actions and splits it in separate slices classified by ActionType
//
// Parameters:
// - A slice of Actions to split (will not be modified)
//
// Returns:
// - A pointer to a slice of ScheduledActions
// - A pointer to a slice of CronAction
func SplitActionsByType(actions []Action) ([]ScheduledAction, []CronAction) {
	var schedActions []ScheduledAction
	var cronActions []CronAction

	for _, action := range actions {
		switch action.GetType() {
		case SCHEDULED_ACTION_TYPE: // Unmarshall as ScheduledAction
			schedActions = append(schedActions, action.(ScheduledAction))
		case CRON_ACTION_TYPE: // Unmarshall as CronAction
			cronActions = append(cronActions, action.(CronAction))
		}
	}

	return schedActions, cronActions
}
