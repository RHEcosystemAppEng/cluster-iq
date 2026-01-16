package dto

// PaginationRequest represents pagination parameters for list endpoints,
// allowing clients to specify page number and page size.
type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"gte=1"`
	PageSize int `form:"page_size,default=10" binding:"gte=1,lte=100"`
} // @name PaginationRequest
