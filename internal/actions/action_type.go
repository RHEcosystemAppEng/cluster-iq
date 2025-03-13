package actions

type ActionType string

const (
	// SCHEDULED_ACTION_TYPE code for labeling ScheduledActions
	SCHEDULED_ACTION_TYPE = "scheduled_action"
	// Cron_ACTION_TYPE code for labeling CronActions
	CRON_ACTION_TYPE = "cron_action"
	// INSTANT_ACTION_TYPE code for labeling InstantActions. InstantActions are NOT stored on DB
	INSTANT_ACTION_TYPE = "instant_action"
)
