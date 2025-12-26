package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/reugn/async"
	"github.com/reugn/async/internal/ptr"
)

func main() {
	runPromiseExample()
	runTaskExample()
	runExecutorExample()
	runTransformationsExample()
	runRecoveryExample()
}

func runPromiseExample() {
	fmt.Println("=== Promise Example ===")
	future := asyncAction()
	result, err := future.Join()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %s\n", result)
}

func runTaskExample() {
	fmt.Println("\n=== Task Example ===")
	task := async.NewTask(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "Task completed", nil
	})
	result, err := task.Call().Join()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %s\n", result)
}

func runExecutorExample() {
	fmt.Println("\n=== Executor Example ===")
	ctx := context.Background()
	executor := async.NewExecutor[*string](ctx, async.NewExecutorConfig(2, 2))

	future, err := executor.Submit(func(_ context.Context) (*string, error) {
		value := "Executor task completed"
		return ptr.Of(value), nil
	})
	if err != nil {
		if shutdownErr := executor.Shutdown(); shutdownErr != nil {
			log.Printf("Error shutting down executor: %v", shutdownErr)
		}
		log.Fatal(err)
	}

	result, err := future.Get(ctx)
	if err != nil {
		if shutdownErr := executor.Shutdown(); shutdownErr != nil {
			log.Printf("Error shutting down executor: %v", shutdownErr)
		}
		log.Fatal(err)
	}
	fmt.Printf("Result: %s\n", *result)

	if err := executor.Shutdown(); err != nil {
		log.Printf("Error shutting down executor: %v", err)
	}
}

func runTransformationsExample() {
	fmt.Println("\n=== Future Transformations ===")
	promise := async.NewPromise[int]()
	go func() {
		time.Sleep(50 * time.Millisecond)
		promise.Success(10)
	}()

	transformed := promise.Future().
		Map(func(v int) (int, error) {
			return v * 2, nil
		}).
		FlatMap(func(v int) (async.Future[int], error) {
			p := async.NewPromise[int]()
			go func() {
				time.Sleep(50 * time.Millisecond)
				p.Success(v + 5)
			}()
			return p.Future(), nil
		})

	result, err := transformed.Join()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transformed result: %d\n", result)
}

func runRecoveryExample() {
	fmt.Println("\n=== Error Recovery ===")
	failingPromise := async.NewPromise[int]()
	go func() {
		time.Sleep(50 * time.Millisecond)
		failingPromise.Failure(errors.New("operation failed"))
	}()

	recovered := failingPromise.Future().Recover(func() (int, error) {
		return 10, nil // fallback value
	})

	result, err := recovered.Join()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recovered result: %d\n", result)
}

func asyncAction() async.Future[string] {
	promise := async.NewPromise[string]()
	go func() {
		time.Sleep(100 * time.Millisecond)
		promise.Success("Promise completed")
	}()
	return promise.Future()
}
