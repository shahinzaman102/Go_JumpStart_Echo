package Channels_unbuffered

import (
	"fmt"
	"time"
)

// Producer: simulates customers placing orders
func takeOrders(orderCh chan string) {
	orders := []string{"Burger", "Pizza", "Pasta", "Salad"}
	for _, order := range orders {
		fmt.Println("Customer ordered:", order)
		orderCh <- order // send order to channel (blocks until chef receives)
		time.Sleep(500 * time.Millisecond)
	}
	close(orderCh) // signal no more orders
}

// Consumer: simulates chef preparing orders
func prepareOrders(orderCh chan string) {
	for order := range orderCh {
		fmt.Println("Chef received:", order)
		time.Sleep(1 * time.Second) // simulate cooking time
		fmt.Println("Chef prepared:", order)
	}
}

func Run() {
	orderCh := make(chan string) // unbuffered channel
	// rderCh := make(chan string, 2) // buffered channel

	// Run producer and consumer concurrently
	go takeOrders(orderCh)
	prepareOrders(orderCh) // main goroutine acts as consumer
}
