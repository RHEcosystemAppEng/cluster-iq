package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RHEcosystemAppEng/cluster-iq/pkg/inventory"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	inven    inventory.Inventory
	router   *gin.Engine
	apiURL   string
	dbURL    string
	dbPass   string
	rdb      *redis.Client
	ctx      context.Context
	redisKey = "Stock"
)

func init() {
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
}

func addHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

// getAccounts returns every account in Stock
func getAccountsCount(c *gin.Context) {
	updateStock()
	c.PureJSON(http.StatusOK, len(inven.Accounts))
}

// getAccounts returns every account in Stock
func getClusters(c *gin.Context) {
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
	updateStock()
	addHeaders(c)
	c.PureJSON(http.StatusOK, inven.Accounts)
}

// getAccountsByName returns an account by its name in Stock
func getAccountsByName(c *gin.Context) {
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
		fmt.Println(err)
	}

	// Unmarshall from JSON to inventory.Inventory type
	json.Unmarshal([]byte(val), &inven)
}

func main() {
	log.Println("Starting Openshift Inventory API")
	log.Println("API URL: ", apiURL)
	log.Println("DB URL: ", dbURL)

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
	log.Println("API Ready to serve")
	router.Run(apiURL)
}
