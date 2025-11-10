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
	"context"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

func TestRun(t *testing.T) {
	type testCase struct {
		desp     string
		count    int
		workFunc func(int) error
		opts     []Option
		err      string
	}
	testCases := []testCase{
		{
			desp:     "option error",
			count:    0,
			workFunc: nil,
			err:      "count '0'",
		},
		{
			desp:  "all success",
			count: 5,
			workFunc: func(_ int) error {
				return nil
			},
			opts: []Option{
				WithConcurrent(2),
			},
			err: "",
		},
		{
			desp:  "some failed",
			count: 10,
			workFunc: func(i int) error {
				if i == 7 {
					return xerrors.ErrContinue
				}
				return nil
			},
			opts: []Option{
				WithConcurrent(4),
			},
			err: "Continue",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := Run(context.Background(), tc.count, tc.workFunc, tc.opts...)

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestRunWorkFuncNotBlock(t *testing.T) {
	g := NewWithT(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count := int32(0)
	now := time.Now()
	err := Run(ctx, 10, func(_ int) error {
		time.Sleep(1000 * time.Second)
		atomic.AddInt32(&count, 1)
		return nil
	})

	g.Expect(time.Since(now).Milliseconds() - 10000).To(BeNumerically("<=", 100))
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(0)))
	g.Expect(err).To(HaveOccurred())
}

// this unit test is used to avoid 'send on closed channel'.
func TestSendOnClosedErrChannel(_ *testing.T) {
	funcs := []func(i int) error{
		func(_ int) error {
			time.Sleep(1 * time.Second)
			return xerrors.ErrContinue
		},
		func(_ int) error {
			time.Sleep(3 * time.Second)
			return xerrors.ErrContinue
		},
	}
	ctx := context.TODO() // never timeout ctx
	_ = Run(ctx, 2, func(i int) error {
		return funcs[i](i)
	})
	time.Sleep(time.Second * 5)
}
