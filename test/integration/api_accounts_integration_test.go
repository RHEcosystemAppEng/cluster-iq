package integration

import (
	"bytes"
	"encoding/json"
	"io"
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

	t.Run("TestGetAccount", func(t *testing.T) { testGetAccounts(t) })
	t.Run("TestGetAccountByID Success", func(t *testing.T) { testGetAccountByID_Exists(t) })
	t.Run("TestGetAccountByID Not Found", func(t *testing.T) { testGetAccountByID_NoExists(t) })
	t.Run("TestGetAccountClusters", func(t *testing.T) { testGetAccountClusters(t) })
	t.Run("TestPostOneAccount", func(t *testing.T) { testPostOneAccount(t) })
	t.Run("TestPostMultipleAccounts", func(t *testing.T) { testPostMultipleAccounts(t) })
	t.Run("TestDeleteAccount Success", func(t *testing.T) { testDeleteAccount_Exists(t) })
	t.Run("TestDeleteAccount Not Found", func(t *testing.T) { testDeleteAccount_NoExists(t) })
	t.Run("TestPatchAccount", func(t *testing.T) { testPatchAccount(t) })
}

func testGetAccounts(t *testing.T) {
	expectedCount := 3
	// Getting accounts data
	resp, err := http.Get(APIAccountsURL)
	if err != nil {
		t.Fatalf("Failed to make GetAccounts request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.ListResponse[dto.AccountDTOResponse]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GET response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected %d Accounts, got %d", expectedCount, response.Count)
	}
}

func testGetAccountByID_Exists(t *testing.T) {
	expectedAccountID := "gcp-project-1"

	// Getting Clusters data
	resp, err := http.Get(APIAccountsURL + "/" + expectedAccountID)
	if err != nil {
		t.Fatalf("Failed to make GetAccountByID request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.AccountDTOResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetAccountsByID response body: %v", err)
	}

	// Comparing data
	if response.AccountID != expectedAccountID {
		t.Fatalf("Expected: %s, Have: %s", expectedAccountID, response.AccountID)
	}
}

func testGetAccountByID_NoExists(t *testing.T) {
	expectedMsg := "Account not found"

	// Getting Clusters data
	resp, err := http.Get(APIAccountsURL + "/missing-account")
	if err != nil {
		t.Fatalf("Failed to make GetAccountByID request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.GenericErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetAccountsByID response body: %v", err)
	}

	// Comparing data
	if response.Message != expectedMsg {
		t.Fatalf("Expected: %s, Have: %s", expectedMsg, response.Message)
	}
}

func testGetAccountClusters(t *testing.T) {
	expectedCount := 2
	accountID := "subs-00000001"

	// Getting accounts data
	resp, err := http.Get(APIAccountsURL + "/" + accountID + "/clusters")
	if err != nil {
		t.Fatalf("Failed to make GetAccountClusters request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response dto.ClusterDTOResponseList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode GetAccountClusters response body: %v", err)
	}

	// Comparing data
	if response.Count != expectedCount {
		t.Fatalf("Expected: %d, Have: %d", expectedCount, response.Count)
	}
	// TODO Add elements check
}

func postAccounts(t *testing.T, accounts string) *responsetypes.PostResponse {
	// Posting test data
	resp, err := http.Post(APIAccountsURL, "application/json", bytes.NewBuffer([]byte(accounts)))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
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

func testPostOneAccount(t *testing.T) {
	// TODO Transform into dto.AccountDTORequestList
	payload := `
		[
			{
				"accountID": "ACC-001",
				"accountName": "test-account-001",
				"provider": "AWS",
				"lastScanTS": "1993-10-12T00:00:00Z",
				"createdAt": "1993-10-12T00:00:00Z"
			}
		]
	`

	// Posting test data
	response := postAccounts(t, payload)

	// Checks
	if response.Count != 1 {
		t.Fatalf("Expected 1 Posted Account, got %d", response.Count)
	}

	if response.Status != "OK" {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}

func testPostMultipleAccounts(t *testing.T) {
	// TODO Transform into dto.AccountDTORequestList
	payload := `
		[
			{
				"accountID": "ACC-002",
				"accountName": "test-account-002",
				"provider": "GCP",
				"last_scan_timestamp": "2014-10-12T00:00:00Z",
				"createdAt": "1993-10-12T00:00:00Z"
			},
			{
				"accountID": "ACC-003",
				"accountName": "test-account-003",
				"provider": "Azure",
				"last_scan_timestamp": "1970-10-12T00:00:00Z",
				"createdAt": "1993-10-12T00:00:00Z"
			}
		]
	`

	// Posting test data
	response := postAccounts(t, payload)

	if response.Count != 2 {
		t.Fatalf("Expected 2 Posted Account, got %d", response.Count)
	}

	if response.Status != "OK" {
		t.Fatalf("Unexpected Status Message: '%s'", response.Status)
	}
}

func testPatchAccount(t *testing.T) {
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
		t.Fatalf("Failed to create PatchAccount request: %v", err)
	}

	// Executing PATCH request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute PatchAccount request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusNotImplemented {
		t.Fatalf("Expected status 501, got %d", resp.StatusCode)
	}
}

func testDeleteAccount_Exists(t *testing.T) {
	// Preparing DELETE request
	accountID := "ACC-001"
	req, err := http.NewRequest(http.MethodDelete, APIAccountsURL+"/"+accountID, nil)
	if err != nil {
		t.Fatalf("Failed to create DeleteAccount request: %v", err)
	}

	// Executing DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute DeleteAccount request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func testDeleteAccount_NoExists(t *testing.T) {
	// Preparing DELETE request
	accountID := "missing"
	req, err := http.NewRequest(http.MethodDelete, APIAccountsURL+"/"+accountID, nil)
	if err != nil {
		t.Fatalf("Failed to create DeleteAccount request: %v", err)
	}

	// Executing DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute DeleteAccount request: %v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}
