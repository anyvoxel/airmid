package xapp

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/panjf2000/ants/v2"

	"github.com/anyvoxel/airmid/beans"
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
	p.PostProcessBeanDefinition("l", beans.MustNewBeanDefinition(reflect.TypeOf(t)))
	p.PostProcessBeanDefinition("l2", beans.MustNewBeanDefinition(reflect.TypeOf(t), beans.WithBeanScope(beans.ScopePrototype)))
	g.Expect(p.singletonNames).To(Equal(map[string]bool{
		"l":  true,
		"l2": false,
	}))
}
