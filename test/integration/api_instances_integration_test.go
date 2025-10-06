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
	APIInstancesURL = APIBaseURL + "/instances"
)

func TestAPIInstances(t *testing.T) {
	waitForAPIReady(t)

	t.Run("Test List Instances", func(t *testing.T) { testListInstances(t) })
	t.Run("Test List Instances With Pagination", func(t *testing.T) { testListInstancesWithPagination(t) })
	t.Run("Test List Instances By ClusterID", func(t *testing.T) { testListInstancesByClusterID(t) })
	t.Run("Test List Instances By Status", func(t *testing.T) { testListInstancesByStatus(t) })
	t.Run("Test List Instances Multiple Filters", func(t *testing.T) { testListInstancesMultipleFilters(t) })
	t.Run("Test List Instances Wrong Filters", func(t *testing.T) { testListInstancesWrongFilters(t) })
	t.Run("Test Get Instance By ID Success", func(t *testing.T) { testGetInstanceByID_Exists(t) })
	t.Run("Test Get Instance By ID Not Found", func(t *testing.T) { testGetInstanceByID_NoExists(t) })
	t.Run("Test Get Instance With Tags Success", func(t *testing.T) { testGetInstancesWithTags_Exists(t) })
	t.Run("Test Get Instance With Tags Not Found", func(t *testing.T) { testGetInstancesWithTags_NoExists(t) })
	t.Run("Test Post Instances", func(t *testing.T) { testPostInstances(t) })
	t.Run("Test Post Instances With Tags", func(t *testing.T) { testPostInstancesWithTags(t) })
	t.Run("Test Post Instances Wrong values", func(t *testing.T) { testPostInstancesWrongValues(t) })
}

