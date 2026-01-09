package integration

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	APIBaseURL             = "http://localhost:8081/api/v1"
	APIHealthcheckURL      = APIBaseURL + "/healthcheck"
	APIInventoryRefreshURL = APIBaseURL + "/inventory"
)

// loadTestData loads test data files from a JSON file
func loadTestData(filename string) []byte {
	path := filepath.Join("./data_test_files", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func refreshInventory() error {
	_, err := http.Post(APIInventoryRefreshURL, "", bytes.NewReader([]byte{}))

	return err
}

// waitForAPIReady tries to reach the API endpoint until it responds or times out
func waitForAPIReady(t *testing.T) {
	t.Helper()
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(APIHealthcheckURL)
		if err == nil && resp.StatusCode < 500 {
			return
		}
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("API did not become ready at %s", APIHealthcheckURL)
}

func checkHTTPResponseCode(t *testing.T, response *http.Response, expectedHTTPCode int) {
	if response.StatusCode != expectedHTTPCode {
		t.Fatalf("Expected HTTP Response code: %d, got %d", expectedHTTPCode, response.StatusCode)
	}
}
