package Worker_Pool

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkWorkerPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		numWorkers := 3
		numJobs := 20

		jobs := make(chan Job, numJobs)
		results := make(chan Result, numJobs)

		var wg sync.WaitGroup

		for w := 1; w <= numWorkers; w++ {
			wg.Add(1)
			go worker(w, jobs, results, &wg, r)
		}

		for j := 1; j <= numJobs; j++ {
			jobs <- Job{ID: j, Name: "job"}
		}
		close(jobs)

		go func() {
			wg.Wait()
			close(results)
		}()

		for range results {
			// just drain results
		}
	}
}
