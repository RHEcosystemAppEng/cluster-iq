package repositories

type ListOptions struct {
	PageSize int
	Offset   int
	Filters  map[string]interface{}
}
