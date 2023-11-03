package main

import (
	"log"
	"time"

	"github.com/reugn/async"
	"github.com/reugn/async/internal/util"
)

func main() {
	future := asyncAction()
	result, err := future.Join()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(result)
}

func asyncAction() async.Future[string] {
	promise := async.NewPromise[string]()
	go func() {
		time.Sleep(time.Second)
		promise.Success(util.Ptr("OK"))
	}()

	return promise.Future()
}
