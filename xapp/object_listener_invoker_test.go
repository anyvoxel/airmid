package xapp

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/panjf2000/ants/v2"

	"github.com/anyvoxel/airmid/xerrors"
)

func TestObjectListenerInvoker(t *testing.T) {
	g := NewWithT(t)

	count := int32(0)
	obj := &objectListenerInvoker{
		obj: &FuncApplicationListener{
			OnEventFunc: func(ctx context.Context, event ApplicationEvent) {
				atomic.AddInt32(&count, 1)
			},
		},
	}
	obj.Invoke(context.Background(), nil)
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(1)))

	var wg sync.WaitGroup
	obj = &objectListenerInvoker{
		objAsync: &FuncApplicationListener{
			OnEventAsyncFunc: func(ctx context.Context, event ApplicationEvent) {
				defer wg.Done()
				atomic.AddInt32(&count, 2)
			},
		},
		app: &airmidApplication{
			gpool: func() *ants.Pool {
				p, _ := ants.NewPool(10)
				return p
			}(),
		},
	}
	wg.Add(1)
	obj.Invoke(context.Background(), nil)
	wg.Wait()
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(3)))
}

func TestNewObjectListenerInvoker(t *testing.T) {
	g := NewWithT(t)
	count := int32(0)
	var wg sync.WaitGroup

	app := &airmidApplication{
		gpool: func() *ants.Pool {
			p, _ := ants.NewPool(10)
			return p
		}(),
	}

	vobj, err := NewObjectListenerInvoker(&FuncApplicationListener{
		OnEventFunc: func(ctx context.Context, event ApplicationEvent) {
			atomic.AddInt32(&count, 1)
		},
		OnEventAsyncFunc: func(ctx context.Context, event ApplicationEvent) {
			defer wg.Done()
			atomic.AddInt32(&count, 2)
		},
	}, app)
	g.Expect(err).ToNot(HaveOccurred())
	wg.Add(1)
	vobj.Invoke(context.Background(), nil)
	wg.Wait()
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(3)))

	vobj, err = NewObjectListenerInvoker(t, app)
	g.Expect(err).To(HaveOccurred())
	g.Expect(vobj).To(BeNil())
	g.Expect(xerrors.IsContinue(err)).To(BeTrue())
	g.Expect(err.Error()).To(MatchRegexp(`type '\*testing\.T' of object doesn't implement Listener: Continue`))
}
