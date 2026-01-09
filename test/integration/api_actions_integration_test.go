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
	APIActionsURL = APIBaseURL + "/actions"
)

func TestAPIActions(t *testing.T) {
	waitForAPIReady(t)

	if err := refreshInventory(); err != nil {
		t.Fatal("Error refreshing inventory")
	}

	t.Run("List Actions", func(t *testing.T) { testListActions(t) })
	t.Run("List Actions with Pagination", func(t *testing.T) { testListActionsWithPagination(t) })
	t.Run("Get Actions By ID Success", func(t *testing.T) { testGetActionByID_Exists(t) })
	t.Run("Get Actions By ID Not Found", func(t *testing.T) { testGetActionByID_NoExists(t) })
	t.Run("Post One Action", func(t *testing.T) { testPostOneAction(t) })
	t.Run("Post Multiple Action", func(t *testing.T) { testPostMultipleActions(t) })
	t.Run("Post Wrong Action", func(t *testing.T) { testPostWrongAction(t) })
	t.Run("Enable Action", func(t *testing.T) { testEnableAction(t) })
	t.Run("Disable Action", func(t *testing.T) { testDisableAction(t) })
	t.Run("Delete existing Action", func(t *testing.T) { testDeleteAction_Exists(t) })
	t.Run("Delete no existing Action", func(t *testing.T) { testDeleteAction_NoExists(t) })
}

func testListActions(t *testing.T) {
	expectedCount := 3
	expectedHTTPCode := http.StatusOK

	// Getting actions data
	resp, err := http.Get(APIActionsURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ActionDTOResponse]
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

func testListActionsWithPagination(t *testing.T) {
	expectedCount := 2
	expectedHTTPCode := http.StatusOK

	// Getting Actions data
	resp, err := http.Get(APIActionsURL + "?page=1&page_size=2")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ActionDTOResponse]
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

func testGetActionByID_Exists(t *testing.T) {
	expectedActionID := "1"
	expectedHTTPCode := http.StatusOK

	// Getting Clusters data
	resp, err := http.Get(APIActionsURL + "/" + expectedActionID)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.ActionDTOResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.ID != expectedActionID {
		t.Fatalf("Expected ID: '%s', got: '%s'", expectedActionID, response.ID)
	}
}

func testGetActionByID_NoExists(t *testing.T) {
	expectedMsg := "Scheduled action not found"
	expectedHTTPCode := http.StatusNotFound

	// Getting Clusters data
	resp, err := http.Get(APIActionsURL + "/" + "9999")
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
func postActions(t *testing.T, actions []dto.ActionDTORequest, expectedHTTPCode int) *http.Response {
	b, err := json.Marshal(actions)
	if err != nil {
		return nil
	}

	// Posting test data
	resp, err := http.Post(APIActionsURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	return resp
}

func testPostOneAction(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedMsg := "OK"
	expectedCount := 1

	ts, _ := time.Parse(time.RFC3339, "1970-02-02T10:00:00+00:00")
	payload := []dto.ActionDTORequest{
		{
			Type:      "scheduled_action",
			Time:      ts,
			Operation: "PowerOn",
			Status:    "Pending",
			Enabled:   false,
			ClusterID: "gcp-cluster-1-gcp-infra-1",
			Region:    "europe",
			AccountID: "gcp-project-1",
			Instances: []string{},
		},
	}

	// Posting test data
	resp := postActions(t, payload, expectedHTTPCode)
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

func testPostMultipleActions(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedMsg := "OK"
	expectedCount := 2

	ts, _ := time.Parse(time.RFC3339, "1970-02-02T10:00:00+00:00")
	payload := []dto.ActionDTORequest{
		{
			Type:      "cron_action",
			CronExp:   "45 */12 * 7 *",
			Operation: "PowerOn",
			Status:    "Pending",
			Enabled:   false,
			ClusterID: "gcp-cluster-2-gcp-infra-2",
			Region:    "europe",
			AccountID: "gcp-project-1",
			Instances: []string{},
		},
		{
			Type:      "scheduled_action",
			Time:      ts,
			Operation: "PowerOn",
			Status:    "Pending",
			Enabled:   false,
			ClusterID: "gcp-cluster-1-gcp-infra-1",
			Region:    "europe",
			AccountID: "gcp-project-1",
			Instances: []string{},
		},
	}

	// Posting test data
	resp := postActions(t, payload, expectedHTTPCode)
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

func testPostWrongAction(t *testing.T) {
	expectedHTTPCode := http.StatusInternalServerError
	expectedMsg := "Failed to create actions: error when processing action requests data"

	payload := []dto.ActionDTORequest{
		{
			Type:      "generic_action",
			Time:      time.Now(),
			CronExp:   "* * * * *",
			Operation: "power",
			Status:    "Pending",
			Enabled:   false,
			ClusterID: "aws-cluster-2-aws-infra-2",
			Region:    "europe",
			AccountID: "subs-00000001",
			Instances: []string{},
		},
	}

	// Posting test data
	resp := postActions(t, payload, expectedHTTPCode)
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

func testEnableAction(t *testing.T) {
	expectedHTTPCode := http.StatusOK

	actionID := "2"
	req, err := http.NewRequest(http.MethodPatch, APIActionsURL+"/"+actionID+"/enable", nil)
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

func testDisableAction(t *testing.T) {
	expectedHTTPCode := http.StatusOK

	actionID := "1"
	req, err := http.NewRequest(http.MethodPatch, APIActionsURL+"/"+actionID+"/disable", nil)
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

func testDeleteAction_Exists(t *testing.T) {
	expectedHTTPCode := http.StatusNoContent

	// Preparing DELETE request
	actionID := "2"
	req, err := http.NewRequest(http.MethodDelete, APIActionsURL+"/"+actionID, nil)
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

func testDeleteAction_NoExists(t *testing.T) {
	expectedHTTPCode := http.StatusNoContent

	// Preparing DELETE request
	actionID := "9999"
	req, err := http.NewRequest(http.MethodDelete, APIActionsURL+"/"+actionID, nil)
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
