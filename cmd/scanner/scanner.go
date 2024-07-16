package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/stocker"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

const (
	apiAccountEndpoint    = "/accounts"
	apiClusterEndpoint    = "/clusters"
	apiInstanceEndpoint   = "/instances"
	apiExpenseEndpoint    = "/expenses"
	defaultINISectionName = "__DEFAULT__"
)

var (
	// version reflects the current version of the API
	version string
	// commit reflects the git short-hash of the compiled version
	commit string
	//
	logger    *zap.Logger
	apiURL    string
	credsFile string
	// client http
	client http.Client
)

// Scanner models the cloud agnostic Scanner for looking up OCP deployments
type Scanner struct {
	inventory inventory.Inventory
	stockers  []stocker.Stocker
	apiURL    string
	credsFile string
	logger    *zap.Logger
}

// NewScanner creates and returns a new Scanner instance
func NewScanner(apiURL string, credsFile string, logger *zap.Logger) *Scanner {
	return &Scanner{
		inventory: *inventory.NewInventory(),
		stockers:  make([]stocker.Stocker, 0),
		apiURL:    apiURL,
		credsFile: credsFile,
		logger:    logger,
	}
}

func init() {
	// Initialize logging configuration.
	logger = ciqLogger.NewLogger()

	// Load configuration from environment variables.
	apiURL = os.Getenv("CIQ_API_URL")
	credsFile = os.Getenv("CIQ_CREDS_FILE")

	// Setting INI files default section name
	ini.DefaultSection = defaultINISectionName

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = http.Client{Transport: tr}

}

// gerProvider checks a incoming string and returns the corresponding inventory.CloudProvider value
func getProvider(provider string) inventory.CloudProvider {
	switch strings.ToUpper(provider) {
	case "AWS":
		return inventory.AWSProvider
	case "GCP":
		return inventory.GCPProvider
	case "AZURE":
		return inventory.AzureProvider
	}
	return inventory.UnknownProvider
}

// readCloudProviderAccounts reads and loads cloud provider accounts from a credentials file.
func (s *Scanner) readCloudProviderAccounts() error {
	// Load cloud accounts credentials file.
	cfg, err := ini.Load(s.credsFile)
	if err != nil {
		return err
	}

	// Removing the default sections because every account should have its own
	// section, and any parameter outside of an account section will be considered
	cfg.DeleteSection(defaultINISectionName)

	// Read INI file content.
	for _, account := range cfg.Sections() {
		newAccount := inventory.NewAccount(
			"",
			account.Name(),
			getProvider(account.Key("provider").String()),
			account.Key("user").String(),
			account.Key("key").String(),
		)
		if err := s.inventory.AddAccount(newAccount); err != nil {
			return err
		}
	}

	return nil
}

// createStockers creates and configures stocker instances for each provided account to be inventoried.
func (s *Scanner) createStockers() error {
	instances, err := s.getInstances()
	if err != nil {
		return err
	}
	for i := range s.inventory.Accounts {
		account := s.inventory.Accounts[i]
		switch account.Provider {
		case inventory.AWSProvider:
			s.logger.Info("Adding the AWS account to be inventoried", zap.String("account", account.Name))
			s.stockers = append(s.stockers, stocker.NewAWSStocker(account, s.logger, instances))
		case inventory.GCPProvider:
			logger.Warn("Failed to scan GCP account",
				zap.String("account", account.Name),
				zap.String("reason", "not implemented"),
			)
			// TODO: Uncomment line below when Stocker is implemented
			//s.stockers = append(s.stockers, stocker.NewGCPStocker(&account, logger))
		case inventory.AzureProvider:
			logger.Warn("Failed to scan Azure account",
				zap.String("account", account.Name),
				zap.String("reason", "not implemented"),
			)
			// TODO: Uncomment line below when Stocker is implemented
			//s.stockers = append(s.stockers, stocker.NewAzureStocker(&account, logger))
		}
	}

	if len(s.stockers) == 0 {
		return fmt.Errorf("Any account has been provided for scanning on credentials file")
	}

	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("Total Stockers created", zap.Int("count", len(s.stockers)))
		for i, stocker := range s.stockers {
			s.logger.Debug("Stocker", zap.Int("id", i), zap.String("name", stocker.GetResults().Name))
		}
	}

	return nil
}

// startStockers runs every stocker instance
func (s *Scanner) startStockers() error {
	for _, stockerInstance := range s.stockers {
		err := stockerInstance.MakeStock()
		if err != nil {
			return err
		}
	}
	return nil
}

// postNewInstance posts into the API, the new instances obtained after scanning
func (s *Scanner) postNewInstances(instances []inventory.Instance) error {
	s.logger.Debug("Posting new Instances")
	b, err := json.Marshal(instances)
	if err != nil {
		logger.Error("Failed to marshal inventory data from database", zap.Error(err))
		return err
	}

	requestURL := fmt.Sprintf("%s%s", s.apiURL, apiInstanceEndpoint)

	return postData(requestURL, b, s.logger)
}

// postNewInstance posts into the API, the new instances obtained after scanning
func (s *Scanner) postNewExpenses(expenses []inventory.Expense) error {
	s.logger.Debug("Posting new Instances")
	b, err := json.Marshal(expenses)
	if err != nil {
		logger.Error("Failed to marshal inventory data from database", zap.Error(err))
		return err
	}

	requestURL := fmt.Sprintf("%s%s", s.apiURL, apiExpenseEndpoint)

	return postData(requestURL, b, s.logger)
}

