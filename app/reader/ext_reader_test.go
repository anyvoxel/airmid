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

package reader

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestExtReaderRead(t *testing.T) {
	g := NewWithT(t)

	vv := map[string]any{
		"1": 1,
	}
	r := &extReader{
		fn: func(data []byte) (map[string]any, error) {
			return vv, nil
		},
	}
	g.Expect(r.Read([]byte{})).To(Equal(vv))
}

func TestExtReaderMatch(t *testing.T) {
	type testCase struct {
		desp     string
		r        *extReader
		filename string
		err      string
	}
	testCases := []testCase{
		{
			desp: "match",
			r: &extReader{
				exts: []string{"1", "2"},
			},
			filename: "x2",
			err:      "",
		},
		{
			desp: "match",
			r: &extReader{
				exts: []string{"1", "2"},
			},
			filename: "x",
			err:      "'1,2' cannot support filename 'x'",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.r.Match(tc.filename)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestExtReaderName(t *testing.T) {
	g := NewWithT(t)
	r := &extReader{
		exts: []string{"1", "2"},
	}
	g.Expect(r.Name()).To(Equal("1,2"))
}
