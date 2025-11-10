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
	"testing"

	. "github.com/onsi/gomega"
)

func TestComplete(t *testing.T) {
	type testCase struct {
		desp   string
		o      *option
		expect *option
		err    string
	}
	testCases := []testCase{
		{
			desp: "count validate error",
			o:    &option{},
			err:  "count '0'",
		},
		{
			desp: "with default numcpu",
			o: &option{
				concurrent: 0,
				count:      runtime.NumCPU() + 1,
			},
			expect: &option{
				concurrent: runtime.NumCPU(),
				count:      runtime.NumCPU() + 1,
			},
		},
		{
			desp: "with min concurrent",
			o: &option{
				concurrent: 10,
				count:      9,
			},
			expect: &option{
				concurrent: 9,
				count:      9,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.o.Complete()
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.o).To(Equal(tc.expect))
		})
	}
}

func TestWithConcurrent(t *testing.T) {
	g := NewWithT(t)
	o := defaultOption()
	g.Expect(o.concurrent).To(Equal(runtime.NumCPU()))

	WithConcurrent(1)(o)
	g.Expect(o.concurrent).To(Equal(1))
}
