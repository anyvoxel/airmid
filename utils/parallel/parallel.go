// Package parallel provide helper function to execute parallel task
package parallel

import (
	"context"

	"github.com/anyvoxel/airmid/utils"
)

// Run will execute workFunc on every idx in count.
//
//nolint:revive,cyclop
func Run(parentCtx context.Context, count int, workFunc func(int) error, opts ...Option) error {
	o := defaultOption()
	opts = append(opts, func(o *option) {
		o.count = count
	})
	for _, opt := range opts {
		opt(o)
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

	wg := utils.WaitGroupWrapper{}
	for i := 0; i < o.concurrent; i++ {
		wg.Wrap(
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
