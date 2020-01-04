<p align="center">
  <img src="./images/async.png" />
</p>

[![Build Status](https://travis-ci.org/reugn/async.svg?branch=master)](https://travis-ci.org/reugn/async)
[![GoDoc](https://godoc.org/github.com/reugn/async?status.svg)](https://godoc.org/github.com/reugn/async)
[![Go Report Card](https://goreportcard.com/badge/github.com/reugn/async)](https://goreportcard.com/report/github.com/reugn/async)
[![codecov](https://codecov.io/gh/reugn/async/branch/master/graph/badge.svg)](https://codecov.io/gh/reugn/async)

Alternative sync library for Go.

## Overview
* **Future** - A placeholder object for a value that may not yet exist.
* **Promise** - While futures are defined as a type of read-only placeholder object created for a result which doesn’t yet exist, a promise can be thought of as a writable, single-assignment container, which completes a future.
* **Reentrant Lock** - Mutex that allows goroutines to enter into the lock on a resource more than once.
* **Optimistic Lock** - Mutex that allows optimistic reading. Could be retried or switched to RLock in case of failure. Significantly improves performance in case of frequent reads and short writes. See [benchmarks](./benchmarks/README.md).

## Examples
Could be found in the examples directory/tests.

## License
Licensed under the MIT License.