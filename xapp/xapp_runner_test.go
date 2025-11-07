package xapp

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

type testAppRunner struct {
	runFn  func(ctx context.Context)
	stopFn func(ctx context.Context)
}

var (
	_ AppRunner = (*testAppRunner)(nil)
)

func (t *testAppRunner) Run(ctx context.Context) {
	if t.runFn != nil {
		t.runFn(ctx)
	}
}

func (t *testAppRunner) Stop(ctx context.Context) {
	if t.stopFn != nil {
		t.stopFn(ctx)
	}
}

func TestAppRunnerCompositorProcessor(t *testing.T) {
	t.Run("PostProcessBeforeInitialization", func(t *testing.T) {
		g := NewWithT(t)

		p := &AppRunnerCompoistorProcessor{
			appRunnerNames: map[AppRunner]string{},
		}
		p.PostProcessBeforeInitialization(1, "1")
		p.PostProcessBeforeInitialization(&testAppRunner{}, "2")

		g.Expect(p.appRunnerNames).To(HaveLen(1))
	})
}

func TestAppRunnerCompositorGetAppRunnerBeanName(t *testing.T) {
	t.Run("nil map", func(t *testing.T) {
		g := NewWithT(t)
		c := &AppRunnerCompositor{
			appRunnerNames: nil,
		}
		g.Expect(c.GetAppRunnerBeanName(context.Background(), nil)).To(Equal(""))
	})

	t.Run("not exists", func(t *testing.T) {
		g := NewWithT(t)
		r := &testAppRunner{}
		c := &AppRunnerCompositor{
			appRunnerNames: map[AppRunner]string{},
		}
		g.Expect(c.GetAppRunnerBeanName(context.Background(), r)).To(Equal(""))
	})

	t.Run("not exists", func(t *testing.T) {
		g := NewWithT(t)
		r := &testAppRunner{}
		c := &AppRunnerCompositor{
			appRunnerNames: map[AppRunner]string{
				r: "1",
			},
		}
		g.Expect(c.GetAppRunnerBeanName(context.Background(), r)).To(Equal("1"))
	})
}

func TestAppRunnerCompositorRun(t *testing.T) {
	g := NewWithT(t)
	count := 0
	r1 := &testAppRunner{
		runFn: func(ctx context.Context) {
			count++
		},
	}
	r2 := &testAppRunner{
		runFn: func(ctx context.Context) {
			count++
		},
	}
	c := &AppRunnerCompositor{
		runners: []AppRunner{r1, r2},
		appRunnerNames: map[AppRunner]string{
			r1: "r1",
			r2: "r2",
		},
	}
	c.Run(context.Background())
	g.Expect(count).To(Equal(2))
}

func TestAppRunnerCompositorStop(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		var wg sync.WaitGroup
		wg.Add(2)

		count := int32(0)
		r1 := &testAppRunner{
			stopFn: func(ctx context.Context) {
				defer wg.Done()
				atomic.AddInt32(&count, 1)
			},
		}
		r2 := &testAppRunner{
			stopFn: func(ctx context.Context) {
				defer wg.Done()
				atomic.AddInt32(&count, 1)
			},
		}
		c := &AppRunnerCompositor{
			runners: []AppRunner{r1, r2},
			appRunnerNames: map[AppRunner]string{
				r1: "r1",
				r2: "r2",
			},
		}
		err := c.Stop(context.Background())
		wg.Wait()
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(2)))
	})

	t.Run("context failed", func(t *testing.T) {
		g := NewWithT(t)

		var wg sync.WaitGroup
		wg.Add(2)

		count := int32(0)
		r1 := &testAppRunner{
			stopFn: func(ctx context.Context) {
				defer wg.Done()
				atomic.AddInt32(&count, 1)
				<-ctx.Done()
			},
		}
		r2 := &testAppRunner{
			stopFn: func(ctx context.Context) {
				defer wg.Done()
				atomic.AddInt32(&count, 1)
				<-ctx.Done()
			},
		}
		c := &AppRunnerCompositor{
			runners: []AppRunner{r1, r2},
			appRunnerNames: map[AppRunner]string{
				r1: "r1",
				r2: "r2",
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.Stop(ctx)
		wg.Wait()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err).To(Equal(context.DeadlineExceeded))
		g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(2)))
	})
}
