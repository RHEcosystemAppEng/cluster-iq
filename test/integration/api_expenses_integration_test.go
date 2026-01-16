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
	APIExpensesURL = APIBaseURL + "/expenses"
)

func TestAPIExpenses(t *testing.T) {
	waitForAPIReady(t)

	if err := refreshInventory(); err != nil {
		t.Fatal("Error refreshing inventory")
	}

	t.Run("Test List System Expenses", func(t *testing.T) { testListExpenses(t) })
	t.Run("Test List Cluster Expenses", func(t *testing.T) { testListExpensesWithPagination(t) })
	// TODO Add filter testing
	t.Run("Test Post Expenses", func(t *testing.T) { testPostExpenses(t) })
}

func testListExpenses(t *testing.T) {
	expectedCount := 10
	expectedHTTPCode := http.StatusOK

	// Getting accounts data
	resp, err := http.Get(APIExpensesURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ExpenseDTOResponse]
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

func testListExpensesWithPagination(t *testing.T) {
	expectedCount := 8
	expectedHTTPCode := http.StatusOK

	// Getting accounts data
	resp, err := http.Get(APIExpensesURL + "?page=1&page_size=8")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.ExpenseDTOResponse]
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

func testPostExpenses(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedCount := 2
	expectedMsg := "OK"

	ts, _ := time.Parse("2006-01-02", "2025-02-02")
	expense := []dto.ExpenseDTORequest{
		{
			InstanceID: "id-3123456789Z",
			Amount:     5.0,
			Date:       ts,
		},
		{
			InstanceID: "id-2123456789Z",
			Amount:     16.0,
			Date:       ts,
		},
	}
	b, err := json.Marshal(expense)
	if err != nil {
		t.Fatalf("Failed to marshal data in post request: %v", err)
	}

	// Posting test data
	resp, err := http.Post(APIExpensesURL, "application/json", bytes.NewBuffer(b))
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
