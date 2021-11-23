<div align="center" style="margin:0 !important;"><img src="./docs/images/async.png" /></div>
<div align="center">
  <a href="https://github.com/reugn/async/actions/workflows/build.yml"><img src="https://github.com/reugn/async/actions/workflows/build.yml/badge.svg"></a>
  <a href="https://pkg.go.dev/github.com/reugn/async"><img src="https://pkg.go.dev/badge/github.com/reugn/async"></a>
  <a href="https://goreportcard.com/report/github.com/reugn/async"><img src="https://goreportcard.com/badge/github.com/reugn/async"></a>
  <a href="https://codecov.io/gh/reugn/async"><img src="https://codecov.io/gh/reugn/async/branch/master/graph/badge.svg"></a>
</div>
<br/>
Async provides synchronization and asynchronous computation utilities for Go.

The implemented patterns were taken from Scala and Java.

## Overview
* **Future** - A placeholder object for a value that may not yet exist.
* **Promise** - While futures are defined as a type of read-only placeholder object created for a result which doesn’t yet exist, a promise can be thought of as a writable, single-assignment container, which completes a future.
* **Reentrant Lock** - Mutex that allows goroutines to enter into the lock on a resource more than once.
* **Optimistic Lock** - Mutex that allows optimistic reading. Could be retried or switched to RLock in case of failure. Significantly improves performance in case of frequent reads and short writes. See [benchmarks](./benchmarks/README.md).

### [Go 1.18 Generic prototypes](./generic)
* **Task** - A data type for controlling possibly lazy and asynchronous computations.

## Examples
Can be found in the examples directory/tests.

## License
Licensed under the MIT License.