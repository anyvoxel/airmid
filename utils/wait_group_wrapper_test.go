package utils

import (
	"sync/atomic"
	"testing"

	. "github.com/onsi/gomega"
)

func TestWaitGroup(t *testing.T) {
	g := NewWithT(t)

	count := int32(0)
	wg := &WaitGroupWrapper{}
	wg.Wrap(func() {
		atomic.AddInt32(&count, 1)
	})
	wg.Wrap(func() {
		atomic.AddInt32(&count, 2)
	})
	wg.Wait()

	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(3)))
}
