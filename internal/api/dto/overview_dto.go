package dto

type OverviewSummary struct {
	Clusters  ClustersSummary  `json:"clusters"`
	Instances InstancesSummary `json:"instances"`
	Providers ProvidersSummary `json:"providers"`
}

type ClustersSummary struct {
	Running  int `json:"running"`
	Stopped  int `json:"stopped"`
	Archived int `json:"archived"`
}

type InstancesSummary struct {
	Count int `json:"count"`
}

type ProvidersSummary struct {
	AWS   ProviderDetail `json:"aws"`
	GCP   ProviderDetail `json:"gcp"`
	Azure ProviderDetail `json:"azure"`
}

type ProviderDetail struct {
	AccountCount int `json:"account_count"`
	ClusterCount int `json:"cluster_count"`
}
