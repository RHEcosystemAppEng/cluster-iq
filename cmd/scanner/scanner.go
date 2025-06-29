package main

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/credentials"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/stocker"
	"go.uber.org/zap"
)

const (
	apiAccountEndpoint  = "/accounts"
	apiClusterEndpoint  = "/clusters"
	apiInstanceEndpoint = "/instances"
	apiExpenseEndpoint  = "/expenses"
)

var (
	// version reflects the current version of the API
	version string
	// commit reflects the git short-hash of the compiled version
	commit string

	// logger variable across the entire scanner code
	logger *zap.Logger

	// MD5 Checksum of credsFile
	credsFileHash []byte

	// HTTP Client for connecting the scanner to the API
	client http.Client

	// API URL as global var for postData function
	APIURL string
)

// Scanner models the cloud agnostic Scanner for looking up OCP deployments
type Scanner struct {
	inventory inventory.Inventory
	stockers  []stocker.Stocker
	cfg       *config.ScannerConfig
	logger    *zap.Logger
}

// NewScanner creates and returns a new Scanner instance
func NewScanner(cfg *config.ScannerConfig, logger *zap.Logger) *Scanner {
	// Calculate Credentials file MD5 checksum for checking on runtime
	hash := md5.Sum([]byte(cfg.CredentialsFile))
	copy(hash[:], credsFileHash)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = http.Client{Transport: tr}
	APIURL = cfg.APIURL

	return &Scanner{
		inventory: *inventory.NewInventory(),
		stockers:  make([]stocker.Stocker, 0),
		cfg:       cfg,
		logger:    logger,
	}
}

func init() {
	// Initialize logging configuration.
	logger = ciqLogger.NewLogger()
}

// readCloudProviderAccounts reads and loads cloud provider accounts from a credentials file.
func (s *Scanner) readCloudProviderAccounts() error {
	// Load cloud accounts credentials file.
	accounts, err := credentials.ReadCloudAccounts(s.cfg.CredentialsFile)
	if err != nil {
		return err
	}

	// Read INI file content.
	for _, account := range accounts {
		newAccount := inventory.NewAccount(
			"",
			account.Name,
			account.Provider,
			account.User,
			account.Key,
		)
		// Getting billing enabled flag from config
		if account.BillingEnabled {
			newAccount.EnableBilling()
		}

		// Adding account to Inventory for scanning
		if err := s.inventory.AddAccount(newAccount); err != nil {
			return err
		}
	}

	return nil
}

// createStockers creates and configures stocker instances for each provided account to be inventoried.
func (s *Scanner) createStockers() error {
	var skippedAccounts int = 0
	var validStockers []stocker.Stocker
	for _, account := range s.inventory.Accounts {
		switch account.Provider {
		case inventory.AWSProvider:
			s.logger.Info("Processing AWS account", zap.String("account", account.Name))

			// AWS API Stoker
			awsStocker, err := stocker.NewAWSStocker(account, s.cfg.SkipNoOpenShiftInstances, s.logger)
			if err != nil {
				s.logger.Error("Failed to create AWS stocker; skipping this account",
					zap.String("account", account.Name),
					zap.Error(err))
				skippedAccounts++
				continue
			}
			validStockers = append(validStockers, awsStocker)

			// AWS Billing API Stoker
			if account.IsBillingEnabled() {
				s.logger.Warn("Enabled AWS Billing Stocker", zap.String("account", account.Name))
				instancesToScan, err := s.getInstancesForBillingUpdate()
				if err != nil {
					s.logger.Error("Failed to retrieve the list of instances required for billing information from AWS Cost Explorer.",
						zap.String("account", account.Name))
				} else {
					validStockers = append(validStockers, stocker.NewAWSBillingStocker(account, s.logger, instancesToScan))
				}
			}
		case inventory.GCPProvider:
			s.logger.Warn("Failed to scan GCP account",
				zap.String("account", account.Name),
				zap.String("reason", "not implemented"),
			)
			// TODO: Uncomment line below when GCP Stocker is implemented
			// gcpStocker = stocker.NewGCPStocker(account, s.cfg.SkipNoOpenShiftInstances, s.logger))
		case inventory.AzureProvider:
			s.logger.Warn("Failed to scan Azure account",
				zap.String("account", account.Name),
				zap.String("reason", "not implemented"),
			)
			// TODO: Uncomment line below when Azure Stocker is implemented
			// azureStocker = stocker.NewAzureStocker(account, s.cfg.SkipNoOpenShiftInstances, s.logger))

		default:
			s.logger.Warn("Unsupported cloud provider, skipping account",
				zap.String("account", account.Name),
				zap.String("provider", string(account.Provider)))
			continue
		}
	}

	s.stockers = validStockers
	s.logger.Info("Account registration complete",
		zap.Int("registeredAccounts", len(s.inventory.Accounts)),
		zap.Int("registeredStockers", len(s.stockers)),
		zap.Int("skippedAccounts", skippedAccounts))

	// If there are no stockers, nothing to do
	if len(s.stockers) == 0 {
		return fmt.Errorf("No valid accounts found for scanning")
	}

	// Checking the logLevel before entering on the For loop for optimization
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
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.stockers))

	// Running Stockers.MakeStock procedure on parallel
	for _, stockerInstance := range s.stockers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := stockerInstance.MakeStock(); err != nil {
				errChan <- err
			}
		}()
	}

	// Waiting for every Stock
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collecting stockers errors
	var errorList []error
	for err := range errChan {
		errorList = append(errorList, err)
	}

	// Processing errors when every stocker has finished
	if len(errorList) > 0 {
		for _, err := range errorList {
			s.logger.Error("Stocker Error", zap.Error(err))
		}
		return fmt.Errorf("Error when running Scanner stockers. Failed Stockers: (%d)", len(errorList))
	}

	s.logger.Info("Stockers executed correctly")
	return nil
}

