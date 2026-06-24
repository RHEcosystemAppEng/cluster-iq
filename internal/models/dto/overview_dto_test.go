package dto

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// TestToOverviewSummaryDTO verifies ToOverviewSummaryDTO converts inventory.OverviewSummary to DTO.
func TestToOverviewSummaryDTO(t *testing.T) {
	t.Run("Convert OverviewSummary", func(t *testing.T) { testToOverviewSummaryDTO_Correct(t) })
}

func testToOverviewSummaryDTO_Correct(t *testing.T) {
	now := time.Now().UTC()

	model := inventory.OverviewSummary{
		Clusters: inventory.ClustersSummary{
			Running:  1,
			Stopped:  2,
			Archived: 3,
		},
		Instances: inventory.InstancesSummary{
			Running:  10,
			Stopped:  20,
			Archived: 30,
		},
		Providers: inventory.ProvidersSummary{
			AWS:   inventory.ProviderDetails{AccountCount: 1, ClusterCount: 2},
			GCP:   inventory.ProviderDetails{AccountCount: 3, ClusterCount: 4},
			Azure: inventory.ProviderDetails{AccountCount: 5, ClusterCount: 6},
		},
		Scanner: inventory.Scanner{
			LastScanTimestamp: now,
		},
		TopRegions:        []inventory.TopItem{{Name: "us-east-1", ClusterCount: 8}},
		TopOwners:         []inventory.TopItem{{Name: "jsmith", ClusterCount: 5}},
		ClustersByPartner: []inventory.TopItem{{Name: "Acme", ClusterCount: 3}},
		CostPerAccount:    []inventory.AccountCost{{AccountName: "my-account", CurrentMonthCost: 1234.56}},
	}

	dto := ToOverviewSummaryDTO(model)

	assert.Equal(t, 1, dto.Clusters.Running)
	assert.Equal(t, 2, dto.Clusters.Stopped)
	assert.Equal(t, 3, dto.Clusters.Archived)

	assert.Equal(t, 10, dto.Instances.Running)
	assert.Equal(t, 20, dto.Instances.Stopped)
	assert.Equal(t, 30, dto.Instances.Archived)

	assert.Equal(t, 1, dto.Providers.AWS.AccountCount)
	assert.Equal(t, 2, dto.Providers.AWS.ClusterCount)
	assert.Equal(t, 3, dto.Providers.GCP.AccountCount)
	assert.Equal(t, 4, dto.Providers.GCP.ClusterCount)
	assert.Equal(t, 5, dto.Providers.Azure.AccountCount)
	assert.Equal(t, 6, dto.Providers.Azure.ClusterCount)

	assert.Equal(t, now, dto.Scanner.LastScanTimestamp)

	assert.Len(t, dto.TopRegions, 1)
	assert.Equal(t, "us-east-1", dto.TopRegions[0].Name)
	assert.Equal(t, 8, dto.TopRegions[0].ClusterCount)

	assert.Len(t, dto.TopOwners, 1)
	assert.Equal(t, "jsmith", dto.TopOwners[0].Name)

	assert.Len(t, dto.ClustersByPartner, 1)
	assert.Equal(t, "Acme", dto.ClustersByPartner[0].Name)

	assert.Len(t, dto.CostPerAccount, 1)
	assert.Equal(t, "my-account", dto.CostPerAccount[0].AccountName)
	assert.InDelta(t, 1234.56, dto.CostPerAccount[0].CurrentMonthCost, 0.01)
}

// TestToClusterSummaryDTO verifies toClusterSummaryDTO conversion.
func TestToClusterSummaryDTO(t *testing.T) {
	t.Run("Convert ClusterSummary", func(t *testing.T) { testToClusterSummaryDTO_Correct(t) })
}

func testToClusterSummaryDTO_Correct(t *testing.T) {
	model := inventory.ClustersSummary{Running: 1, Stopped: 2, Archived: 3}
	dto := toClusterSummaryDTO(model)

	assert.Equal(t, 1, dto.Running)
	assert.Equal(t, 2, dto.Stopped)
	assert.Equal(t, 3, dto.Archived)
}

// TestToInstancesSummaryDTO verifies toInstancesSummaryDTO conversion.
func TestToInstancesSummaryDTO(t *testing.T) {
	t.Run("Convert InstancesSummary", func(t *testing.T) { testToInstancesSummaryDTO_Correct(t) })
}

