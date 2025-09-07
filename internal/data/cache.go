package data

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

// OrderCache is a global in-memory cache for user orders.
var OrderCache *bigcache.BigCache

// InitCache initializes the global OrderCache using BigCache.
func InitCache() {
	config := bigcache.DefaultConfig(5 * time.Minute) // creates a default config where items expire after 5 minutes.
	config.Shards = 64                                // splits cache into 64 buckets (shards) to reduce lock contention (improves concurrency).
	config.CleanWindow = 1 * time.Minute              // every 1 minute, a background process clears expired items.

	var err error
	ctx := context.Background()                 // you can use a cancellable context if needed
	OrderCache, err = bigcache.New(ctx, config) // creates a new cache with the given configuration.
	if err != nil {                             // If cache initialization fails, the program logs a fatal error and exits.
		log.Fatalf("failed to initialize cache: %v", err)
	}
}

// SetOrdersCache stores orders for a user in the cache.
// `orders` can be any Go struct or slice; it will be JSON-encoded.
func SetOrdersCache(userKey string, orders any) error {
	bytes, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	return OrderCache.Set(userKey, bytes)
}

// GetOrdersCache retrieves cached orders for a user.
// `target` must be a pointer to the expected type (struct or slice).
func GetOrdersCache(userKey string, target any) error {
	entry, err := OrderCache.Get(userKey)
	if err != nil {
		return err
	}
	return json.Unmarshal(entry, target)
}
