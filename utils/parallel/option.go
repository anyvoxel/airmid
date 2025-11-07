package parallel

import (
	"runtime"

	"github.com/anyvoxel/airmid/xerrors"
)

// Option is the configuration helper for parallel.
type Option = func(*option)

type option struct {
	// concurrent is the parallelism of parallel
	concurrent int

	// count is the total number of task
	count int
}

func defaultOption() *option {
	return &option{
		concurrent: runtime.NumCPU(),
		count:      0,
	}
}

// Complete will validate the option and correct it.
func (o *option) Complete() error {
	if o.count <= 0 {
		return xerrors.Errorf("option.Complete: count '%d' must greater than zero", o.count)
	}

	if o.concurrent <= 0 {
		o.concurrent = runtime.NumCPU()
	}

	// We should use at most count workers
	if o.count < o.concurrent {
		o.concurrent = o.count
	}
	return nil
}

// WithConcurrent set the concurrent of parallel.
func WithConcurrent(v int) Option {
	return func(o *option) {
		o.concurrent = v
	}
}
