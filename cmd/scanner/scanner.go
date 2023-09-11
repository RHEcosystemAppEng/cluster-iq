package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/redis"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/stocker"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

var (
	// version reflects the current version of the API
	version string
	// commit reflects the git short-hash of the compiled version
	commit string
	//
	logger    *zap.Logger
	dbURL     string
	dbPass    string
	credsFile string
)

// Scanner models the cloud agnostic Scanner for looking up OCP deployments
type Scanner struct {
	inventory inventory.Inventory
	stockers  []stocker.Stocker
	dbURL     string
	dbPass    string
	credsFile string
	logger    *zap.Logger
}

// NewScanner creates and returns a new Scanner instance
func NewScanner(dbURL string, dbPass string, credsFile string, logger *zap.Logger) *Scanner {
	return &Scanner{
		inventory: *inventory.NewInventory(),
		stockers:  make([]stocker.Stocker, 0),
		dbURL:     dbURL,
		dbPass:    dbPass,
		credsFile: credsFile,
		logger:    logger,
	}
}

func init() {
	// Initialize logging configuration.
	logger = ciqLogger.NewLogger()

	// Load configuration from environment variables.
	dbHost := os.Getenv("CIQ_DB_HOST")
	dbPort := os.Getenv("CIQ_DB_PORT")
	dbURL = fmt.Sprintf("%s:%s", dbHost, dbPort)
	dbPass = os.Getenv("CIQ_DB_PASS")
	credsFile = os.Getenv("CIQ_CREDS_FILE")
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

	// Read INI file content.
	for _, account := range cfg.Sections() {
		newAccount := inventory.NewAccount(
			account.Name(),
			getProvider(account.Key("provider").String()),
			account.Key("user").String(),
			account.Key("key").String(),
		)
		if err := s.inventory.AddAccount(*newAccount); err != nil {
			return err
		}
	}

	return nil
}

// createStockers creates and configures stocker instances for each provided account to be inventoried.
func (s *Scanner) createStockers() error {
	for i := range s.inventory.Accounts {
		account := s.inventory.Accounts[i]
		switch account.Provider {
		case inventory.AWSProvider:
			s.logger.Info("Adding the AWS account to be inventoried", zap.String("account", account.Name))
			s.stockers = append(s.stockers, stocker.NewAWSStocker(&account, logger))
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

func main() {
	// Ignore Logger sync error
	defer func() { _ = logger.Sync() }()

	scan := NewScanner(dbURL, dbPass, credsFile, logger)
	scan.logger.Info("Starting ClusterIQ Scanner",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("dbURL", dbURL),
		zap.String("credentials file", credsFile),
	)

	rdb, err := redis.InitDatabase(dbURL, dbPass)
	if err != nil {
		logger.Fatal("Failed to establish database connection", zap.Error(err))
	}

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

	b, err := json.Marshal(scan.inventory)
	if err != nil {
		logger.Error("Failed to marshal inventory data from database", zap.Error(err))
		return
	}

	ctx := context.Background()
	logger.Info("Saving results to the database")
	// TODO Refactor into dedicated function
	err = rdb.Set(ctx, "Stock", string(b), redis.DataExpirationTTL).Err()
	if err != nil {
		logger.Error("Failed to write results to the database", zap.Error(err))
		return
	}

	logger.Info("Scanner finished successfully")
}
