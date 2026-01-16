package actions

type ActionType string

const (
	// ScheduledActionType code for labeling ScheduledActions
	ScheduledActionType ActionType = "scheduled_action"
	// CronActionType code for labeling CronActions
	CronActionType ActionType = "cron_action"
	// InstantActionType code for labeling InstantActions. InstantActions are NOT stored on DB
	InstantActionType ActionType = "instant_action"
)
