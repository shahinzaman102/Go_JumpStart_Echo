// Goroutines + WaitGroup
// - Example: Multiple users making payment requests.
// - Shows: spawning goroutines, waiting for completion.
// spawning goroutines: means creating and starting new lightweight threads of execution by prefixing a function call with the keyword go.

package Goroutines_WaitGroup

import (
	"fmt"
	"sync" // see about this at below -->
	"time"
)

// processPayment simulates a payment processing task with panic recovery.
func processPayment(userID int, amount float64, wg *sync.WaitGroup) {
	defer wg.Done() // mark this goroutine as done

	// Recover from panic (invalid payment)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("User %d: recovered from failure: %v\n", userID, r)
		}
	}()

	fmt.Printf("User %d: Connecting to payment gateway...\n", userID)

	// Simulate invalid payment
	if amount <= 0 {
		panic("invalid payment amount")
	}

	// Simulate processing time
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("User %d: Payment processed successfully: $%.2f\n", userID, amount)
}

// Run demonstrates using WaitGroup to wait for multiple goroutines.
func Run() {
	var wg sync.WaitGroup

	payments := []struct {
		userID int
		amount float64
	}{
		{1, 0},   // invalid
		{2, 50},  // valid
		{3, -10}, // invalid
		{4, 100}, // valid
	}

	for _, p := range payments {
		wg.Add(1)
		go processPayment(p.userID, p.amount, &wg)
	}

	wg.Wait()
	fmt.Println("All payment requests handled ✅")
}

// ----------------------------------------------------
// ✅ Primitives → the basic building blocks provided by a language or library to control concurrency (e.g., Mutex, RWMutex, WaitGroup).
//     Think of them as the "raw ingredients."
// ✅ A package like sync is called a concurrency package or concurrency toolkit →
// 	it provides these primitives (tools) for managing goroutines safely.

// So:
//  - Primitives = the tools (Mutex, RWMutex, etc.).
//  - sync = the package (toolbox) that gives you those tools.
// ----------------------------------------------------
// Key tools inside sync:

//  - sync.Mutex → mutual exclusion lock (one goroutine at a time).
//  - sync.RWMutex → read/write lock (many readers OR one writer).
//  - sync.WaitGroup → wait until a group of goroutines finish.
//  - sync.Once → ensures something runs only once (even across goroutines).
//  - sync.Cond → condition variable (goroutines wait for signals).
//  - sync.Map → a concurrency-safe map (no need for manual locking).
//  - sync.Pool → object reuse to reduce allocations (performance).

// 👉 In short: sync = Go’s toolbox for concurrency control.
