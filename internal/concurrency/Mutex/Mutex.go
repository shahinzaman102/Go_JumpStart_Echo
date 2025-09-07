package mutex

import (
	"fmt"
	"sync"
)

// TicketSystem simulates a ticket booking system with a mutex protecting the seat count.
type TicketSystem struct {
	mu    sync.Mutex
	seats int
}

// Book attempts to reserve n seats.
func (t *TicketSystem) Book(n int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.seats >= n {
		t.seats -= n
		fmt.Printf("Booked %d seats → Remaining: %d\n", n, t.seats)
	} else {
		fmt.Printf("Booking failed for %d seats → Remaining: %d\n", n, t.seats)
	}
}

// Cancel adds n seats back to the system.
func (t *TicketSystem) Cancel(n int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.seats += n
	fmt.Printf("Canceled %d seats → Remaining: %d\n", n, t.seats)
}

// CheckAvailableSeats returns the current number of available seats.
func (t *TicketSystem) CheckAvailableSeats() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.seats
}

// Run demonstrates concurrent bookings and cancellations.
func Run() {
	system := &TicketSystem{seats: 50}
	var wg sync.WaitGroup

	// Many writers (book/cancel)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if id%2 == 0 {
				system.Book(2) // even IDs → book 2 seats
			} else {
				system.Cancel(1) // odd IDs → cancel 1 seat
			}
		}(i) // pass loop variable to avoid closure capture issues
	}

	// Few readers
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("Seats Available → %d\n", system.CheckAvailableSeats())
		}()
	}

	wg.Wait()
	fmt.Printf("Final Seats → %d\n", system.CheckAvailableSeats())
}
