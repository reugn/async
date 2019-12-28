package main

import (
	"fmt"
	"sync"

	"github.com/reugn/async"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			id, _ := async.GoroutineID()
			fmt.Println(id)
			wg.Done()
		}()
	}
	wg.Wait()
}
