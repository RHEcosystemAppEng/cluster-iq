package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

const (
	APIBaseURL        = "http://localhost:8081/api/v1"
	APIHealthcheckURL = APIBaseURL + "/healthcheck"
	APIAccountsURL    = APIBaseURL + "/accounts"
)

type ClusterListResponse struct {
	Count    int                 `json:"count"`
	Clusters []inventory.Cluster `json:"clusters"`
}

type AccountListResponse struct {
	Count    int                 `json:"count"`
	Accounts []inventory.Account `json:"accounts"`
}

// waitForAPIReady tries to reach the API endpoint until it responds or times out
func waitForAPIReady(t *testing.T) {
	url := APIHealthcheckURL
	t.Helper()
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode < 500 {
			return
		}
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("API did not become ready at %s", url)
}
