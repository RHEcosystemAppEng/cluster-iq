package dto

// GenericErrorResponse represents a generic error response returned by the API.
//
// This structure is used to provide a consistent error message format in the API responses.
// It includes a single field, `Message`, that contains a descriptive error message.
type GenericErrorResponse struct {
	Message string `json:"message"`
}

// NewGenericErrorResponse creates a new instance of GenericErrorResponse.
//
// This function is a utility for initializing a GenericErrorResponse with a specified error message.
//
// Parameters:
// - message: The error message to include in the response.
//
// Returns:
// - A pointer to a new GenericErrorResponse instance containing the provided message.
func NewGenericErrorResponse(message string) *GenericErrorResponse {
	return &GenericErrorResponse{
		Message: message,
	}
}
