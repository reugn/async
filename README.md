<div align="center" style="margin:0 !important;"><img src="docs/images/async.png" width="240" /></div>
<div align="center">
  <a href="https://github.com/reugn/async/actions/workflows/build.yml"><img src="https://github.com/reugn/async/actions/workflows/build.yml/badge.svg"></a>
  <a href="https://pkg.go.dev/github.com/reugn/async"><img src="https://pkg.go.dev/badge/github.com/reugn/async"></a>
  <a href="https://goreportcard.com/report/github.com/reugn/async"><img src="https://goreportcard.com/badge/github.com/reugn/async"></a>
  <a href="https://codecov.io/gh/reugn/async"><img src="https://codecov.io/gh/reugn/async/branch/master/graph/badge.svg"></a>
</div>
<br/>

Async provides a comprehensive set of synchronization primitives and asynchronous computation utilities for Go, complementing the standard library with additional concurrency patterns and data structures.

## Features
* **ConcurrentMap** - Implements the generic `async.Map` interface in a thread-safe manner by delegating load/store operations to the underlying `sync.Map`.
* **ShardedMap** - Implements the generic `async.Map` interface in a thread-safe manner, delegating load/store operations to one of the underlying `async.SynchronizedMap`s (shards), using a key hash to calculate the shard number.
* **Future** - A placeholder object for a value that may not yet exist.
* **Promise** - While futures are defined as a type of read-only placeholder object created for a result which doesn’t yet exist, a promise can be thought of as a writable, single-assignment container, which completes a future.
* **Executor** - A worker pool for executing asynchronous tasks, where each submission returns a Future instance representing the result of the task.
* **Task** - A data type for controlling possibly lazy and asynchronous computations.
* **Once** - An object similar to sync.Once having the Do method taking `f func() (T, error)` and returning `(T, error)`.
* **Value** - An object similar to atomic.Value, but without the consistent type constraint.
* **CyclicBarrier** - A reusable synchronization primitive that allows a group of goroutines to wait for each other to reach a common barrier point.
* **WaitGroupContext** - A WaitGroup with the `context.Context` support for graceful unblocking.
* **ReentrantLock** - A mutex that allows goroutines to enter into the lock on a resource more than once.
* **PriorityLock** - A non-reentrant mutex that allows for the specification of lock acquisition priority.

## Examples
Runnable examples are available in the [examples](./examples) directory. See the test files for additional usage examples.

## License
Licensed under the MIT License.
