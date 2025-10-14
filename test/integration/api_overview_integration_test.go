package integration

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

const (
	APIOverviewURL = APIBaseURL + "/overview"
)

func TestAPIOverview(t *testing.T) {
	waitForAPIReady(t)

	if err := refreshInventory(); err != nil {
		t.Fatal("Error refreshing inventory")
	}

	t.Run("Test Get Overview", func(t *testing.T) { testGetOverview(t) })
}

func testGetOverview(t *testing.T) {
	expectedHTTPCode := http.StatusOK
	lastScanTS, _ := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
	expectedOverviewResponse := dto.OverviewSummary{
		Clusters: dto.ClusterSummary{
			Running:  0,
			Stopped:  0,
			Archived: 6,
		},
		Instances: dto.InstancesSummary{
			Running:  10,
			Stopped:  3,
			Archived: 0,
		},
		Providers: dto.ProvidersSummary{
			AWS: dto.ProviderDetails{
				AccountCount: 1,
				ClusterCount: 2,
			},
			GCP: dto.ProviderDetails{
				AccountCount: 2,
				ClusterCount: 2,
			},
			Azure: dto.ProviderDetails{
				AccountCount: 2,
				ClusterCount: 2,
			},
		},
		Scanner: dto.Scanner{
			LastScanTimestamp: lastScanTS,
		},
	}

	// Getting accounts data
	resp, err := http.Get(APIOverviewURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.OverviewSummary
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if !reflect.DeepEqual(response, expectedOverviewResponse) {
		t.Fatalf("Expected Overview: '%+v', got: '%+v'", expectedOverviewResponse, response)
	}
}
