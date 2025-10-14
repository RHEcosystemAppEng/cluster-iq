package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

const (
	APIEventsURL = APIBaseURL + "/events"
)

func TestAPIEvents(t *testing.T) {
	waitForAPIReady(t)

	if err := refreshInventory(); err != nil {
		t.Fatal("Error refreshing inventory")
	}

	t.Run("Test List System Events", func(t *testing.T) { testListSystemEvents(t) })
	t.Run("Test List Cluster Events", func(t *testing.T) { testListClusterEvents(t) })
	t.Run("Test Post Events", func(t *testing.T) { testPostEvents(t) })
	t.Run("Test Update Event", func(t *testing.T) { testUpdateEvent(t) })
}

func testListSystemEvents(t *testing.T) {
	expectedCount := 2
	expectedHTTPCode := http.StatusOK

	// Getting accounts data
	resp, err := http.Get(APIEventsURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.SystemEventDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}

	if len := len(response.Items); len != expectedCount {
		t.Fatalf("Expected Items: '%d', got: '%d'", expectedCount, len)
	}
}

func testListClusterEvents(t *testing.T) {
	expectedCount := 1
	clusterID := "aws-cluster-1-aws-infra-1"
	expectedHTTPCode := http.StatusOK

	// Getting accounts data
	resp, err := http.Get(APIClustersURL + "/" + clusterID + "/events")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.SystemEventDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}

	if len := len(response.Items); len != expectedCount {
		t.Fatalf("Expected Items: '%d', got: '%d'", expectedCount, len)
	}
}

func testPostEvents(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedCount := 1
	expectedMsg := "OK"

	event := dto.EventDTORequest{
		Action:         "TestAction",
		ResourceID:     "aws-cluster-2-aws-infra-2",
		ResourceType:   "cluster",
		EventTimestamp: time.Now(),
		Result:         "Pending",
		Severity:       "info",
		TriggeredBy:    "tester",
		Description:    nil,
	}
	b, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal data in post request: %v", err)
	}

	// Posting test data
	resp, err := http.Post(APIEventsURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

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

func testUpdateEvent(t *testing.T) {
	expectedHTTPCode := http.StatusOK
	expectedCount := 1
	expectedMsg := "OK"

	event := dto.EventDTORequest{
		ID:             2,
		Action:         "TestAction",
		ResourceID:     "aws-cluster-2-aws-infra-2",
		ResourceType:   "cluster",
		EventTimestamp: time.Now(),
		Result:         "OK-Updated",
		Severity:       "info",
		TriggeredBy:    "tester",
		Description:    nil,
	}
	b, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal data in post request: %v", err)
	}

	// Posting test data
	req, err := http.NewRequest(http.MethodPatch, APIEventsURL, bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to create PATCH request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PATCH request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

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
