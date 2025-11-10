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
	"context"
	"runtime"
	"sync/atomic"
	"time"
)

type cacheTimer struct {
	nowInMs uint64
	ctx     context.Context

	// doneFunc for unit test
	doneFunc func()
}

func (t *cacheTimer) CurrentTimeMills() uint64 {
	return atomic.LoadUint64(&(t.nowInMs))
}

func (t *cacheTimer) startTicker() {
	go func() {
		for {
			select {
			case <-t.ctx.Done():
				if t.doneFunc != nil {
					t.doneFunc()
				}
				return
			default:
			}

			now := uint64(time.Now().UnixNano()) / UnixTimeMilliOffset
			atomic.StoreUint64(&(t.nowInMs), now)
			time.Sleep(time.Millisecond)
		}
	}()
}

type cacheTimeWrapper struct {
	*cacheTimer
}

// NewCacheTimer will return the implement of cache time.
func NewCacheTimer() Timer {
	return newCacheTimerWithFunc(nil)
}

// We make this function for test, to avoid assign doneFunc after the timer is created,
// because that will be an 'DATA RACE'.
func newCacheTimerWithFunc(fn func()) Timer {
	ctx, cancel := context.WithCancel(context.Background())

	// See more information:
	// 1. https://blog.kezhuw.name/2016/02/18/go-memory-leak-in-background-goroutine/
	// 2. https://zhuanlan.zhihu.com/p/76504936
	w := &cacheTimeWrapper{
		cacheTimer: &cacheTimer{
			nowInMs:  0,
			ctx:      ctx,
			doneFunc: fn,
		},
	}
	w.startTicker()
	runtime.SetFinalizer(w, func(_ *cacheTimeWrapper) {
		cancel()
	})

	return w
}