func testListInstances(t *testing.T) {
	expectedCount := 10
	expectedHTTPCode := http.StatusOK

	// Listing Instances
	resp, err := http.Get(APIInstancesURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListInstancesWithPagination(t *testing.T) {
	expectedCount := 8
	expectedHTTPCode := http.StatusOK

	// Listing Instances
	resp, err := http.Get(APIInstancesURL + "?page=1&page_size=8")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListInstancesByClusterID(t *testing.T) {
	expectedCount := 2
	expectedHTTPCode := http.StatusOK

	// Listing Instances
	resp, err := http.Get(APIInstancesURL + "?cluster_id=aws-cluster-1-aws-infra-1")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListInstancesByStatus(t *testing.T) {
	expectedCount := 6
	expectedHTTPCode := http.StatusOK

	// Listing Instances
	resp, err := http.Get(APIInstancesURL + "?status=Running")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListInstancesMultipleFilters(t *testing.T) {
	expectedCount := 1
	expectedHTTPCode := http.StatusOK

	// Listing Instances
	resp, err := http.Get(APIInstancesURL + "?cluster_id=aws-cluster-1-aws-infra-1&status=Running")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func testListInstancesWrongFilters(t *testing.T) {
	expectedMsg := "Failed to retrieve instances"
	expectedHTTPCode := http.StatusInternalServerError

	// Listing Instances
	resp, err := http.Get(APIInstancesURL + "?status=Stopping")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func testGetInstanceByID_Exists(t *testing.T) {
	instanceID := "id-3123456789X"
	expectedHTTPCode := http.StatusOK

	// Getting instances data
	resp, err := http.Get(APIInstancesURL + "/" + instanceID)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.InstanceDTOResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if instanceID != response.InstanceID {
		t.Fatalf("Expected InstanceID: '%s', got: '%s'", instanceID, response.InstanceID)
	}
}

func testGetInstanceByID_NoExists(t *testing.T) {
	expectedMsg := "Instance not found"
	expectedHTTPCode := http.StatusNotFound

	// Getting instances data
	resp, err := http.Get(APIInstancesURL + "/" + "missing-instance")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func testGetInstancesWithTags_Exists(t *testing.T) {
	instanceID := "id-3123456789X"
	expectedHTTPCode := http.StatusOK
	expectedTags := 2

	// Getting instances data
	resp, err := http.Get(APIInstancesURL + "/" + instanceID)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.InstanceDTOResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Comparing data
	if instanceID != response.InstanceID {
		t.Fatalf("Expected InstanceID: '%s', got: '%s'", instanceID, response.InstanceID)
	}

	if count := len(response.Tags); count != expectedTags {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedTags, count)
	}
}

func testGetInstancesWithTags_NoExists(t *testing.T) {
	expectedMsg := "Instance not found"
	expectedHTTPCode := http.StatusNotFound

	// Getting instances data
	resp, err := http.Get(APIInstancesURL + "/" + "missing-instance")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func postInstances(t *testing.T, instances []dto.InstanceDTORequest, expectedHTTPCode int) *http.Response {
	jsonData, err := json.Marshal(instances)
	if err != nil {
		t.Fatal(err.Error())
		return nil
	}

	// Posting test data
	resp, err := http.Post(APIInstancesURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	return resp
}

func testPostInstances(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedMsg := "OK"
	expectedCount := 2
	payload := []dto.InstanceDTORequest{
		{
			InstanceID:       "id-instance-01",
			InstanceName:     "name-instance-01",
			InstanceType:     "c.small",
			Provider:         "AWS",
			AvailabilityZone: "eu-west-1a",
			Status:           inventory.Running,
			ClusterID:        "aws-cluster-1-aws-infra-1",
			LastScanTS:       time.Now(),
			CreatedAt:        time.Now(),
			Age:              1,
			Owner:            "John Doe",
			Tags:             []dto.TagDTORequest{},
		},
		{
			InstanceID:       "id-instance-02",
			InstanceName:     "name-instance-0",
			InstanceType:     "c.small",
			Provider:         "AWS",
			AvailabilityZone: "eu-west-1a",
			Status:           inventory.Running,
			ClusterID:        "aws-cluster-1-aws-infra-1",
			LastScanTS:       time.Now(),
			CreatedAt:        time.Now(),
			Age:              1,
			Owner:            "John Doe",
			Tags:             []dto.TagDTORequest{},
		},
	}

	// Posting test data
	resp := postInstances(t, payload, expectedHTTPCode)
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

func testPostInstancesWithTags(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedMsg := "OK"
	expectedCount := 2
	payload := []dto.InstanceDTORequest{
		{
			InstanceID:       "id-instance-05",
			InstanceName:     "name-instance-05",
			InstanceType:     "c.small",
			Provider:         "AWS",
			AvailabilityZone: "eu-west-1a",
			Status:           inventory.Running,
			ClusterID:        "aws-cluster-1-aws-infra-1",
			LastScanTS:       time.Now(),
			CreatedAt:        time.Now(),
			Age:              1,
			Owner:            "John Doe",
			Tags: []dto.TagDTORequest{
				{
					Key:   "key-test-C",
					Value: "false",
				},
				{
					Key:   "key-test-D",
					Value: "true",
				},
			},
		},
		{
			InstanceID:       "id-instance-06",
			InstanceName:     "name-instance-06",
			InstanceType:     "c.small",
			Provider:         "AWS",
			AvailabilityZone: "eu-west-1a",
			Status:           inventory.Running,
			ClusterID:        "aws-cluster-1-aws-infra-1",
			LastScanTS:       time.Now(),
			CreatedAt:        time.Now(),
			Age:              1,
			Owner:            "John Doe",
			Tags: []dto.TagDTORequest{
				{
					Key:   "key-test-E",
					Value: "false",
				},
				{
					Key:   "key-test-F",
					Value: "true",
				},
			},
		},
	}

	// Posting test data
	resp := postInstances(t, payload, expectedHTTPCode)
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

func testPostInstancesWrongValues(t *testing.T) {
	expectedHTTPCode := http.StatusInternalServerError
	expectedMsg := "Failed to create instances: named-exec INSERT error: pq: invalid input value for enum cloud_provider: \"PROVIDER\""
	payload := []dto.InstanceDTORequest{
		{
			InstanceID:       "error-instance",
			InstanceName:     "error-instance-name",
			InstanceType:     "c.small",
			Provider:         "PROVIDER",
			AvailabilityZone: "eu-west-1a",
			Status:           inventory.Running,
			ClusterID:        "missing-cluster",
			LastScanTS:       time.Now(),
			CreatedAt:        time.Now(),
			Age:              1,
			Owner:            "John Doe",
			Tags: []dto.TagDTORequest{
				{
					Key:   "key-test-C",
					Value: "false",
				},
				{
					Key:   "key-test-D",
					Value: "true",
				},
			},
		},
	}

	// Posting test data
	resp := postInstances(t, payload, expectedHTTPCode)
	defer resp.Body.Close()

	// Decode the JSON response
	var response dto.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Checks
	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}
