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

package xtime

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
)

type testTimer struct {
	ms uint64
}

func (t *testTimer) CurrentTimeMills() uint64 {
	return t.ms
}

func TestDefault(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		tt := &testTimer{
			ms: 100,
		}
		SetDefault(tt)

		g.Expect(Default()).To(gomega.Equal(tt))
	})
}

func TestCurrentTimeMills(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		tt := &testTimer{
			ms: 100,
		}
		SetDefault(tt)

		g.Expect(CurrentTimeMills()).To(gomega.Equal(tt.ms))
	})
}

func TestCurrentTime(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		tt := &testTimer{
			ms: 10100,
		}
		SetDefault(tt)

		g.Expect(CurrentTime()).To(gomega.Equal(time.Unix(10, 100000000)))
	})
}
