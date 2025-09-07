package Context_Cancellation

import (
	"context"
	"fmt"
	"time"
)

func fetchData(ctx context.Context, userID int, resultChan chan<- string) {
	select {
	case <-time.After(3 * time.Second): // simulate slow API call
		resultChan <- fmt.Sprintf("User %d: data fetched ✅", userID)
	case <-ctx.Done(): // context cancelled or deadline exceeded
		resultChan <- fmt.Sprintf("User %d: fetch cancelled ❌ (%v)", userID, ctx.Err())
	}
}

func Run() {
	resultChan := make(chan string)
	userID := 101

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // -->
	// If we increase timeout to 4 seconds, it will fetch successfully instead.
	// Change 2*time.Second to --> 4*time.Second

	defer cancel() // ensure cleanup

	go fetchData(ctx, userID, resultChan)

	// Wait for result or cancellation
	result := <-resultChan
	fmt.Println(result)

	// Send a hint for learners in the browser
	if ctx.Err() != nil {
		fmt.Println("💡 Hint: Increase timeout to 4*time.Second to fetch successfully ✅")
	}
}
