<div align="center" style="margin:0 !important;"><img src="./docs/images/async.png" /></div>
<div align="center">
  <a href="https://github.com/reugn/async/actions/workflows/build.yml"><img src="https://github.com/reugn/async/actions/workflows/build.yml/badge.svg"></a>
  <a href="https://pkg.go.dev/github.com/reugn/async"><img src="https://pkg.go.dev/badge/github.com/reugn/async"></a>
  <a href="https://goreportcard.com/report/github.com/reugn/async"><img src="https://goreportcard.com/badge/github.com/reugn/async"></a>
  <a href="https://codecov.io/gh/reugn/async"><img src="https://codecov.io/gh/reugn/async/branch/master/graph/badge.svg"></a>
</div>
<br/>
Async is a synchronization and asynchronous computation package for Go.

## Overview
* **ConcurrentMap** - Implements the generic `async.Map` interface in a thread-safe manner by delegating load/store operations to the underlying `sync.Map`.
* **Future** - A placeholder object for a value that may not yet exist.
* **Promise** - While futures are defined as a type of read-only placeholder object created for a result which doesn’t yet exist, a promise can be thought of as a writable, single-assignment container, which completes a future.
* **Task** - A data type for controlling possibly lazy and asynchronous computations.
* **Once** - An object similar to sync.Once having the Do method taking `f func() (T, error)` and returning `(T, error)`.
* **Value** - An object similar to atomic.Value, but without the consistent type constraint.
* **WaitGroupContext** - A WaitGroup with the `context.Context` support for graceful unblocking.
* **ReentrantLock** - A Mutex that allows goroutines to enter into the lock on a resource more than once.

## Examples
Can be found in the examples directory/tests.

## License
Licensed under the MIT License.
