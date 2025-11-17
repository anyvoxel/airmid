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
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/ioc"
)

type testApplicationEvent1 struct {
	ApplicationEvent
	count int
}

type testApplicationEvent2 struct {
	ApplicationEvent
	count int
}

type testApplicationListener struct {
	ev1 int
	ev2 int
	ev  int
}

func (l *testApplicationListener) OnEvent1(_ context.Context, ev testApplicationEvent1) {
	l.ev1 += ev.count
}

func (l *testApplicationListener) OnEvent2(_ context.Context, ev testApplicationEvent2) {
	l.ev2 += ev.count
}

func (l *testApplicationListener) OnEvent(_ context.Context, _ ApplicationEvent) {
	l.ev++
}

func TestPublishEvent(t *testing.T) {
	g := NewWithT(t)
	ev1 := testApplicationEvent1{
		ApplicationEvent: NewDefaultApplicationEvent(nil),
		count:            1,
	}
	ev2 := testApplicationEvent2{
		ApplicationEvent: NewDefaultApplicationEvent(nil),
		count:            2,
	}
	app := NewApplication().(*airmidApplication)
	app.AddBeanPostProcessor(&applicationListenerDetector{
		app: app,
		singletonNames: map[string]bool{
			"l": true,
		},
	})

	app.RegisterBeanDefinition("l", ioc.MustNewBeanDefinition(
		reflect.TypeOf((*testApplicationListener)(nil)),
	))
	object, err := app.GetBean(context.Background(), "l")
	g.Expect(err).ToNot(HaveOccurred())
	l := object.(*testApplicationListener)

	app.PublishEvent(context.TODO(), ev1)
	app.PublishEvent(context.TODO(), ev2)
	g.Expect(l.ev1).To(Equal(1))
	g.Expect(l.ev2).To(Equal(2))
	g.Expect(l.ev).To(Equal(2))
}
