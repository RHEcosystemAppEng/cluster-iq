package responsetypes

// GenericResponse is a simple structure used to return a textual message in API responses.
type GenericResponse struct {
	Message string `json:"message"`
} // @name GenericResponse

// PatchResponse represents a generic response for successful PATCH operations.
type PatchResponse struct {
	Count  int    `json:"count"`
	Status string `json:"status"`
} // @name PatchResponse

// PostResponse represents a generic response for successful POST operations.
type PostResponse struct {
	Count  int    `json:"count"`
	Status string `json:"status"`
} // @name PostResponse

// GenericErrorResponse provides a standardized structure for returning error messages.
type GenericErrorResponse struct {
	Message string `json:"message"`
} // @name GenericErrorResponse

// ListResponse defines a reusable response type for endpoints returning lists of items.
type ListResponse[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
} // @name ListResponse

// NewListResponse returns a ListResponse initialized with the given items and total count.
// Ensures that Items is never nil to avoid null arrays in JSON responses.
func NewListResponse[T any](items []T, total int) *ListResponse[T] {
	if items == nil {
		items = []T{}
	}
	return &ListResponse[T]{
		Count: total,
		Items: items,
	}
}
