package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

const (
	APIActionRunsURL = APIBaseURL + "/action-runs"
)

func TestAPIActionRuns(t *testing.T) {
	waitForAPIReady(t)

	if err := refreshInventory(); err != nil {
		t.Fatal("Error refreshing inventory")
	}

	t.Run("List Action Runs", func(t *testing.T) { testListActionRuns(t) })
	t.Run("List Action Runs with Pagination", func(t *testing.T) { testListActionRunsWithPagination(t) })
	t.Run("List Action Runs filtered by Status", func(t *testing.T) { testListActionRunsFilteredByStatus(t) })
	t.Run("Get Action Run By ID Success", func(t *testing.T) { testGetActionRunByID_Exists(t) })
	t.Run("Get Action Run By ID Not Found", func(t *testing.T) { testGetActionRunByID_NoExists(t) })
	t.Run("Post Action Run", func(t *testing.T) { testPostActionRun(t) })
	t.Run("Update Action Run", func(t *testing.T) { testUpdateActionRun(t) })
	t.Run("Update Action Run Not Found", func(t *testing.T) { testUpdateActionRun_NoExists(t) })
}

func testListActionRuns(t *testing.T) {
	expectedCount := 2
	expectedHTTPCode := http.StatusOK

	resp, err := http.Get(APIActionRunsURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	var response responsetypes.ListResponse[dto.ActionRunDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}

	if len := len(response.Items); len != expectedCount {
		t.Fatalf("Expected Items: '%d', got: '%d'", expectedCount, len)
	}
}

func testListActionRunsWithPagination(t *testing.T) {
	expectedCount := 1
	expectedHTTPCode := http.StatusOK

	resp, err := http.Get(APIActionRunsURL + "?page=1&page_size=1")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	var response responsetypes.ListResponse[dto.ActionRunDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}

	if len := len(response.Items); len != expectedCount {
		t.Fatalf("Expected Items: '%d', got: '%d'", expectedCount, len)
	}
}

func testListActionRunsFilteredByStatus(t *testing.T) {
	expectedCount := 1
	expectedHTTPCode := http.StatusOK

	resp, err := http.Get(APIActionRunsURL + "?status=Running")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	var response responsetypes.ListResponse[dto.ActionRunDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}

	if len := len(response.Items); len != expectedCount {
		t.Fatalf("Expected Items: '%d', got: '%d'", expectedCount, len)
	}

	if response.Items[0].Status != "Running" {
		t.Fatalf("Expected Status: 'Running', got: '%s'", response.Items[0].Status)
	}
}

func testGetActionRunByID_Exists(t *testing.T) {
	expectedRunID := "1"
	expectedHTTPCode := http.StatusOK

	resp, err := http.Get(APIActionRunsURL + "/" + expectedRunID)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	var response dto.ActionRunDTOResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if response.ID != expectedRunID {
		t.Fatalf("Expected ID: '%s', got: '%s'", expectedRunID, response.ID)
	}
}

func testGetActionRunByID_NoExists(t *testing.T) {
	expectedMsg := "Action run not found"
	expectedHTTPCode := http.StatusNotFound

	resp, err := http.Get(APIActionRunsURL + "/" + "9999")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	var response responsetypes.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func testPostActionRun(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedCount := 1

	payload := dto.ActionRunDTORequest{
		ScheduleID: "1",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data in post request: %v", err)
	}

	resp, err := http.Post(APIActionRunsURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	var response responsetypes.PostResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got '%d'", expectedCount, response.Count)
	}
}

func testUpdateActionRun(t *testing.T) {
	expectedHTTPCode := http.StatusOK

	payload := dto.ActionRunDTORequest{
		ScheduleID: "1",
		Status:     "Failed",
		ErrorMsg:   "test error",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data in patch request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, APIActionRunsURL+"/1", bytes.NewBuffer(b))
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

	checkHTTPResponseCode(t, resp, expectedHTTPCode)
}

func testUpdateActionRun_NoExists(t *testing.T) {
	expectedMsg := "Action run not found"
	expectedHTTPCode := http.StatusNotFound

	payload := dto.ActionRunDTORequest{
		ScheduleID: "1",
		Status:     "Failed",
		ErrorMsg:   "test error",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data in patch request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, APIActionRunsURL+"/9999", bytes.NewBuffer(b))
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

	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	var response responsetypes.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}
