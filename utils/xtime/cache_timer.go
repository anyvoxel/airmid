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
