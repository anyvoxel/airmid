package xtime

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
)

func TestSystemTimerMills(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)

		tt := NewSystemTimer()
		tt.(*systemTimer).nowFn = func() time.Time {
			return time.Unix(10, 100100000)
		}
		g.Expect(tt.CurrentTimeMills()).To(gomega.Equal(uint64(10100)))
	})
}
