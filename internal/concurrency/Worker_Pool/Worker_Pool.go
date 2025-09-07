package Worker_Pool

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents an image processing task
type Job struct {
	ID   int
	Name string
}

// Result represents the output of processing
type Result struct {
	JobID   int
	Outcome string
}

// worker consumes jobs and sends results
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, r *rand.Rand) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d: processing job %d (%s)\n", id, job.ID, job.Name)

		// Simulate processing time
		time.Sleep(time.Duration(r.Intn(500)+200) * time.Millisecond)

		results <- Result{
			JobID:   job.ID,
			Outcome: fmt.Sprintf("Processed %s by worker %d", job.Name, id),
		}
	}
}

func Run() {
	// Create a local random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numWorkers := 3
	numJobs := 8

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg, r)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Name: fmt.Sprintf("image_%d.png", j)}
	}
	close(jobs) // no more jobs

	// Wait for all workers to finish in background
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for res := range results {
		fmt.Printf("Result: Job %d → %s\n", res.JobID, res.Outcome)
	}

	fmt.Println("All jobs processed.")
}
