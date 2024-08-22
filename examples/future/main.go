package main

import (
	"context"
	"log"
	"time"

	"github.com/reugn/async"
)

const ok = "OK"

func main() {
	// using a promise
	future1 := asyncAction()
	result1, err := future1.Join()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(result1)

	// using a task
	task := async.NewTask(func() (string, error) { return ok, nil })
	result2, err := task.Call().Join()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(result2)

	// using the executor
	ctx := context.Background()
	executor := async.NewExecutor[*string](ctx, async.NewExecutorConfig(2, 2))

	future3, err := executor.Submit(func(_ context.Context) (*string, error) {
		value := ok
		return &value, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	result3, err := future3.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(*result3)
}

func asyncAction() async.Future[string] {
	promise := async.NewPromise[string]()
	go func() {
		time.Sleep(time.Second)
		promise.Success(ok)
	}()

	return promise.Future()
}
