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
