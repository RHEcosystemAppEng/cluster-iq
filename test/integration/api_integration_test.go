package integration

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	APIBaseURL        = "http://localhost:8081/api/v1"
	APIHealthcheckURL = APIBaseURL + "/healthcheck"
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

// waitForAPIReady tries to reach the API endpoint until it responds or times out
func waitForAPIReady(t *testing.T) {
	url := APIHealthcheckURL
	t.Helper()
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode < 500 {
			return
		}
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("API did not become ready at %s", url)
}
