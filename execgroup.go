package execgroup

import (
	"errors"
	"fmt"
	"sync"
)

// ExecGroup uses a WaitGroup to allow multiple functions to be called
// concurrently.
type ExecGroup struct {
	wg sync.WaitGroup
	me MultiError
	m  sync.Mutex
}

// Do executes the provided function in a goroutine.
func (e *ExecGroup) Do(fn func()) {
	e.wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					e.TrackError(err)
				} else {
					msg := fmt.Sprintf("panic in execgroup: %s", r)
					e.TrackError(errors.New(msg))
				}
			}
			e.wg.Done()
		}()
		fn()
	}()
}

// TrackError tracks an error that occurred.
func (e *ExecGroup) TrackError(err error) {
	if err != nil {
		e.m.Lock()
		defer e.m.Unlock()
		e.me = e.me.Append(err)
	}
}

// Wait for all functions to complete.
func (e *ExecGroup) Wait() error {
	e.wg.Wait()
	if len(e.me) == 0 {
		return nil
	}
	return e.me
}

// Error returns a MultiError if any errors occurred.
func (e *ExecGroup) Error() error {
	if len(e.me) != 0 {
		return e.me
	}
	return nil
}
