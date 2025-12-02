package actions

type ActionStatus string

const (
	StatusPending   ActionStatus = "Pending"
	StatusRunning   ActionStatus = "Running"
	StatusFailed    ActionStatus = "Failed"
	StatusCompleted ActionStatus = "Completed"
)
