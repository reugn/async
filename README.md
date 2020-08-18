<p align="center"><img src="./images/async.png" /></p>
<p align="center">Alternative sync library for Go.</p>
<p align="center">
  <a href="https://travis-ci.org/reugn/async"><img src="https://travis-ci.org/reugn/async.svg?branch=master"></a>
  <a href="https://godoc.org/github.com/reugn/async"><img src="https://godoc.org/github.com/reugn/async?status.svg"></a>
  <a href="https://goreportcard.com/report/github.com/reugn/async"><img src="https://goreportcard.com/badge/github.com/reugn/async"></a>
  <a href="https://codecov.io/gh/reugn/async"><img src="https://codecov.io/gh/reugn/async/branch/master/graph/badge.svg"></a>
</p>

## Overview
* **Future** - A placeholder object for a value that may not yet exist.
* **Promise** - While futures are defined as a type of read-only placeholder object created for a result which doesn’t yet exist, a promise can be thought of as a writable, single-assignment container, which completes a future.
* **Reentrant Lock** - Mutex that allows goroutines to enter into the lock on a resource more than once.
* **Optimistic Lock** - Mutex that allows optimistic reading. Could be retried or switched to RLock in case of failure. Significantly improves performance in case of frequent reads and short writes. See [benchmarks](./benchmarks/README.md).

## Examples
Can be found in the examples directory/tests.

## License
Licensed under the MIT License.