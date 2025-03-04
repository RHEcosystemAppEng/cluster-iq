package actions

type ActionTarget struct {
	AccountName string   `db:"accountName" json:"accountName"`
	Region      string   `db:"region" json:"region"`
	ClusterID   string   `db:"clusterID" json:"clusterID"`
	Instances   []string `db:"instances" json:"instances"`
}

func NewActionTarget(accountName string, region string, clusterID string, instances []string) *ActionTarget {
	return &ActionTarget{
		AccountName: accountName,
		Region:      region,
		ClusterID:   clusterID,
		Instances:   instances,
	}
}

func (at *ActionTarget) GetAccountName() string {
	return at.AccountName
}

func (at *ActionTarget) GetRegion() string {
	return at.Region
}

func (at *ActionTarget) GetClusterID() string {
	return at.ClusterID
}

func (at *ActionTarget) GetInstances() []string {
	return at.Instances
}
