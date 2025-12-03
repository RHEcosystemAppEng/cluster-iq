package models

// TODO: Check if there's a better place for this struct
type ListOptions struct {
	PageSize int
	Offset   int
	Filters  map[string]interface{}
}

type ProviderDetail struct {
	AccountCount int `json:"account_count"`
	ClusterCount int `json:"cluster_count"`
}
