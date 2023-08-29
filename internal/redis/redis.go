package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DatabaseConnectionTimeout = 10 * time.Second
)

func InitDatabase(address, password string) (*redis.Client, error) {
	rdc := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), DatabaseConnectionTimeout)
	defer cancel()

	if _, err := rdc.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	return rdc, nil
}
