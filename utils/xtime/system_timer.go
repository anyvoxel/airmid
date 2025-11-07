package xtime

import (
	"time"
)

type systemTimer struct {
	nowFn func() time.Time
}

func (t *systemTimer) CurrentTimeMills() uint64 {
	return uint64(t.nowFn().UnixNano()) / UnixTimeMilliOffset
}

// NewSystemTimer will return the implement of system time.
func NewSystemTimer() Timer {
	return &systemTimer{
		nowFn: time.Now,
	}
}
