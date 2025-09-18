package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	t.Run("TestGetInstance", func(t *testing.T) { testGetInstances(t) })
	t.Run("TestGetInstanceByID", func(t *testing.T) { testGetInstanceByID(t) })
	t.Run("TestGetInstanceWithTags", func(t *testing.T) { testGetInstancesWithTags(t) })
	t.Run("TestPostOneInstance", func(t *testing.T) { testPostOneInstance(t) })
	t.Run("TestPostOneInstanceWithTags", func(t *testing.T) { testPostOneInstanceWithTags(t) })
	t.Run("TestPostMultipleInstances", func(t *testing.T) { testPostMultipleInstances(t) })
	t.Run("TestPostMultipleInstancesWithTags", func(t *testing.T) { testPostMultipleInstancesWithTags(t) })
	t.Run("TestDeleteInstance", func(t *testing.T) { testDeleteInstance(t) })
	t.Run("TestPatchInstance", func(t *testing.T) { testPatchInstance(t) })
}

func testGetInstances(t *testing.T) {
	expectedCount := 12
	// Getting instances data
	resp, err := http.Get(APIInstancesURL)
	if err != nil {
		t.Fatalf("Failed to make GetInstances request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetInstances response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected %d Instances, got %d", expectedCount, response.Count)
	}
}

func testGetInstanceByID(t *testing.T) {
	expectedCount := 1
	instanceID := "id-3123456789X"

	// Getting instances data
	resp, err := http.Get(APIInstancesURL + "/" + instanceID)
	if err != nil {
		t.Fatalf("Failed to make GetInstanceByID request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetInstancesByID response body: %v", err)
	}

	// Comparing data
	if len(response.Instances) != expectedCount {
		t.Fatalf("Expected: %d, Have: %d", expectedCount, response.Count)
	}
}

func testGetInstancesWithTags(t *testing.T) {
	expectedCount := 2
	instanceID := "id-3123456789X"

	// Getting instances data
	resp, err := http.Get(APIInstancesURL + "/" + instanceID)
	if err != nil {
		t.Fatalf("Failed to make GetInstancesWithTags request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.InstanceDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetInstancesWithTags response body: %v", err)
	}

	// Comparing data
	if len(response.Instances[0].Tags) != expectedCount {
		t.Fatalf("Expected: %d, Have: %d", expectedCount, response.Count)
	}
	// TODO Add elements check
}

func postInstances(t *testing.T, instances dto.InstanceDTORequestList) *responsetypes.PostResponse {
	jsonData, err := json.Marshal(instances)
	if err != nil {
		return nil
	}

	// Posting test data
	resp, err := http.Post(APIInstancesURL, "application/json", bytes.NewBuffer(jsonData))
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

	var response responsetypes.PostResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal("Can't Unmarshal API Response")
	}

	return &response
}

func testPostOneInstance(t *testing.T) {
	payload := dto.InstanceDTORequestList{
		Instances: []dto.InstanceDTORequest{
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
		},
	}

	// Posting test data
	response := postInstances(t, payload)

	// Checks
	if response.Count != 1 {
		t.Fatalf("Expected 1 Posted Instance, got %d", response.Count)
	}

	if response.Status != "Instance(s) Post OK" {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}

func testPostOneInstanceWithTags(t *testing.T) {
	payload := dto.InstanceDTORequestList{
		Instances: []dto.InstanceDTORequest{
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
				Tags: []dto.TagDTORequest{
					{
						Key:   "key-test-A",
						Value: "false",
					},
					{
						Key:   "key-test-B",
						Value: "true",
					},
				},
			},
		},
	}

	// Posting test data
	response := postInstances(t, payload)

	// Checks
	if response.Count != 1 {
		t.Fatalf("Expected 1 Posted Instance, got %d", response.Count)
	}

	if response.Status != "Instance(s) Post OK" {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}

func testPostMultipleInstances(t *testing.T) {
	payload := dto.InstanceDTORequestList{
		Instances: []dto.InstanceDTORequest{
			{
				InstanceID:       "id-instance-03",
				InstanceName:     "name-instance-03",
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
				InstanceID:       "id-instance-04",
				InstanceName:     "name-instance-04",
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
		},
	}

	// Posting test data
	response := postInstances(t, payload)

	if response.Count != 2 {
		t.Fatalf("Expected 2 Posted Instance, got %d", response.Count)
	}

	if response.Status != "Instance(s) Post OK" {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}

func testPostMultipleInstancesWithTags(t *testing.T) {
	payload := dto.InstanceDTORequestList{
		Instances: []dto.InstanceDTORequest{
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
		},
	}

	// Posting test data
	response := postInstances(t, payload)

	if response.Count != 2 {
		t.Fatalf("Expected 2 Posted Instance, got %d", response.Count)
	}

	if response.Status != "Instance(s) Post OK" {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}

func testPatchInstance(t *testing.T) {
	patchInstance := dto.InstanceDTORequest{
		InstanceID:   "ACC-001",
		InstanceName: "test-instance-003",
		Provider:     inventory.AWSProvider,
		LastScanTS:   time.Now(),
	}

	patchBody, err := json.Marshal(patchInstance)
	if err != nil {
		t.Fatalf("Failed to marshal updated instance: %v", err)
	}

	// Preparing PATCH request
	req, err := http.NewRequest(http.MethodPatch, APIInstancesURL+"/ACC-003", bytes.NewBuffer(patchBody))
	if err != nil {
		t.Fatalf("Failed to create PatchInstance request: %v", err)
	}

	// Executing PATCH request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute PatchInstance request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusNotImplemented {
		t.Fatalf("Expected status 501, got %d", resp.StatusCode)
	}
}

func testDeleteInstance(t *testing.T) {
	// Preparing DELETE request
	instanceID := "id-instance-03"
	req, err := http.NewRequest(http.MethodDelete, APIInstancesURL+"/"+instanceID, nil)
	if err != nil {
		t.Fatalf("Failed to create DeleteInstance request: %v", err)
	}

	// Executing DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute DeleteInstance request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var response responsetypes.DeleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode DeleteInstance response body: %v", err)
	}

	if response.Count != 1 {
		t.Fatalf("Expected 1 Deleted Instance, got %d", response.Count)
	}

	if response.Status != fmt.Sprintf("Instance '%s' Delete OK", instanceID) {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}
