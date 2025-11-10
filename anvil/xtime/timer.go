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

// Package xtime provide the helper function to access time
package xtime

import (
	"sync/atomic"
	"time"
	"unsafe"
)

// UnixTimeMilliOffset is the millisecond offset of unix time.
var UnixTimeMilliOffset = uint64(time.Millisecond / time.Nanosecond)

// Timer is the interface to retrieve current time.
type Timer interface {
	// CurrentTimeMills will return current milli second
	CurrentTimeMills() uint64
}

// We cannot use atomic.Value, because it must use same concrete type
// var defaultTimer atomic.Value

var defaultTimer unsafe.Pointer

// Default will return the default timer.
func Default() Timer {
	return *(*Timer)(atomic.LoadPointer(&defaultTimer))
}

// SetDefault will update the default timer.
func SetDefault(t Timer) {
	atomic.StorePointer(&defaultTimer, unsafe.Pointer(&t))
}

// CurrentTimeMills will return the current milli second.
func CurrentTimeMills() uint64 {
	return Default().CurrentTimeMills()
}

// CurrentTime will return the current milli second in Time format.
func CurrentTime() time.Time {
	ms := CurrentTimeMills()
	return time.Unix(int64(ms/1000), int64(ms%1000)*int64(UnixTimeMilliOffset))
}

func init() {
	SetDefault(NewCacheTimer())
}
