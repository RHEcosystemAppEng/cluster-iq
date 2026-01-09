package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

const (
	APIClustersURL = APIBaseURL + "/clusters"
)

func TestClusters(t *testing.T) {
	waitForAPIReady(t)

	if err := refreshInventory(); err != nil {
		t.Fatal("Error refreshing inventory")
	}

	t.Run("Test List Clusters", func(t *testing.T) { testListClusters(t) })
	t.Run("Test List Clusters With Pagination", func(t *testing.T) { testListClustersWithPagination(t) })
	t.Run("Test List Clusters By Status", func(t *testing.T) { testListClustersByStatus(t) })
	t.Run("Test List Clusters By Provider", func(t *testing.T) { testListClustersByProvider(t) })
	t.Run("Test List Clusters By Region", func(t *testing.T) { testListClustersByRegion(t) })
	t.Run("Test List Clusters By Account", func(t *testing.T) { testListClustersByAccount(t) })
	t.Run("Test List Clusters Multiple Filters", func(t *testing.T) { testListClustersMultipleFilters(t) })
	t.Run("Test List Clusters Wrong Filters", func(t *testing.T) { testListClustersWrongFilters(t) })
	t.Run("Test Get Cluster By ID Success", func(t *testing.T) { testGetClusterByID_Exists(t) })
	t.Run("Test Get Cluster By ID Not Found", func(t *testing.T) { testGetClusterByID_NoExists(t) })
	t.Run("Test Get Cluster Instances Success", func(t *testing.T) { testGetClusterInstances_Exists(t) })
	t.Run("Test Get Cluster Instances Not Found", func(t *testing.T) { testGetClusterInstances_NoExists(t) })
	t.Run("Test Post One Cluster", func(t *testing.T) { testPostOneCluster(t) })
	t.Run("Test Post Multiple Clusters", func(t *testing.T) { testPostMultipleClusters(t) })
	t.Run("Test Post Wrong Cluster", func(t *testing.T) { testPostWrongCluster(t) })
	t.Run("Test Patch Cluster", func(t *testing.T) { testPatchCluster(t) })
	t.Run("Test Delete Cluster Success", func(t *testing.T) { testDeleteCluster_Exists(t) })
	t.Run("Test Delete Cluster Not Found", func(t *testing.T) { testDeleteCluster_NoExists(t) })
}

