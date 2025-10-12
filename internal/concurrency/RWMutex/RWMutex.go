// RWMutex (read-heavy workload)
// 👉 Real-world analogy: a product catalog service (like Amazon search).
// - Thousands of users are reading product prices.
// - Very few operations are updating prices.
// - RWMutex allows all readers to run concurrently, which is crucial at scale.

// Use RWMutex → when reads are frequent and writes are rare. For examples -->
// ---------------------------------------
// - Flight schedule lookup → many people checking flight times, few updates from the airline.
// - E-commerce product catalog → thousands of customers viewing items, few price/stock changes.
// - Configuration values (feature flags, app settings) → often read, rarely updated.
// - Weather or news feeds → many users read, system occasionally updates.
// - Dictionary / cache service → lots of lookups, rare invalidations/updates.
// ---------------------------------------
// RWMutex
// - Example: Product catalog → many readers fetch prices, a few writers update prices.
// - Shows: multiple readers can access simultaneously, but writers get exclusive access (read–write lock).

package rwmutex

import (
	"fmt"
	"sync"
	"time"
)

// ProductCatalog simulates a catalog where reads are frequent and writes are rare.
type ProductCatalog struct {
	mu     sync.RWMutex
	prices map[string]int
}

// UpdatePrice updates the price of a product. Writer lock is exclusive.
func (c *ProductCatalog) UpdatePrice(product string, price int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prices[product] = price
	fmt.Printf("Updated %s → %d\n", product, price)
}

// GetPrice returns the price of a product. Readers can share the lock.
func (c *ProductCatalog) GetPrice(product string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.prices[product]
}

// Run demonstrates multiple readers and few writers using RWMutex.
func Run() {
	catalog := &ProductCatalog{
		prices: map[string]int{
			"laptop": 1000,
			"phone":  500,
		},
	}

	var wg sync.WaitGroup

	// Many readers
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id%5) * 50 * time.Millisecond) // stagger readers
			fmt.Printf("Reader %d → Laptop Price: %d\n", id, catalog.GetPrice("laptop"))
		}(i) // pass loop variable to avoid closure capture issues
	}

	// Few writers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id) * 200 * time.Millisecond) // stagger writers
			catalog.UpdatePrice("laptop", 900+id*10)               // simulate price updates
		}(i)
	}

	wg.Wait()
	fmt.Println("All reads and writes complete ✅")
}