// postNewAccount posts into the API an account, its clusters, instances and expenses
func (s *Scanner) postNewAccount(account inventory.Account) error {
	s.logger.Debug("Posting new Account", zap.String("account", account.Name))

	// Converting to Array because API handler assumes a list of accounts
	var accounts []inventory.Account
	accounts = append(accounts, account)
	b, err := json.Marshal(accounts)
	if err != nil {
		s.logger.Error("Failed to marshal account", zap.String("account", account.Name), zap.Error(err))
		return err
	}

	// Posting Account data
	if err := postData(apiAccountEndpoint, b); err != nil {
		return err
	}

	var clusters []inventory.Cluster
	var instances []inventory.Instance
	var expenses []inventory.Expense
	for _, cluster := range account.Clusters {
		for _, instance := range cluster.Instances {
			for _, expense := range instance.Expenses {
				expenses = append(expenses, expense)
			}
			instances = append(instances, instance)

		}
		clusters = append(clusters, *cluster)
	}

	// Posting Clusters
	if len(clusters) > 0 {
		if err := postClusters(clusters); err != nil {
			return err
		}
	}

	// Posting Instances
	if len(instances) > 0 {
		if err := postInstances(instances); err != nil {
			return err
		}
	}

	// Posting Expenses
	if len(expenses) > 0 {
		if err := postExpenses(expenses); err != nil {
			return err
		}
	}
	return nil
}

// postClusters posts into the API, the new instances obtained after scanning
func postClusters(clusters []inventory.Cluster) error {
	b, err := json.Marshal(clusters)
	if err != nil {
		return err
	}

	return postData(apiClusterEndpoint, b)
}

// postInstances posts into the API, the instances obtained after scanning
func postInstances(instances []inventory.Instance) error {
	b, err := json.Marshal(instances)
	if err != nil {
		return err
	}

	return postData(apiInstanceEndpoint, b)
}

// postExpenses posts into the API, the expenses obtained after scanning
func postExpenses(expenses []inventory.Expense) error {
	b, err := json.Marshal(expenses)
	if err != nil {
		return err
	}

	return postData(apiExpenseEndpoint, b)
}

// postScannerInventory posts to ClusterIQ API the information obtained of the scanning process
// This function parallelizes the post operations creating a thread by account(or stocker)
func (s *Scanner) postScannerInventory() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.inventory.Accounts))

	for _, account := range s.inventory.Accounts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.postNewAccount(*account); err != nil {
				errChan <- err
			}
		}()

	}
	// Waiting for every Stock
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collecting account posting errors
	var errorList []error
	for err := range errChan {
		errorList = append(errorList, err)
	}

	// Processing errors when every post account operation has finished
	if len(errorList) > 0 {
		for _, err := range errorList {
			s.logger.Error("Post Account Error", zap.Error(err))
		}
		return fmt.Errorf("Error when posting Scanner inventory")
	}

	s.logger.Info("Inventory posted correctly")
	return nil
}

func postData(path string, b []byte) error {
	url := fmt.Sprintf("%s%s", APIURL, path)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}

	return nil
}

// getInstances fetches instances from the backend API
func (s *Scanner) getInstancesForBillingUpdate() ([]inventory.Instance, error) {
	s.logger.Debug("Fetching instances for update billing from backend")

	requestURL := fmt.Sprintf("%s%s", s.cfg.APIURL, apiInstanceEndpoint+"/expense_update")

	resp, err := http.Get(requestURL)
	if err != nil {
		s.logger.Error("Failed to get last expenses from API", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Failed to get last expenses from API", zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to get last expenses, status code: %d", resp.StatusCode)
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

	s.logger.Debug("Successfully fetched instances from backend", zap.Int("instances_num", len(instances)))
	return instances, nil
}

// signalHandler for managing incoming OS signals
func signalHandler(sig os.Signal) {
	if sig == syscall.SIGTERM {
		logger.Fatal("SIGTERM signal received. Stopping ClusterIQ Scanner")
		os.Exit(0)
	}

	logger.Warn("Ignoring signal: ", zap.String("signal_id", sig.String()))
}

// Main method
func main() {
	// Ignore Logger sync error
	defer func() { _ = logger.Sync() }()

	var err error

	cfg, err := config.LoadScannerConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	scan := NewScanner(cfg, logger)

	scan.logger.Info("==================== Starting ClusterIQ Scanner ====================",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("credentials_file_path", cfg.CredentialsFile),
		zap.ByteString("credentials_file_hash", credsFileHash),
	)

	// Listen Signals block for receive OS signals. This is used by K8s/OCP for
	// interacting with this software when it's deployed on a Pod
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGTERM)
		s := <-quitChan
		signalHandler(s)
		logger.Info("Scanner stopped")
	}()

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
	if err := scan.postScannerInventory(); err != nil {
		logger.Error("Can't post scanned results", zap.Error(err))
		return
	}

	logger.Info("Scanner finished successfully")
}