func testListClusters(t *testing.T) {
	expectedCount := 6
	expectedHTTPCode := http.StatusOK

	// Listing Clusters
	resp, err := http.Get(APIClustersURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ClusterDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListClustersWithPagination(t *testing.T) {
	expectedCount := 4
	expectedHTTPCode := http.StatusOK

	// Listing Clusters
	resp, err := http.Get(APIClustersURL + "?page=1&page_size=4")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ClusterDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListClustersByStatus(t *testing.T) {
	expectedCount := 0
	expectedHTTPCode := http.StatusOK

	// Listing Clusters
	resp, err := http.Get(APIClustersURL + "?status=Running")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ClusterDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListClustersByProvider(t *testing.T) {
	expectedCount := 2
	expectedHTTPCode := http.StatusOK

	// Listing Clusters
	resp, err := http.Get(APIClustersURL + "?provider=AWS")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ClusterDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListClustersByRegion(t *testing.T) {
	expectedCount := 1
	expectedHTTPCode := http.StatusOK

	// Listing Clusters
	resp, err := http.Get(APIClustersURL + "?region=us-east-1")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ClusterDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListClustersByAccount(t *testing.T) {
	expectedCount := 2
	expectedHTTPCode := http.StatusOK

	// Listing Clusters
	resp, err := http.Get(APIClustersURL + "?account=gcp-project-demo")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ClusterDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListClustersMultipleFilters(t *testing.T) {
	expectedCount := 0
	expectedHTTPCode := http.StatusOK

	// Listing Clusters
	resp, err := http.Get(APIClustersURL + "?status=Running&provider=GCP&account=gcp-project-demo")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ClusterDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListClustersWrongFilters(t *testing.T) {
	expectedMsg := "Failed to retrieve clusters"
	expectedHTTPCode := http.StatusInternalServerError

	// Listing Clusters
	resp, err := http.Get(APIClustersURL + "?provider=ANYONE")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func testGetClusterByID_Exists(t *testing.T) {
	clusterID := "aws-cluster-1-aws-infra-1"
	expectedHTTPCode := http.StatusOK

	// Getting Cluster by ID
	resp, err := http.Get(APIClustersURL + "/" + clusterID)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.ClusterDTOResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if clusterID != response.ClusterID {
		t.Fatalf("Expected ClusterID: '%s', got: '%s'", clusterID, response.ClusterID)
	}
}

func testGetClusterByID_NoExists(t *testing.T) {
	expectedMsg := "Cluster not found"
	expectedHTTPCode := http.StatusNotFound

	// Getting Cluster by ID
	resp, err := http.Get(APIClustersURL + "/" + "missing-cluster")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func testGetClusterInstances_Exists(t *testing.T) {
	expectedCount := 2
	clusterID := "gcp-cluster-1-gcp-infra-1"
	expectedHTTPCode := http.StatusOK

	// Getting Cluster instances list
	resp, err := http.Get(APIClustersURL + "/" + clusterID + "/instances")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response []dto.InstanceDTOResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	count := len(response)
	if count != expectedCount {
		t.Fatalf("Expected instances count: %d, Have: %d", expectedCount, count)
	}
}

func testGetClusterInstances_NoExists(t *testing.T) {
	expectedMsg := "Cluster not found"
	expectedHTTPCode := http.StatusNotFound

	// Getting Cluster instances list
	resp, err := http.Get(APIClustersURL + "/missing-cluster" + "/instances")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func postClusters(t *testing.T, clusters []dto.ClusterDTORequest, expectedHTTPCode int) *http.Response {
	b, err := json.Marshal(clusters)
	if err != nil {
		t.Fatal(err.Error())
		return nil
	}

	// Posting test data
	resp, err := http.Post(APIClustersURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	return resp
}

func testPostOneCluster(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedMsg := "OK"
	expectedCount := 1

	ts, _ := time.Parse(time.RFC3339, "2025-08-02T10:00:00+00:00")
	payload := []dto.ClusterDTORequest{
		{
			ClusterID:         "test-cluster-infra-1234",
			ClusterName:       "test-cluster",
			InfraID:           "infra-1234",
			Provider:          "AWS",
			Status:            "Running",
			Region:            "us-west-2",
			AccountID:         "111111111111",
			ConsoleLink:       "http://test-cluster.domain",
			LastScanTimestamp: ts,
			CreatedAt:         ts,
			Age:               1,
			Owner:             "John Doe",
		},
	}

	// Posting test data
	resp := postClusters(t, payload, expectedHTTPCode)
	defer resp.Body.Close()

	// Decode the JSON response
	var response responsetypes.PostResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Checks
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got '%d'", expectedCount, response.Count)
	}

	if response.Status != expectedMsg {
		t.Fatalf("Expected Status: '%s', got: '%s'", expectedMsg, response.Status)
	}
}

func testPostMultipleClusters(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedMsg := "OK"
	expectedCount := 2

	ts, _ := time.Parse(time.RFC3339, "2025-08-02T10:00:00+00:00")
	payload := []dto.ClusterDTORequest{
		{
			ClusterID:         "test-cluster-infra-2345",
			ClusterName:       "test-cluster",
			InfraID:           "infra-1234",
			Provider:          "Azure",
			Status:            "Running",
			Region:            "us-west-2",
			AccountID:         "subs-00000001",
			ConsoleLink:       "http://test-cluster.domain",
			LastScanTimestamp: ts,
			CreatedAt:         ts,
			Age:               1,
			Owner:             "John Doe",
		},
		{
			ClusterID:         "test-cluster-infra-3456",
			ClusterName:       "test-cluster",
			InfraID:           "infra-1234",
			Provider:          "GCP",
			Status:            "Running",
			Region:            "us-west-2",
			AccountID:         "gcp-project-1",
			ConsoleLink:       "http://test-cluster.domain",
			LastScanTimestamp: ts,
			CreatedAt:         ts,
			Age:               1,
			Owner:             "John Doe",
		},
	}

	// Posting test data
	resp := postClusters(t, payload, expectedHTTPCode)
	defer resp.Body.Close()

	// Decode the JSON response
	var response responsetypes.PostResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Checks
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got '%d'", expectedCount, response.Count)
	}

	if response.Status != expectedMsg {
		t.Fatalf("Expected Status: '%s', got: '%s'", expectedMsg, response.Status)
	}
}

func testPostWrongCluster(t *testing.T) {
	expectedHTTPCode := http.StatusInternalServerError
	expectedMsg := "Failed to create clusters: named-exec INSERT error: pq: invalid input value for enum cloud_provider: \"Provider\""

	ts, _ := time.Parse(time.RFC3339, "2025-08-02T10:00:00+00:00")
	payload := []dto.ClusterDTORequest{
		{
			ClusterID:         "test-cluster-infra-1234",
			ClusterName:       "test-cluster",
			InfraID:           "infra-1234",
			Provider:          "Provider",
			Status:            "Running",
			Region:            "us-west-2",
			AccountID:         "111111111111",
			ConsoleLink:       "http://test-cluster.domain",
			LastScanTimestamp: ts,
			CreatedAt:         ts,
			Age:               1,
			Owner:             "John Doe",
		},
	}

	// Posting test data
	resp := postClusters(t, payload, expectedHTTPCode)
	defer resp.Body.Close()

	// Decode the JSON response
	var response responsetypes.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Checks
	if response.Message != expectedMsg {
		t.Fatalf("Expected Status: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func testPatchCluster(t *testing.T) {
	expectedHTTPCode := http.StatusNotImplemented

	patchCluster := dto.ClusterDTORequest{
		ClusterID:         "test-cluster-infra-2345",
		ClusterName:       "test-Cluster-003",
		Provider:          inventory.AWSProvider,
		LastScanTimestamp: time.Now(),
	}

	patchBody, err := json.Marshal(patchCluster)
	if err != nil {
		t.Fatalf("Failed to marshal updated Cluster: %v", err)
	}

	// Preparing PATCH request
	req, err := http.NewRequest(http.MethodPatch, APIClustersURL+"/aws-cluster-1", bytes.NewBuffer(patchBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Executing PATCH request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)
}

func testDeleteCluster_Exists(t *testing.T) {
	expectedHTTPCode := http.StatusNoContent

	// Preparing DELETE request
	clusterID := "test-cluster-infra-1234"
	req, err := http.NewRequest(http.MethodDelete, APIClustersURL+"/"+clusterID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Executing DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)
}

func testDeleteCluster_NoExists(t *testing.T) {
	expectedHTTPCode := http.StatusNoContent

	// Preparing DELETE request
	clusterID := "missing"
	req, err := http.NewRequest(http.MethodDelete, APIClustersURL+"/"+clusterID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Executing DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)
}
