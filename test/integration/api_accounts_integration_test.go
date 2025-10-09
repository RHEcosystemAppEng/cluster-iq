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
	APIAccountsURL = APIBaseURL + "/accounts"
)

func TestAPIAccounts(t *testing.T) {
	waitForAPIReady(t)

	t.Run("Test List Accounts", func(t *testing.T) { testListAccounts(t) })
	t.Run("Test List Accounts with Pagination", func(t *testing.T) { testListAccountsWithPagination(t) })
	t.Run("Test List Accounts By Provider", func(t *testing.T) { testListAccountsByProvider(t) })
	t.Run("Test List Accounts By Provider Wrong", func(t *testing.T) { testListAccountsByWrongProvider(t) })
	t.Run("Test Get Account By ID Success", func(t *testing.T) { testGetAccountByID_Exists(t) })
	t.Run("Test Get Account By ID Not Found", func(t *testing.T) { testGetAccountByID_NoExists(t) })
	t.Run("Test Get Account Clusters", func(t *testing.T) { testGetAccountClusters(t) })
	t.Run("Test Post One Account", func(t *testing.T) { testPostOneAccount(t) })
	t.Run("Test Post Multiple Accounts", func(t *testing.T) { testPostMultipleAccounts(t) })
	t.Run("Test Post Wrong Accounts", func(t *testing.T) { testPostWrongAccount(t) })
	t.Run("Test Patch Account", func(t *testing.T) { testPatchAccount(t) })
	t.Run("Test Delete Account Success", func(t *testing.T) { testDeleteAccount_Exists(t) })
	t.Run("Test Delete Account Not Found", func(t *testing.T) { testDeleteAccount_NoExists(t) })
}

func testListAccounts(t *testing.T) {
	expectedCount := 3
	expectedHTTPCode := http.StatusOK

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.AccountDTOResponse]
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

func testListAccountsWithPagination(t *testing.T) {
	expectedCount := 2
	expectedHTTPCode := http.StatusOK

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL + "?page=1&page_size=2")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.AccountDTOResponse]
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

func testListAccountsByProvider(t *testing.T) {
	expectedCount := 1
	expectedHTTPCode := http.StatusOK

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL + "?provider=AWS")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.AccountDTOResponse]
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

func testListAccountsByWrongProvider(t *testing.T) {
	expectedMsg := "Failed to retrieve accounts"
	expectedHTTPCode := http.StatusInternalServerError

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL + "?provider=ANYONE")
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

	if response.Message != expectedMsg {
		t.Fatalf("Expected Message: '%s', got: '%s'", expectedMsg, response.Message)
	}
}

func testGetAccountByID_Exists(t *testing.T) {
	expectedAccountID := "gcp-project-1"
	expectedHTTPCode := http.StatusOK

	// Getting Clusters data
	resp, err := http.Get(APIAccountsURL + "/" + expectedAccountID)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response dto.AccountDTOResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.AccountID != expectedAccountID {
		t.Fatalf("Expected AccountID: '%s', got: '%s'", expectedAccountID, response.AccountID)
	}
}

func testGetAccountByID_NoExists(t *testing.T) {
	expectedMsg := "Account not found"
	expectedHTTPCode := http.StatusNotFound

	// Getting Clusters data
	resp, err := http.Get(APIAccountsURL + "/" + "missing-account")
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

func testGetAccountClusters(t *testing.T) {
	expectedCount := 2
	expectedHTTPCode := http.StatusOK
	expectedAccountID := "subs-00000001"

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL + "/" + expectedAccountID + "/clusters")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	// Decode the JSON response
	var response responsetypes.ListResponse[dto.AccountDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected Count: '%d', got: '%d'", expectedCount, response.Count)
	}
}

func postAccounts(t *testing.T, accounts []dto.AccountDTORequest, expectedHTTPCode int) *http.Response {
	b, err := json.Marshal(accounts)
	if err != nil {
		return nil
	}

	// Posting test data
	resp, err := http.Post(APIAccountsURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}

	// Check response code
	checkHTTPResponseCode(t, resp, expectedHTTPCode)

	return resp
}

func testPostOneAccount(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedMsg := "OK"
	expectedCount := 1

	ts, _ := time.Parse(time.RFC3339, "2025-08-02T10:00:00+00:00")
	payload := []dto.AccountDTORequest{
		{
			AccountID:   "ACC-001",
			AccountName: "test-account-001",
			Provider:    "AWS",
			LastScanTS:  ts,
			CreatedAt:   ts,
		},
	}

	// Posting test data
	resp := postAccounts(t, payload, expectedHTTPCode)
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

func testPostMultipleAccounts(t *testing.T) {
	expectedHTTPCode := http.StatusCreated
	expectedMsg := "OK"
	expectedCount := 2

	ts, _ := time.Parse(time.RFC3339, "2025-08-02T10:00:00+00:00")
	payload := []dto.AccountDTORequest{
		{
			AccountID:   "ACC-002",
			AccountName: "test-account-002",
			Provider:    "GCP",
			LastScanTS:  ts,
			CreatedAt:   ts,
		},
		{
			AccountID:   "ACC-003",
			AccountName: "test-account-003",
			Provider:    "Azure",
			LastScanTS:  ts,
			CreatedAt:   ts,
		},
	}

	// Posting test data
	resp := postAccounts(t, payload, expectedHTTPCode)
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

func testPostWrongAccount(t *testing.T) {
	expectedHTTPCode := http.StatusInternalServerError
	expectedMsg := "Failed to create accounts: named-exec INSERT error: pq: invalid input value for enum cloud_provider: \"Provider\""

	ts, _ := time.Parse(time.RFC3339, "2025-08-02T10:00:00+00:00")
	payload := []dto.AccountDTORequest{
		{
			AccountID:   "ACC-004",
			AccountName: "test-account-002",
			Provider:    "Provider",
			LastScanTS:  ts,
			CreatedAt:   ts,
		},
	}

	// Posting test data
	resp := postAccounts(t, payload, expectedHTTPCode)
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

func testPatchAccount(t *testing.T) {
	expectedHTTPCode := http.StatusNotImplemented

	patchAccount := dto.AccountDTORequest{
		AccountID:   "ACC-001",
		AccountName: "test-account-003",
		Provider:    inventory.AWSProvider,
		LastScanTS:  time.Now(),
	}

	patchBody, err := json.Marshal(patchAccount)
	if err != nil {
		t.Fatalf("Failed to marshal updated account: %v", err)
	}

	// Preparing PATCH request
	req, err := http.NewRequest(http.MethodPatch, APIAccountsURL+"/ACC-003", bytes.NewBuffer(patchBody))
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

func testDeleteAccount_Exists(t *testing.T) {
	expectedHTTPCode := http.StatusNoContent

	// Preparing DELETE request
	accountID := "ACC-001"
	req, err := http.NewRequest(http.MethodDelete, APIAccountsURL+"/"+accountID, nil)
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

func testDeleteAccount_NoExists(t *testing.T) {
	expectedHTTPCode := http.StatusNoContent

	// Preparing DELETE request
	accountID := "missing"
	req, err := http.NewRequest(http.MethodDelete, APIAccountsURL+"/"+accountID, nil)
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
