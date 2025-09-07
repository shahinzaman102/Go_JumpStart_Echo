package Buffered_Channels

import (
	"fmt"
	"time"
)

// Log writer: consumes logs from channel
func logWriter(logCh chan string) {
	for log := range logCh {
		fmt.Println("Writing log:", log)
		time.Sleep(800 * time.Millisecond) // simulate slow I/O
	}
}

func Run() {
	// Buffered channel with capacity 3
	logCh := make(chan string, 3)

	// Start log writer in a separate goroutine
	go logWriter(logCh)

	// Multiple producers sending logs
	for i := 1; i <= 6; i++ {
		logMsg := fmt.Sprintf("Log message %d", i)
		fmt.Println("Generated:", logMsg)
		logCh <- logMsg                    // won't block until buffer is full
		time.Sleep(200 * time.Millisecond) // simulate app activity
	}

	close(logCh)                // no more logs
	time.Sleep(3 * time.Second) // wait for writer to finish
}
