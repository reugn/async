package main

import (
	"fmt"
	"sync"

	"github.com/reugn/async"
)

func main() {
	fmt.Println("=== Goroutine ID Example ===")
	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			goroutineID, err := async.GoroutineID()
			if err != nil {
				fmt.Printf("Goroutine %d: failed to get ID: %v\n", id, err)
				return
			}
			fmt.Printf("Goroutine %d: ID = %d\n", id, goroutineID)
		}(i)
	}

	wg.Wait()
}
