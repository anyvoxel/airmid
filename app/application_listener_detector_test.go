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
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/panjf2000/ants/v2"

	"github.com/anyvoxel/airmid/ioc"
)

type testPostProcessListener struct {
	wg    *sync.WaitGroup
	count *int32
}

func (l *testPostProcessListener) OnEvent1(ctx context.Context, event ApplicationEvent) {
	atomic.AddInt32(l.count, 1)
}

func (l *testPostProcessListener) OnEvent1Async(ctx context.Context, event ApplicationEvent) {
	defer l.wg.Done()

	atomic.AddInt32(l.count, 5)
}

func (l *testPostProcessListener) OnEvent2(ctx context.Context, event *testFnEventListenerEvent) {
	atomic.AddInt32(l.count, 2)
}

func (l *testPostProcessListener) OnApplicationEvent(ctx context.Context, event ApplicationEvent) {
	atomic.AddInt32(l.count, 3)
}

func (l *testPostProcessListener) OnApplicationEventAsync(ctx context.Context, event ApplicationEvent) {
	defer l.wg.Done()

	atomic.AddInt32(l.count, 4)
}

func TestPostProcessAfterInitialization(t *testing.T) {
	count := int32(0)
	var wg sync.WaitGroup

	l := &testPostProcessListener{
		count: &count,
		wg:    &wg,
	}
	app := &airmidApplication{
		listenerInvoker: make([]ListenerInvoker, 0),
		gpool: func() *ants.Pool {
			p, _ := ants.NewPool(10)
			return p
		}(),
	}
	p := &applicationListenerDetector{
		app: app,
		singletonNames: map[string]bool{
			"l":  true,
			"l2": false,
		},
	}
	g := NewWithT(t)
	p.PostProcessAfterInitialization(l, "l")
	g.Expect(app.listenerInvoker).To(HaveLen(4))

	wg.Add(2)
	app.PublishEvent(context.Background(), NewDefaultApplicationEvent(nil))
	wg.Wait()
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(13)))
	atomic.StoreInt32(&count, 0)

	wg.Add(2)
	app.PublishEvent(context.Background(), &testFnEventListenerEvent{})
	wg.Wait()
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(15)))

	p.PostProcessAfterInitialization(l, "l2")
	g.Expect(app.listenerInvoker).To(HaveLen(4))
}

func TestPostProcessBeanDefinition(t *testing.T) {
	g := NewWithT(t)

	app := &airmidApplication{
		listenerInvoker: make([]ListenerInvoker, 0),
	}
	p := &applicationListenerDetector{
		app:            app,
		singletonNames: map[string]bool{},
	}
	p.PostProcessBeanDefinition("l", ioc.MustNewBeanDefinition(reflect.TypeOf(t)))
	p.PostProcessBeanDefinition("l2", ioc.MustNewBeanDefinition(reflect.TypeOf(t), ioc.WithBeanScope(ioc.ScopePrototype)))
	g.Expect(p.singletonNames).To(Equal(map[string]bool{
		"l":  true,
		"l2": false,
	}))
}
