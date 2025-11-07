package utils

import (
	"sync"
)

// WaitGroupWrapper is a help struct for sync.WaitGroup.
type WaitGroupWrapper struct {
	NoCopy
	sync.WaitGroup
}

// Wrap will handle sync.WaitGroup Add and Done.
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		defer w.Done()

		cb()
	}()
}
