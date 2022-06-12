package main

import (
	"log"
	"time"

	"github.com/reugn/async"
)

func main() {
	future := asyncAction()
	result, err := future.Get()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(result)
}

func asyncAction() async.Future[string] {
	promise := async.NewPromise[string]()
	go func() {
		time.Sleep(time.Second)
		promise.Success("OK")
	}()

	return promise.Future()
}
