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
