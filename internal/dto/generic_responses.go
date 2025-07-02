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
