// test/integration/post_account_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
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

// loadTestData loads test data files from a JSON file
func loadTestData(filename string) []byte {
	path := filepath.Join("./data_test_files", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
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

func TestPostAccount(t *testing.T) {
	waitForAPIReady(t)

	// Loading test data
	filename := "test_accounts_01.json"
	payload := loadTestData(filename)
	if payload == nil {
		t.Fatalf("Error when loading %s test data file", filename)
	}

	// Posting test data
	resp, err := http.Post(APIAccountsURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetAccount(t *testing.T) {
	waitForAPIReady(t)

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response AccountListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GET response body: %v", err)
	}

	// Loading and decoding test file to compare
	testData := loadTestData("test_accounts_01.json")
	var accounts []inventory.Account
	if err := json.Unmarshal(testData, &accounts); err != nil {
		t.Fatalf("failed to unmarshal accounts JSON: %v", err)
	}

	// Comparing data
	if response.Count != len(accounts) {
		t.Fatalf("Accounts arrays doesn't have same number of elements. Expected: %d, Have: %d", len(accounts), response.Count)
	}
	if !reflect.DeepEqual(accounts, response.Accounts) {
		t.Errorf("Accounts arrays are not equal")
	}
}

func TestGetAccountByID(t *testing.T) {
	waitForAPIReady(t)

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL + "/test-account-001")
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response AccountListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GET response body: %v", err)
	}

	// Comparing data
	if len(response.Accounts) != 1 {
		t.Fatalf("Accounts arrays doesn't have same number of elements. Expected: %d, Have: %d", 1, response.Count)
	}
}

func TestDeleteAccount(t *testing.T) {
	waitForAPIReady(t)

	// Preparing DELETE request
	req, err := http.NewRequest(http.MethodDelete, APIAccountsURL+"/ACC-002", nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}

	// Executing DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute DELETE request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestPatchAccount(t *testing.T) {
	waitForAPIReady(t)

	patchAccount := inventory.Account{
		ID:           "ACC-003",
		Name:         "test-account-003",
		Provider:     "AWS",
		ClusterCount: 4,
		TotalCost:    50.67,
	}

	patchBody, err := json.Marshal(patchAccount)
	if err != nil {
		t.Fatalf("Failed to marshal updated account: %v", err)
	}

	// Preparing PATCH request
	req, err := http.NewRequest(http.MethodPatch, APIAccountsURL+"/ACC-003", bytes.NewBuffer(patchBody))
	if err != nil {
		t.Fatalf("Failed to create PATCH request: %v", err)
	}

	// Executing PATCH request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute PATCH request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusNotImplemented {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetAccountClusters(t *testing.T) {
	waitForAPIReady(t)

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL + "/test-account-001/clusters")
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response ClusterListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GET response body: %v", err)
	}

	// Comparing data
	if len(response.Clusters) != 0 {
		t.Fatalf("Accounts arrays doesn't have same number of elements. Expected: %d, Have: %d", 0, response.Count)
	}
}
