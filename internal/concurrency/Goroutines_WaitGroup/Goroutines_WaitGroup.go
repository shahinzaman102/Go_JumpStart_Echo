package Goroutines_WaitGroup

import (
	"fmt"
	"sync"
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
