package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cluster-iq/pkg/inventory"
	"github.com/cluster-iq/pkg/stocker"
	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
)

var (
	inven     inventory.Inventory
	stockers  []stocker.Stocker
	dbURL     string
	dbPass    string
	credsFile string
	rdb       *redis.Client
	ctx       context.Context
)

func init() {
	// Getting config
	dbHost := os.Getenv("CIQ_DB_HOST")
	dbPort := os.Getenv("CIQ_DB_PORT")
	dbPass = os.Getenv("CIQ_DB_PASSWORD")
	credsFile = os.Getenv("CIQ_CLOUD_CREDS")
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
		log.Fatal("Can't Open credentials file", err.Error())
	}

	// Reading INI file content
	for _, account := range cfg.Sections() {
		newAccount := inventory.NewAccount(
			account.Name(),
			getProvider(account.Key("provider").String()),
			account.Key("user").String(),
			account.Key("key").String(),
		)
		accounts = append(accounts, newAccount)
	}

	return accounts
}

func createStockers(accounts []inventory.Account) error {
	for _, account := range accounts {
		switch account.Provider {
		case inventory.AWSProvider:
			log.Printf("Adding AWS account '%s' to be inventored\n", account.Name)
			stockers = append(stockers, stocker.NewAWSStocker(account))
		case inventory.GCPProvider:
			err := fmt.Errorf("Google Cloud Platform (GCP) Stocker not implemented! Account %s will not be scanned", account.Name)
			log.Println(err)
			return err
		case inventory.AzureProvider:
			err := fmt.Errorf("Microsoft Azure Stocker not implemented! Account %s will not be scanned", account.Name)
			log.Println(err)
			return err
		}
	}

	return nil
}

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
	// Prepare New Stock
	inven = inventory.NewInventory()

	// Get Cloud Accounts from credentials file
	accounts := GetCloudProviderAccounts()

	// Running Stockers
	createStockers(accounts)
	err := startStockers()
	if err != nil {
		log.Fatal(err)
	}

	inven.PrintInventory()
	fmt.Println("Stock Maker Finished")

	// TODO move redis module
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     dbURL,
		Password: dbPass,
		DB:       0,
	})

	b, err := json.Marshal(inven)
	if err != nil {
		return
	}

	log.Println("Writing results into Redis...")
	err = rdb.Set(ctx, "Stock", string(b), 0).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Done!")
}
