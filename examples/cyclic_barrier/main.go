package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/reugn/async"
)

func main() {
	fmt.Println("=== CyclicBarrier Example ===")
	barrier := async.NewCyclicBarrier(3)
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d: starting work\n", id)
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Goroutine %d: reached barrier\n", id)
			if err := barrier.Await(); err != nil {
				fmt.Printf("Goroutine %d: barrier error: %v\n", id, err)
				return
			}
			fmt.Printf("Goroutine %d: passed barrier\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All goroutines completed")

	// Example with context
	fmt.Println("\n=== CyclicBarrier with Context ===")
	barrier2 := async.NewCyclicBarrier(2)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	wg.Add(2)

	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		if err := barrier2.AwaitContext(ctx); err != nil {
			fmt.Printf("Goroutine 1: %v\n", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := barrier2.AwaitContext(ctx); err != nil {
			fmt.Printf("Goroutine 2: %v\n", err)
		}
	}()

	wg.Wait()
}
