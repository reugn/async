package main

import (
	"fmt"
	"time"

	"github.com/reugn/async"
)

func main() {
	future := asyncAction()
	rt, _ := future.Get()
	fmt.Println(rt)
}

func asyncAction() async.Future {
	promise := async.NewPromise()
	go func() {
		time.Sleep(time.Second)
		promise.Success("OK")
	}()
	return promise.Future()
}