func testToInstancesSummaryDTO_Correct(t *testing.T) {
	model := inventory.InstancesSummary{Running: 10, Stopped: 20, Archived: 30}
	dto := toInstancesSummaryDTO(model)

	assert.Equal(t, 10, dto.Running)
	assert.Equal(t, 20, dto.Stopped)
	assert.Equal(t, 30, dto.Archived)
}

// TestToProvidersSummaryDTO verifies toProvidersSummaryDTO conversion.
func TestToProvidersSummaryDTO(t *testing.T) {
	t.Run("Convert ProvidersSummary", func(t *testing.T) { testToProvidersSummaryDTO_Correct(t) })
}

func testToProvidersSummaryDTO_Correct(t *testing.T) {
	model := inventory.ProvidersSummary{
		AWS:   inventory.ProviderDetails{AccountCount: 1, ClusterCount: 2},
		GCP:   inventory.ProviderDetails{AccountCount: 3, ClusterCount: 4},
		Azure: inventory.ProviderDetails{AccountCount: 5, ClusterCount: 6},
	}

	dto := toProvidersSummaryDTO(model)

	assert.Equal(t, 1, dto.AWS.AccountCount)
	assert.Equal(t, 2, dto.AWS.ClusterCount)
	assert.Equal(t, 3, dto.GCP.AccountCount)
	assert.Equal(t, 4, dto.GCP.ClusterCount)
	assert.Equal(t, 5, dto.Azure.AccountCount)
	assert.Equal(t, 6, dto.Azure.ClusterCount)
}

// TestToProviderDetailsDTO verifies toProviderDetailsDTO conversion.
func TestToProviderDetailsDTO(t *testing.T) {
	t.Run("Convert ProviderDetails", func(t *testing.T) { testToProviderDetailsDTO_Correct(t) })
}

func testToProviderDetailsDTO_Correct(t *testing.T) {
	model := inventory.ProviderDetails{AccountCount: 7, ClusterCount: 8}
	dto := toProviderDetailsDTO(model)

	assert.Equal(t, 7, dto.AccountCount)
	assert.Equal(t, 8, dto.ClusterCount)
}

// TestToScannerDTO verifies toScannerDTO conversion.
func TestToScannerDTO(t *testing.T) {
	t.Run("Convert Scanner", func(t *testing.T) { testToScannerDTO_Correct(t) })
}

func testToScannerDTO_Correct(t *testing.T) {
	now := time.Now().UTC()
	model := inventory.Scanner{LastScanTimestamp: now}
	dto := toScannerDTO(model)

	assert.Equal(t, now, dto.LastScanTimestamp)
}

// TestToTopItemsDTO verifies toTopItemsDTO conversion.
func TestToTopItemsDTO(t *testing.T) {
	t.Run("Convert TopItems", func(t *testing.T) {
		items := []inventory.TopItem{
			{Name: "us-east-1", ClusterCount: 8},
			{Name: "eu-west-1", ClusterCount: 3},
		}
		result := toTopItemsDTO(items)

		assert.Len(t, result, 2)
		assert.Equal(t, "us-east-1", result[0].Name)
		assert.Equal(t, 8, result[0].ClusterCount)
		assert.Equal(t, "eu-west-1", result[1].Name)
		assert.Equal(t, 3, result[1].ClusterCount)
	})

	t.Run("Convert empty TopItems", func(t *testing.T) {
		result := toTopItemsDTO([]inventory.TopItem{})
		assert.Len(t, result, 0)
	})
}

// TestToAccountCostsDTO verifies toAccountCostsDTO conversion.
func TestToAccountCostsDTO(t *testing.T) {
	t.Run("Convert AccountCosts", func(t *testing.T) {
		costs := []inventory.AccountCost{
			{AccountName: "prod", CurrentMonthCost: 5000.50},
			{AccountName: "dev", CurrentMonthCost: 1200.00},
		}
		result := toAccountCostsDTO(costs)

		assert.Len(t, result, 2)
		assert.Equal(t, "prod", result[0].AccountName)
		assert.InDelta(t, 5000.50, result[0].CurrentMonthCost, 0.01)
		assert.Equal(t, "dev", result[1].AccountName)
		assert.InDelta(t, 1200.00, result[1].CurrentMonthCost, 0.01)
	})

	t.Run("Convert empty AccountCosts", func(t *testing.T) {
		result := toAccountCostsDTO([]inventory.AccountCost{})
		assert.Len(t, result, 0)
	})
}