// postNewCluster posts into the API, the new instances obtained after scanning
func (s *Scanner) postNewClusters(clusters []inventory.Cluster) error {
	s.logger.Debug("Posting new Clusters")
	b, err := json.Marshal(clusters)
	if err != nil {
		logger.Error("Failed to marshal inventory data from database", zap.Error(err))
		return err
	}

	requestURL := fmt.Sprintf("%s%s", s.apiURL, apiClusterEndpoint)

	return postData(requestURL, b, s.logger)
}

// postNewAccount posts into the API, the new instances obtained after scanning
func (s *Scanner) postNewAccounts(accounts []inventory.Account) error {
	s.logger.Debug("Posting new Accounts")
	b, err := json.Marshal(accounts)
	if err != nil {
		s.logger.Error("Failed to marshal inventory data from database", zap.Error(err))
		return err
	}

	requestURL := fmt.Sprintf("%s%s", s.apiURL, apiAccountEndpoint)

	return postData(requestURL, b, s.logger)
}

func postData(url string, b []byte, logger *zap.Logger) error {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
		if err != nil {
			logger.Error("Request Failed", zap.String("response", response.Status), zap.Error(err))
			return err
		}
	} else if err != nil {
		logger.Error("Can't request to API. Response is Null", zap.Error(err))
		return err
	}

	return nil
}

func (s *Scanner) postScannerResults() error {
	var accounts []inventory.Account
	var clusters []inventory.Cluster
	var instances []inventory.Instance
	var expenses []inventory.Expense
	for _, account := range s.inventory.Accounts {
		for _, cluster := range account.Clusters {
			for _, instance := range cluster.Instances {
				for _, expense := range instance.Expenses {
					expenses = append(expenses, expense)
				}
				instance.Expenses = nil
				instances = append(instances, instance)

			}
			cluster.Instances = nil
			clusters = append(clusters, *cluster)
		}
		account.Clusters = nil
		accounts = append(accounts, *account)
	}

	if err := s.postNewAccounts(accounts); err != nil {
		s.logger.Debug("Adding Accounts", zap.Int("accounts_count", len(accounts)))
		return err
	}

	if err := s.postNewClusters(clusters); err != nil {
		s.logger.Debug("Adding Clusters", zap.Int("clusters_count", len(accounts)))
		return err
	}

	if err := s.postNewInstances(instances); err != nil {
		s.logger.Debug("Adding Instances", zap.Int("instances_count", len(accounts)))
		return err
	}

	if err := s.postNewExpenses(expenses); err != nil {
		s.logger.Debug("Adding Expenses", zap.Any("Expenses", expenses))
		return err
	}

	return nil
}

// signalHandler for managing incoming OS signals
func signalHandler(signal os.Signal) {
	if signal == syscall.SIGTERM {
		logger.Fatal("SIGTERM signal received. Stopping ClusterIQ Scanner")
		os.Exit(0)
	} else {
		logger.Warn("Ignoring signal: ", zap.String("signal_id", signal.String()))
	}
}

func main() {
	// Ignore Logger sync error
	defer func() { _ = logger.Sync() }()

	scan := NewScanner(apiURL, credsFile, logger)
	scan.logger.Info("Starting ClusterIQ Scanner",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("credentials file", credsFile),
	)

	// Listen Signals
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGTERM)
		s := <-quitChan
		signalHandler(s)
		logger.Info("Scanner stopped")
	}()

	var err error

	// Get Cloud Accounts from credentials file
	err = scan.readCloudProviderAccounts()
	if err != nil {
		logger.Error("Failed to get cloud provider accounts", zap.Error(err))
		return
	}

	// Run Stockers
	err = scan.createStockers()
	if err != nil {
		logger.Error("Failed to create stockers", zap.Error(err))
		return
	}

	err = scan.startStockers()
	if err != nil {
		logger.Error("Failed to start up stocker instances", zap.Error(err))
		return
	}

	// Writing into DB
	scan.inventory.PrintInventory()
	if err := scan.postScannerResults(); err != nil {
		logger.Error("Can't post scanned results", zap.Error(err))
		return
	}

	logger.Info("Scanner finished successfully")
}

// getInstances fetches instances from the backend API
func (s *Scanner) getInstances() ([]inventory.Instance, error) {
	s.logger.Debug("Fetching instances from backend")

	requestURL := fmt.Sprintf("%s%s", s.apiURL, apiInstanceEndpoint)

	resp, err := http.Get(requestURL)
	if err != nil {
		s.logger.Error("Failed to get instances from API", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Failed to get instances from API", zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to get instances, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.Error(err))
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		s.logger.Error("Failed to unmarshal response body", zap.Error(err))
		return nil, err
	}

	instancesData, ok := result["instances"]
	if !ok {
		s.logger.Error("Instances field not found in the response")
		return nil, fmt.Errorf("instances field not found in the response")
	}

	instancesJSON, err := json.Marshal(instancesData)
	if err != nil {
		s.logger.Error("Failed to marshal instances data", zap.Error(err))
		return nil, err
	}

	var instances []inventory.Instance
	err = json.Unmarshal(instancesJSON, &instances)
	if err != nil {
		s.logger.Error("Failed to unmarshal instances JSON", zap.Error(err))
		return nil, err
	}

	s.logger.Debug("Successfully fetched instances from backend", zap.Any("the instances are ", instances))
	return instances, nil
}
