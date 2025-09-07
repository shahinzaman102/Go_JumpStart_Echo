package Cond_syncCond

import (
	"fmt"
	"sync"
	"time"
)

var (
	queue []int
	mutex sync.Mutex
	cond  = sync.NewCond(&mutex)
)

func consumer(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for { // outer infinite loop - it'll continuously run by ->
		// cond.Signal() & cond.Wait()
		mutex.Lock()
		for len(queue) == 0 { // inner loop: wait while queue is empty
			// Wait until producer signals
			cond.Wait()
		}

		// Consume item
		item := queue[0]
		queue = queue[1:]
		mutex.Unlock()

		// Stop condition: exit if item == -1 (poison pill)
		if item == -1 {
			return
		}

		fmt.Printf("Consumer %d consumed item %d\n", id, item)
	}
}

func producer(count int) {
	for i := 1; i <= count; i++ {
		mutex.Lock()
		queue = append(queue, i)
		mutex.Unlock()

		// Signal one waiting consumer
		cond.Signal()
		time.Sleep(200 * time.Millisecond)
	}

	for range []int{0, 1} { // send stop signals for 2 consumers
		mutex.Lock()
		queue = append(queue, -1) // poison pill
		mutex.Unlock()
		cond.Signal()
	}
}

func Run() {
	var wg sync.WaitGroup

	// Two consumers
	wg.Add(2)
	go consumer(1, &wg)
	go consumer(2, &wg)

	// One producer
	go producer(5)

	wg.Wait()
	fmt.Println("✅ All consumers finished.")
}
