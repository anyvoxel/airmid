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
