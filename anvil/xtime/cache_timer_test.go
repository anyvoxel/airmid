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
	"runtime"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/onsi/gomega"
)

func TestCacheTimerMills(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		var done int32
		var doneCh = make(chan struct{})
		var cnt int32
		var p unsafe.Pointer

		// NOTE: we use pointer to avoid 'DATA RACE' in test, only the atomic Add/Load cann't fix it
		atomic.StorePointer(&p, unsafe.Pointer(&cnt))

		guard := gomonkey.ApplyFunc(time.Now, func() time.Time {
			pp := (*int32)(atomic.LoadPointer(&p))
			atomic.AddInt32(pp, 1)
			return time.Unix(10, 100100000)
		})
		defer guard.Reset()

		tt := newCacheTimerWithFunc(func() {
			atomic.AddInt32(&done, 1)
			close(doneCh)
		})
		time.Sleep(time.Second)

		pp := (*int32)(atomic.LoadPointer(&p))
		cntt := atomic.LoadInt32(pp)
		g.Expect(func() bool {
			return cntt >= 900
		}()).To(gomega.BeTrue())

		g.Expect(tt.CurrentTimeMills()).To(gomega.Equal(uint64(10100)))

		tt = nil

		// NOTE: we must force gc twice, to ensure the Finalizer && object gc have done
		runtime.GC() // force gc
		runtime.GC() // force gc
		<-doneCh
		g.Expect(atomic.LoadInt32(&done)).To(gomega.Equal(int32(1)))
	})
}
