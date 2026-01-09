package models

// ListOptions is a struct that contains the options for a list operation.
type ListOptions struct {
	PageSize int
	Offset   int
	Filters  map[string]interface{}
}
