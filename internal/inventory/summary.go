package inventory

// ClustersSummary provides a summary of cluster counts by status.
type ClustersSummary struct {
	Running  int `db:"running"`
	Stopped  int `db:"stopped"`
	Archived int `db:"archived"`
}

// InstancesSummary provides a summary of instance counts.
type InstancesSummary struct {
	Count int `db:"count"`
}
