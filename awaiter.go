package awaiter

import (
	"sync"
)

type awaiter struct {
	wg     *sync.WaitGroup
	cancel chan interface{}
}

// Awaiter is similar to a WaitGroup but simplifies
// the resource handling
type Awaiter interface {
	// AwaitSync blocks until all routines started with "Go"
	// are finished (same as WaitGroup.Wait)
	AwaitSync()

	// Go executes "impl" in a go routine and increments
	// the internal WaitGroup counter by one
	Go(impl func())

	// Cancel signals the executed go routines to cancel
	// their execution.
	// After Cancel was called, IsCancelRequested returns true and the channel
	// returned by CancelRequested is closed
	Cancel()

	// IsCancelRequested returns true if Cancel was called
	IsCancelRequested() bool

	// CancelRequested returns a channel than never provides data
	// but will be closed when Cancel was called
	CancelRequested() <-chan interface{}
}

// New creates a new Awaiter instance
func New() Awaiter {
	a := new(awaiter)
	a.cancel = make(chan interface{})
	a.wg = new(sync.WaitGroup)
	return a
}

func (a *awaiter) Go(impl func()) {
	a.wg.Add(1)
	go func(a *awaiter) {
		defer a.wg.Done()
		impl()
	}(a)
}

func (a *awaiter) AwaitSync() {
	a.wg.Wait()
}

func (a *awaiter) Cancel() {
	select {
	case <-a.cancel:
		return
	default:
		close(a.cancel)
	}
}

func (a *awaiter) IsCancelRequested() bool {
	select {
	case <-a.cancel:
		return true
	default:
		return false
	}
}

func (a *awaiter) CancelRequested() <-chan interface{} {
	return a.cancel
}
