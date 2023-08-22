package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cluster-iq/pkg/stock"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Prepare New Stock
	cloudStock := stock.NewStock()

	// Get Cloud Accounts from credentials file
	fmt.Println("Starting")

	// TODO move redis module
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	b, err := json.Marshal(cloudStock)
	if err != nil {
		return
	}

	err = rdb.Set(ctx, "Stock", string(b), 0).Err()
	if err != nil {
		panic(err)
	}

	//val, err := rdb.Get(ctx, "Stock").Result()
	//if err != nil {
	//	fmt.Println(err)
	//}

	//fmt.Println(val)
}
