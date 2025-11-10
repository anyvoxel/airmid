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

package app

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
	_ Runner = (*testAppRunner)(nil)
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

		p := &RunnerCompoistorProcessor{
			appRunnerNames: map[Runner]string{},
		}
		p.PostProcessBeforeInitialization(1, "1")
		p.PostProcessBeforeInitialization(&testAppRunner{}, "2")

		g.Expect(p.appRunnerNames).To(HaveLen(1))
	})
}

func TestAppRunnerCompositorGetAppRunnerBeanName(t *testing.T) {
	t.Run("nil map", func(t *testing.T) {
		g := NewWithT(t)
		c := &RunnerCompositor{
			appRunnerNames: nil,
		}
		g.Expect(c.GetAppRunnerBeanName(context.Background(), nil)).To(Equal(""))
	})

	t.Run("not exists", func(t *testing.T) {
		g := NewWithT(t)
		r := &testAppRunner{}
		c := &RunnerCompositor{
			appRunnerNames: map[Runner]string{},
		}
		g.Expect(c.GetAppRunnerBeanName(context.Background(), r)).To(Equal(""))
	})

	t.Run("not exists", func(t *testing.T) {
		g := NewWithT(t)
		r := &testAppRunner{}
		c := &RunnerCompositor{
			appRunnerNames: map[Runner]string{
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
	c := &RunnerCompositor{
		runners: []Runner{r1, r2},
		appRunnerNames: map[Runner]string{
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
		c := &RunnerCompositor{
			runners: []Runner{r1, r2},
			appRunnerNames: map[Runner]string{
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
		c := &RunnerCompositor{
			runners: []Runner{r1, r2},
			appRunnerNames: map[Runner]string{
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
