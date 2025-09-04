package dto

// ListResponse represents a paginated list of items.
type ListResponse[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
}

// NewListResponse creates a new ListResponse.
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

// GenericErrorResponse represents a generic error response.
type GenericErrorResponse struct {
	Message string `json:"message"`
}

// NewGenericErrorResponse creates a new GenericErrorResponse.
func NewGenericErrorResponse(message string) *GenericErrorResponse {
	return &GenericErrorResponse{
		Message: message,
	}
}
