// Copyright (c) 2025 The anyvoxel Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package parallel provide helper function to execute parallel task
package parallel

import (
	"context"
	"sync"

	"github.com/anyvoxel/airmid/anvil"
)

// Run will execute workFunc on every idx in count.
//
//nolint:revive,cyclop
func Run(parentCtx context.Context, count int, workFunc func(int) error, opts ...anvil.Option[parallelOption]) error {
	o := defaultOption()
	applyOpts := make([]anvil.Option[parallelOption], len(opts)+1)
	applyOpts[0] = anvil.NewFnOption(func(o *parallelOption) {
		o.count = count
	})
	copy(applyOpts[1:], opts)

	for _, opt := range applyOpts {
		opt.Apply(o)
	}
	if err := o.Complete(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	// NOTE:
	// 1. optimize this when count is too large
	// 2. +1 to make sure the send never hang
	errCh := make(chan error, o.count+1)

	idxCh := make(chan int, o.count+1)
	for i := 0; i < o.count; i++ {
		idxCh <- i
	}
	close(idxCh)

	wg := sync.WaitGroup{}
	for i := 0; i < o.concurrent; i++ {
		wg.Go(
			func() {
				for idx := range idxCh {
					select {
					case <-ctx.Done():
						return
					default:
					}

					err := workFunc(idx)
					if err != nil {
						errCh <- err
						cancel()
						return
					}
				}
			},
		)
	}

	// We must wait on another goroutine, to avoid the workFunc block
	go func() {
		wg.Wait()
		cancel()
		close(errCh)
	}()

	// At there, we have such scenario:
	//  1. all workFunc have been done
	//  2. some workFunc cancel
	//  3. parentCtx Done
	// So we must distinct it.
	<-ctx.Done()

	select {
	case err := <-errCh:
		// NOTE: we must check err again, because when we read the errCh,
		// It maybe closed by another goroutine.
		if err != nil {
			return err
		}
	default:
	}

	// If the parentCtx is cancel, we must return err to caller
	return parentCtx.Err()
}
