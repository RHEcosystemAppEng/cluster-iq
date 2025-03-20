package actions

type ActionType string

const (
	// ScheduledActionType code for labeling ScheduledActions
	ScheduledActionType = "scheduled_action"
	// CronActionType code for labeling CronActions
	CronActionType = "cron_action"
	// InstantActionType code for labeling InstantActions. InstantActions are NOT stored on DB
	InstantActionType = "instant_action"
)
