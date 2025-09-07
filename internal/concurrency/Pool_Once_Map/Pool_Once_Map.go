package Pool_Once_Map

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

// --- sync.Pool (buffer reuse for efficiency) ---
var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer) // allocate when empty
	},
}

// --- sync.Once (init something only once) ---
var initOnce sync.Once

func initDB() {
	fmt.Println("Initializing database connection...")
	time.Sleep(time.Second) // simulate heavy setup
	fmt.Println("Database connection ready ✅")
	fmt.Println("")
}

// --- sync.Map (concurrent safe map for user sessions) ---
var sessionCache sync.Map

func handleRequest(userID int, data string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Ensure DB initialized only once
	initOnce.Do(initDB)

	// Get a buffer from pool
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset() // always reset before use
	defer bufPool.Put(buf)

	fmt.Fprintf(buf, "User %d data: %s", userID, data)

	// Save to session cache
	sessionCache.Store(userID, buf.String())

	fmt.Printf("Processed request for user %d\n", userID)
}

func Run() {
	var wg sync.WaitGroup

	// Simulate multiple requests
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go handleRequest(i, fmt.Sprintf("payload-%d", i), &wg)
	}

	wg.Wait()

	fmt.Println("\n--- Sessions in Cache ---")
	sessionCache.Range(func(key, value any) bool {
		fmt.Printf("User %v → %v\n", key, value)
		return true
	})
}
