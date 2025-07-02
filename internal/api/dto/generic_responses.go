package dto

type ListResponse[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
}

func NewListResponse[T any](items []T, total int) *ListResponse[T] {
	if items == nil {
		items = []T{}
	}
	return &ListResponse[T]{
		Count: total,
		Items: items,
	}
}

type GenericResponse struct {
	Message string `json:"message"`
}

// NewGenericResponse creates a new generic response with the given message.
func NewGenericResponse(message string) GenericResponse {
	return GenericResponse{Message: message}
}

// ErrorResponse represents a generic error response with a message and a code.
