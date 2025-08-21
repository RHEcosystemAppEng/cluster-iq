package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/apiresponsetypes"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

const (
	APIClustersURL = APIBaseURL + "/clusters"
)

func TestClusterAccounts(t *testing.T) {
	waitForAPIReady(t)

	t.Run("TestGetCluster", func(t *testing.T) { testGetClusters(t) })
	t.Run("TestGetClusterByID", func(t *testing.T) { testGetClusterByID(t) })
	t.Run("TestGetClusterInstances", func(t *testing.T) { testGetClusterInstances(t) })
	t.Run("TestPostOneCluster", func(t *testing.T) { testPostOneCluster(t) })
	t.Run("TestPostMultipleClusters", func(t *testing.T) { testPostMultipleClusters(t) })
	t.Run("TestDeleteCluster", func(t *testing.T) { testDeleteCluster(t) })
	t.Run("TestPatchCluster", func(t *testing.T) { testPatchCluster(t) })
}

func testGetClusters(t *testing.T) {
	expectedCount := 6
	// Getting Clusters data
	resp, err := http.Get(APIClustersURL)
	if err != nil {
		t.Fatalf("Failed to make GetClusters request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.ClusterDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetClusters response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected %d Clusters, got %d", expectedCount, response.Count)
	}
}

func testGetClusterByID(t *testing.T) {
	expectedCount := 1
	clusterID := "aws-cluster-1-aws-infra-1"

	// Getting Clusters data
	resp, err := http.Get(APIClustersURL + "/" + clusterID)
	if err != nil {
		t.Fatalf("Failed to make GetClusterByID request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.ClusterDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetClusterByID response body: %v", err)
	}

	// Comparing data
	if len(response.Clusters) != expectedCount {
		t.Fatalf("Expected: %d, Have: %d", expectedCount, response.Count)
	}
}

func testGetClusterInstances(t *testing.T) {
	expectedCount := 2
	clusterID := "gcp-cluster-1-gcp-infra-1"

	// Getting Clusters data
	resp, err := http.Get(APIClustersURL + "/" + clusterID + "/instances")
	if err != nil {
		t.Fatalf("Failed to make GetClusterInstances request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetClusterInstances response body: %v", err)
	}

	// Comparing data
	if len(response.Instances) != expectedCount {
		t.Fatalf("Expected: %d, Have: %d", expectedCount, response.Count)
	}

	// TODO Add elements check
}

func postClusters(t *testing.T, accounts string) *apiresponsetypes.PostResponse {
	// Posting test data
	resp, err := http.Post(APIClustersURL, "application/json", bytes.NewBuffer([]byte(accounts)))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Cant Read API Response Body")
	}

	var response apiresponsetypes.PostResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal("Can't Unmarshal API Response")
	}

	return &response
}

func testPostOneCluster(t *testing.T) {
	payload := `
	{
		"clusters": [
			{
				"clusterID": "test-cluster-infra-1234",
				"clusterName": "test-cluster",
				"infraID": "infra-1234",
				"provider": "AWS",
				"status": "Running",
				"region": "us-west-2",
				"accountID": "111111111111",
				"consoleLink": "http://test-cluster.domain",
				"lastScanTS": "2025-08-02 10:00:00+00",
				"createdAt": "2025-08-01 10:00:00+00",
				"age": 1,
				"owner": "John Doe"
			}
		]
	}
	`

	// Loading test data
	response := postClusters(t, payload)

	// Checks
	if response.Count != 1 {
		t.Fatalf("Expected 1 Posted Cluster, got %d", response.Count)
	}

	if response.Status != "Cluster(s) Post OK" {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}

func testPostMultipleClusters(t *testing.T) {
	payload := `
	{
		"clusters": [
			{
				"clusterID": "test-cluster-infra-2345",
				"clusterName": "test-cluster",
				"infraID": "infra-2345",
				"provider": "Azure",
				"status": "Stopped",
				"region": "westeurope",
				"accountID": "subs-00000001",
				"consoleLink": "http://test-cluster.domain",
				"lastScanTS": "2025-08-02 10:00:00+00",
				"createdAt": "2025-08-01 10:00:00+00",
				"age": 1,
				"owner": "John Doe"
			},
			{
				"clusterID": "test-cluster-infra-3456",
				"clusterName": "test-cluster",
				"infraID": "infra-3456",
				"provider": "GCP",
				"status": "Running",
				"region": "europe-west-2",
				"accountID": "gcp-project-1",
				"consoleLink": "http://test-cluster.domain",
				"lastScanTS": "2025-08-02 10:00:00+00",
				"createdAt": "2025-08-01 10:00:00+00",
				"age": 1,
				"owner": "John Doe"
			}
		]
	}

	`

	// Loading test data
	response := postClusters(t, payload)

	// Checks
	if response.Count != 2 {
		t.Fatalf("Expected 2 Posted Cluster, got %d", response.Count)
	}

	if response.Status != "Cluster(s) Post OK" {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}

func testPatchCluster(t *testing.T) {
	patchCluster := dto.ClusterDTORequest{
		ClusterID:   "test-cluster-infra-2345",
		ClusterName: "test-Cluster-003",
		Provider:    inventory.AWSProvider,
		LastScanTS:  time.Now(),
	}

	patchBody, err := json.Marshal(patchCluster)
	if err != nil {
		t.Fatalf("Failed to marshal updated Cluster: %v", err)
	}

	// Preparing PATCH request
	req, err := http.NewRequest(http.MethodPatch, APIClustersURL+"/aws-cluster-1", bytes.NewBuffer(patchBody))
	if err != nil {
		t.Fatalf("Failed to create PatchCluster request: %v", err)
	}

	// Executing PATCH request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute PatchCluster request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusNotImplemented {
		t.Fatalf("Expected status 501, got %d", resp.StatusCode)
	}
}

func testDeleteCluster(t *testing.T) {
	// Preparing DELETE request
	clusterID := "test-cluster-infra-1234"
	req, err := http.NewRequest(http.MethodDelete, APIClustersURL+"/"+clusterID, nil)
	if err != nil {
		t.Fatalf("Failed to create DeleteCluster request: %v", err)
	}

	// Executing DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute DeleteCluster request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var response apiresponsetypes.DeleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode DeleteCluster response body: %v", err)
	}

	if response.Count != 1 {
		t.Fatalf("Expected 1 Deleted Cluster, got %d", response.Count)
	}

	if response.Status != fmt.Sprintf("Cluster '%s' Delete OK", clusterID) {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}
