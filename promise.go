package async

import "sync"

// Promise represents a writable, single-assignment container, which completes a Future.
type Promise interface {

	// Success completes the underlying Future with a value.
	Success(interface{})

	// Failure fails the underlying Future with an error.
	Failure(error)

	// Future returns the underlying Future.
	Future() Future
}

type promiseStatus uint8

const (
	ready promiseStatus = iota
	completed
)

// PromiseImpl implements the Promise interface.
type PromiseImpl struct {
	sync.Mutex
	future Future
	status promiseStatus
}

// NewPromise returns a new PromiseImpl.
func NewPromise() Promise {
	return &PromiseImpl{
		future: NewFuture(),
		status: ready,
	}
}

// Success completes the underlying Future with a given value.
func (p *PromiseImpl) Success(value interface{}) {
	p.Lock()
	defer p.Unlock()

	if p.status != completed {
		p.future.complete(value, nil)
		p.status = completed
	}
}

// Failure fails the underlying Future with a given error.
func (p *PromiseImpl) Failure(e error) {
	p.Lock()
	defer p.Unlock()

	if p.status != completed {
		p.future.complete(nil, e)
		p.status = completed
	}
}

// Future returns the underlying Future.
func (p *PromiseImpl) Future() Future {
	return p.future
}
