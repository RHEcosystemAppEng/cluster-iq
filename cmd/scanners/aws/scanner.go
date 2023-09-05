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
	// TODO: comment rest of global vars
	inven     *inventory.Inventory
	stockers  []stocker.Stocker
	dbURL     string
	dbPass    string
	credsFile string
	logger    *zap.Logger
)

func init() {
	// Initialize logging configuration.
	logger = ciqLogger.NewLogger()

	// Load configuration from environment variables.
	dbHost := os.Getenv("CIQ_DB_HOST")
	dbPort := os.Getenv("CIQ_DB_PORT")
	dbPass = os.Getenv("CIQ_DB_PASS")
	credsFile = os.Getenv("CIQ_CREDS_FILE")
	dbURL = fmt.Sprintf("%s:%s", dbHost, dbPort)
}

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

// getCloudProviderAccounts retrieves cloud provider accounts from a credentials file.
func getCloudProviderAccounts() ([]inventory.Account, error) {
	accounts := make([]inventory.Account, 0)

	// Load cloud accounts credentials file.
	cfg, err := ini.Load(credsFile)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Read INI file content.
	for _, account := range cfg.Sections() {
		newAccount := inventory.NewAccount(
			account.Name(),
			getProvider(account.Key("provider").String()),
			account.Key("user").String(),
			account.Key("key").String(),
		)
		accounts = append(accounts, *newAccount)
	}

	return accounts, nil
}

// createStockers creates and configures stocker instances for each provided account to be inventoried.
func createStockers(accounts []inventory.Account) error {
	for _, account := range accounts {
		switch account.Provider {
		case inventory.AWSProvider:
			logger.Info("Adding the AWS account to be inventoried", zap.String("account", account.Name))
			stockers = append(stockers, stocker.NewAWSStocker(account))
		case inventory.GCPProvider:
			logger.Warn("Failed to scan GCP account",
				zap.String("account", account.Name),
				zap.String("reason", "not implemented"),
			)
		case inventory.AzureProvider:
			logger.Warn("Failed to scan Azure account",
				zap.String("account", account.Name),
				zap.String("reason", "not implemented"),
			)
		}
	}

	if len(stockers) == 0 {
		return fmt.Errorf("no valid stockers created")
	}
	return nil
}

// startStockers runs every stocker instance
func startStockers() error {
	for _, stockerInstance := range stockers {
		err := stockerInstance.MakeStock()
		if err != nil {
			return err
		}
		inven.AddAccount(stockerInstance.GetResults())
	}
	return nil
}

func main() {
	// Ignore Logger sync error
	defer func() { _ = logger.Sync() }()

	rdb, err := redis.InitDatabase(dbURL, dbPass)
	if err != nil {
		logger.Fatal("Failed to establish database connection", zap.Error(err))
	}
	// Prepare New Stock
	inven = inventory.NewInventory()

	// Get Cloud Accounts from credentials file
	accounts, err := getCloudProviderAccounts()
	if err != nil {
		logger.Error("Failed to get cloud provider accounts", zap.Error(err))
		return
	}

	// Run Stockers
	err = createStockers(accounts)

	if err != nil {
		logger.Error("Failed to add stockers", zap.Error(err))
		return
	}

	err = startStockers()
	if err != nil {
		logger.Error("Failed to start up stocker instances", zap.Error(err))
		return
	}

	b, err := json.Marshal(inven)
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
