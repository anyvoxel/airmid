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
	"errors"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/anvil"
	"github.com/anyvoxel/airmid/anvil/xerrors"
)

func TestRun(t *testing.T) {
	type testCase struct {
		desp               string
		count              int
		workFunc           func(int) error
		opts               []anvil.Option[parallelOption]
		expectedTypedError any
	}
	testCases := []testCase{
		{
			desp:               "option error",
			count:              0,
			workFunc:           nil,
			expectedTypedError: errors.New("count '0'"),
		},
		{
			desp:  "all success",
			count: 5,
			workFunc: func(_ int) error {
				return nil
			},
			opts: []anvil.Option[parallelOption]{
				WithConcurrent(2),
			},
			expectedTypedError: nil,
		},
		{
			desp:  "some failed",
			count: 10,
			workFunc: func(i int) error {
				if i == 7 {
					return xerrors.NewTyped(xerrors.Continue{})
				}
				return nil
			},
			opts: []anvil.Option[parallelOption]{
				WithConcurrent(4),
			},
			expectedTypedError: &xerrors.TypedError[xerrors.Continue]{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := Run(context.Background(), tc.count, tc.workFunc, tc.opts...)

			if tc.expectedTypedError == nil {
				g.Expect(err).ToNot(HaveOccurred())
				return
			}

			g.Expect(err).To(HaveOccurred())
			switch expected := tc.expectedTypedError.(type) {
			case *xerrors.TypedError[xerrors.Continue]:
				var target *xerrors.TypedError[xerrors.Continue]
				g.Expect(xerrors.As(err, &target)).To(BeTrue())
			case error: // Catches errors.New() and other standard errors
				g.Expect(err.Error()).To(MatchRegexp(expected.Error()))
			default:
				t.Fatalf("Unexpected expected error type: %T", tc.expectedTypedError)
			}
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
			return xerrors.NewTyped(xerrors.Continue{})
		},
		func(_ int) error {
			return xerrors.NewTyped(xerrors.Continue{})
		},
	}
	ctx := context.TODO() // never timeout ctx
	_ = Run(ctx, 2, func(i int) error {
		return funcs[i](i)
	})
	time.Sleep(time.Second * 5)
}
