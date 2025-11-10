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

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

func TestFnListenerInvoker(t *testing.T) {
	g := NewWithT(t)

	count := int32(0)
	var wg sync.WaitGroup
	syncFunc := func(ctx context.Context, event ApplicationEvent) {
		atomic.AddInt32(&count, 1)
	}
	asyncFunc := func(ctx context.Context, event ApplicationEvent) {
		defer wg.Done()
		atomic.AddInt32(&count, 2)
	}
	obj := &fnListenerInvoker{
		async: false,
		fn:    reflect.ValueOf(syncFunc),
		arg1:  reflect.ValueOf(syncFunc).Type().In(1),
	}
	obj.Invoke(context.Background(), NewDefaultApplicationEvent(nil))
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(1)))
	obj.Invoke(context.Background(), nil)
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(2)))

	obj = &fnListenerInvoker{
		async: true,
		fn:    reflect.ValueOf(asyncFunc),
		arg1:  reflect.ValueOf(asyncFunc).Type().In(1),
		app: &airmidApplication{
			gpool: func() *ants.Pool {
				p, _ := ants.NewPool(10)
				return p
			}(),
		},
	}
	wg.Add(1)
	obj.Invoke(context.Background(), NewDefaultApplicationEvent(nil))
	wg.Wait()
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(4)))
}

type testFnEventListener struct {
	wg    *sync.WaitGroup
	count *int32
}

func (l *testFnEventListener) OnEvent1(ctx context.Context, event ApplicationEvent) {
	atomic.AddInt32(l.count, 1)
}

func (l *testFnEventListener) OnEvent1Async(ctx context.Context, event ApplicationEvent) {
	defer l.wg.Done()
	atomic.AddInt32(l.count, 2)
}

func (l *testFnEventListener) OnEvent2(ctx context.Context, event *testFnEventListenerEvent) {
	atomic.AddInt32(l.count, 3)
}

func (l *testFnEventListener) On1(ctx context.Context) {

}

func (l *testFnEventListener) On2(ctx context.Context, event *testFnEventListenerEvent) error {
	return nil
}

func (l *testFnEventListener) On1Type(int, ApplicationEvent) {

}

func (l *testFnEventListener) On2Type(ctx context.Context, i int) {

}

type testFnEventListenerEvent struct {
	ApplicationEvent
}

func TestNewFnListenerInvoker(t *testing.T) {
	g := NewWithT(t)
	count := int32(0)
	var wg sync.WaitGroup

	obj := &testFnEventListener{
		count: &count,
		wg:    &wg,
	}
	app := &airmidApplication{
		gpool: func() *ants.Pool {
			p, _ := ants.NewPool(10)
			return p
		}(),
	}
	vobj, err := NewFnListenerInvoker(reflect.ValueOf(obj).MethodByName("OnEvent1"), "OnEvent1", app)
	g.Expect(err).ToNot(HaveOccurred())
	vobj.Invoke(context.Background(), nil)
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(1)))
	vobj, err = NewFnListenerInvoker(reflect.ValueOf(obj).MethodByName("OnEvent1Async"), "OnEvent1Async", app)
	g.Expect(err).ToNot(HaveOccurred())
	wg.Add(1)
	vobj.Invoke(context.Background(), nil)
	wg.Wait()
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(3)))
	vobj, err = NewFnListenerInvoker(reflect.ValueOf(obj).MethodByName("OnEvent2"), "OnEvent2", app)
	g.Expect(err).ToNot(HaveOccurred())
	vobj.Invoke(context.Background(), nil)
	g.Expect(atomic.LoadInt32(&count)).To(Equal(int32(6)))

	vobj, err = NewFnListenerInvoker(reflect.ValueOf(t), "o1", app)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp(`type '\*testing\.T' of object 'o1' doesn't match kind func: Continue`))
	g.Expect(xerrors.IsContinue(err)).To(BeTrue())
	g.Expect(vobj).To(BeNil())

	vobj, err = NewFnListenerInvoker(reflect.ValueOf(obj).MethodByName("On1"), "On1", app)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp(`func 'On1' expect in '2' and out '0', got in '1' and out '0': Continue`))
	g.Expect(xerrors.IsContinue(err)).To(BeTrue())
	g.Expect(vobj).To(BeNil())

	vobj, err = NewFnListenerInvoker(reflect.ValueOf(obj).MethodByName("On2"), "On2", app)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp(`func 'On2' expect in '2' and out '0', got in '2' and out '1': Continue`))
	g.Expect(xerrors.IsContinue(err)).To(BeTrue())
	g.Expect(vobj).To(BeNil())

	vobj, err = NewFnListenerInvoker(reflect.ValueOf(obj).MethodByName("On1Type"), "On1Type", app)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp(`func 'On1Type' arg0 type 'int', expect 'context.Context': Continue`))
	g.Expect(xerrors.IsContinue(err)).To(BeTrue())
	g.Expect(vobj).To(BeNil())

	vobj, err = NewFnListenerInvoker(reflect.ValueOf(obj).MethodByName("On2Type"), "On2Type", app)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp(`func 'On2Type' arg1 type 'int', expect implement ApplicationEvent: Continue`))
	g.Expect(xerrors.IsContinue(err)).To(BeTrue())
	g.Expect(vobj).To(BeNil())

	vobj, err = NewFnListenerInvoker(reflect.ValueOf(&FuncApplicationListener{}).MethodByName("OnApplicationEvent"), "OnApplicationEvent", app)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp(`func 'OnApplicationEvent' is object listener interface: Continue`))
	g.Expect(xerrors.IsContinue(err)).To(BeTrue())
	g.Expect(vobj).To(BeNil())

	vobj, err = NewFnListenerInvoker(reflect.ValueOf(&FuncApplicationListener{}).MethodByName("OnApplicationEventAsync"), "OnApplicationEventAsync", app)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp(`func 'OnApplicationEventAsync' is object listener interface: Continue`))
	g.Expect(xerrors.IsContinue(err)).To(BeTrue())
	g.Expect(vobj).To(BeNil())
}
