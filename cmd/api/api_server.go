package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	inven    *inventory.Inventory
	router   *gin.Engine
	apiURL   string
	dbURL    string
	dbPass   string
	rdb      *redis.Client
	ctx      context.Context
	redisKey = "Stock"
	logger   *zap.Logger
)

func init() {
	// Logging config
	logger, _ = zap.NewProduction()

	// Getting config
	apiHost := os.Getenv("CIQ_API_HOST")
	apiPort := os.Getenv("CIQ_API_PORT")
	dbHost := os.Getenv("CIQ_DB_HOST")
	dbPort := os.Getenv("CIQ_DB_PORT")
	dbPass = os.Getenv("CIQ_DB_PASS")
	apiURL = fmt.Sprintf("%s:%s", apiHost, apiPort)
	dbURL = fmt.Sprintf("%s:%s", dbHost, dbPort)

	// Initializaion global vars
	inven = inventory.NewInventory()
	router = gin.Default()
	// Configure GIN to use ZAP
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
}

func addHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

// getAccounts returns every account in Stock
func getAccountsCount(c *gin.Context) {
	logger.Debug("Retrieving accounts count")
	updateStock()
	c.PureJSON(http.StatusOK, len(inven.Accounts))
}

// getAccounts returns every account in Stock
func getClusters(c *gin.Context) {
	logger.Debug("Retrieving complete cluster inventory")
	updateStock()
	addHeaders(c)

	var clusters []inventory.Cluster
	for _, account := range inven.Accounts {
		for _, cluster := range account.Clusters {
			clusters = append(clusters, *cluster)
		}
	}

	c.PureJSON(http.StatusOK, clusters)
}

// getAccounts returns every account in Stock
func getAccounts(c *gin.Context) {
	logger.Debug("Retrieving complete accounts inventory")
	updateStock()
	addHeaders(c)
	c.PureJSON(http.StatusOK, inven.Accounts)
}

// getAccountsByName returns an account by its name in Stock
func getAccountsByName(c *gin.Context) {
	logger.Debug("Retrieving accounts by name")
	updateStock()
	addHeaders(c)
	name := c.Param("name")

	account, ok := inven.Accounts[name]
	if ok {
		c.PureJSON(http.StatusOK, account)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// updateStock updates the cache of the API
func updateStock() {
	// Getting Redis Results
	val, err := rdb.Get(ctx, redisKey).Result()
	if err != nil {
		logger.Error("Can't update Inventory", zap.Error(err))
	}

	// Unmarshall from JSON to inventory.Inventory type
	json.Unmarshal([]byte(val), &inven)
}

func main() {
	defer logger.Sync()
	logger.Info("Starting Openshift Inventory API")
	logger.Info("API URL: ", zap.String("API-URL", apiURL))
	logger.Info("DB URL: ", zap.String("DB-URL", dbURL))

	// Preparing API Endpoints
	router.GET("/accounts", getAccounts)
	router.GET("/accounts/:name", getAccountsByName)
	router.GET("/accountsCount", getAccountsCount)
	router.GET("/clusters", getClusters)
	router.GET("/mock", getMockCluster)

	// RedisDB connection
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     dbURL,
		Password: dbPass,
		DB:       0,
	})

	// Start API
	logger.Info("API Ready to serve")
	router.Run(apiURL)
}
