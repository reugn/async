# async
Alternative sync library for Go 

## Overview
* Future - A placeholder object for a value that may not yet exist.
* Promise - While futures are defined as a type of read-only placeholder object created for a result which doesn’t yet exist, a promise can be thought of as a writable, single-assignment container, which completes a future.
* Reentrant Lock - Mutex that allows goroutines to enter into lock on a resource more than once.
* Optimistic Lock - Mutex that allows optimistic reading. Could be retried or switched to RLock in case of failure.

## Examples
Could be found in the examples directory.

## License
Licensed under the MIT License.