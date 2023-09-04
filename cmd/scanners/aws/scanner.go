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
	// Logging config
	logger = ciqLogger.NewLogger()

	// Getting config
	dbHost := os.Getenv("CIQ_DB_HOST")
	dbPort := os.Getenv("CIQ_DB_PORT")
	dbPass = os.Getenv("CIQ_DB_PASS")
	credsFile = os.Getenv("CIQ_CREDS_FILE")
	dbURL = fmt.Sprintf("%s:%s", dbHost, dbPort)
}

// getProvider return a inventory.CloudProvider based on a string
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

// GetCloudProviderAccounts TODO
func GetCloudProviderAccounts() []inventory.Account {
	accounts := make([]inventory.Account, 0)

	// Getting cloud accounts credentials file
	cfg, err := ini.Load(credsFile)
	if err != nil {
		logger.Fatal("Can't Open credentials file", zap.Error(err))
	}

	// Reading INI file content
	for _, account := range cfg.Sections() {
		newAccount := inventory.NewAccount(
			account.Name(),
			getProvider(account.Key("provider").String()),
			account.Key("user").String(),
			account.Key("key").String(),
		)
		accounts = append(accounts, *newAccount)
	}

	return accounts
}

// createStockers creates and configures a stocker instance for each Account to be scraped
func createStockers(accounts []inventory.Account) error {
	for _, account := range accounts {
		switch account.Provider {
		case inventory.AWSProvider:
			logger.Info("Adding AWS account to be inventored\n", zap.String("account", account.Name))
			stockers = append(stockers, stocker.NewAWSStocker(account))
		case inventory.GCPProvider:
			err := fmt.Errorf("Google Cloud Platform (GCP) Stocker not implemented! Account %s will not be scanned", account.Name)
			logger.Error("Can't scan GCP account", zap.Error(err))
			return err
		case inventory.AzureProvider:
			err := fmt.Errorf("Microsoft Azure Stocker not implemented! Account %s will not be scanned", account.Name)
			logger.Error("Can't scan Azure account", zap.Error(err))
			return err
		}
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
		// TODO handle error properly
		inven.AddAccount(stockerInstance.GetResults())
	}
	return nil
}

func main() {
	defer logger.Sync()
	rdb, err := redis.InitDatabase(dbURL, dbPass)
	if err != nil {
		logger.Error("Failed to establish database connection", zap.Error(err))
	}
	// Prepare New Stock
	inven = inventory.NewInventory()

	// Get Cloud Accounts from credentials file
	accounts := GetCloudProviderAccounts()

	// Running Stockers
	// TODO Handle error properly
	createStockers(accounts)
	err = startStockers()
	if err != nil {
		logger.Fatal("Failed to start up stocker instances", zap.Error(err))
		return
	}

	b, err := json.Marshal(inven)
	if err != nil {
		logger.Fatal("Failed to marshal inventory data from DB", zap.Error(err))
		return
	}

	ctx := context.Background()
	logger.Info("Writing scraped resources into redis")
	// TODO Refactor into dedicated function
	err = rdb.Set(ctx, "Stock", string(b), redis.DataExpirationTTL).Err()
	if err != nil {
		logger.Fatal("Failed to write results into DB", zap.Error(err))
		return
	}

	logger.Info("Scanner finished successfully")
}
