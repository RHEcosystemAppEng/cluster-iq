package dto

import "time"

// ActionRunDTORequest represents a request to create or update an action run.
type ActionRunDTORequest struct {
	ScheduleID string `json:"scheduleId" binding:"required"`
	Status     string `json:"status"`
	ErrorMsg   string `json:"errorMsg"`
} // @name ActionRunRequest

// ActionRunDTOResponse represents the response for an action execution record.
type ActionRunDTOResponse struct {
	ID         string    `json:"id"`
	ScheduleID string    `json:"scheduleId"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
	Status     string    `json:"status"`
	ErrorMsg   string    `json:"errorMsg"`
} // @name ActionRunResponse
