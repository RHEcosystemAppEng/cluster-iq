package actions

type ActionStatus string

const (
	StatusPending ActionStatus = "Pending"
	StatusRunning ActionStatus = "Running"
	StatusFailed  ActionStatus = "Failed"
	StatusSuccess ActionStatus = "Success"
	StatusUnknown ActionStatus = "Unknown"
)
