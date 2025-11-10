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

package anvil

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

func TestMust(t *testing.T) {
	type testCase struct {
		desp   string
		err    error
		expect error
	}
	testCases := []testCase{
		{
			desp:   "panic with err",
			err:    xerrors.ErrNotFound,
			expect: xerrors.ErrNotFound,
		},
		{
			desp:   "not panic",
			err:    nil,
			expect: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			if tc.expect != nil {
				g.Ω(func() {
					Must(tc.err)
				}).To(PanicWith(tc.expect))
			} else {
				g.Ω(func() {
					Must(tc.err)
				}).NotTo(Panic())
			}
		})
	}
}
