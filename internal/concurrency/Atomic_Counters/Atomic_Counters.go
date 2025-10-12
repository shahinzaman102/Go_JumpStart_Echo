// Atomic Counters
// - Example: Multiple workers increment a shared counter safely.
// - Shows: using atomic operations to avoid data races without a mutex.
// Race condition: outcome depends on the timing/order of concurrent operations.
// Data race: a specific type of race condition where two goroutines access the same variable
// 			  simultaneously and at least one writes → causes undefined behavior.

// Deadlock: A deadlock happens when two or more goroutines are stuck waiting on each other’s
// 			 resources or locks, so none of them can ever proceed.

package Atomic_Counters

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// worker simulates work and increments the counter atomically.
func worker(id int, wg *sync.WaitGroup, counter *int64, r *rand.Rand) {
	defer wg.Done()

	// Simulate variable work duration
	time.Sleep(time.Duration(r.Intn(400)+100) * time.Millisecond)

	// Atomically increment counter
	atomic.AddInt64(counter, 1)
	// ✅ Use atomic when your shared variable is a simple number or pointer and you only need atomic increments/reads.
	// ✅ Use mutex/rwmutex when you need to protect more complex shared state (like updating a struct, map, or multiple fields at once).

	fmt.Printf("Worker %d finished a task\n", id)
}

// Run executes 10 concurrent workers demonstrating atomic counters.
func Run() {
	// Use a locked source to make rand safe for concurrent use
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var wg sync.WaitGroup
	var completedTasks int64 = 0

	numWorkers := 10

	// Launch 10 workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, &completedTasks, r)
	}

	wg.Wait()

	fmt.Printf("✅ Total tasks completed: %d\n", completedTasks)
}
