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

package parallel

import (
	"runtime"

	"github.com/anyvoxel/airmid/anvil/xerrors"
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
