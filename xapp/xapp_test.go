package xapp

import (
	"context"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/beans"
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

	app.RegisterBeanDefinition("l", beans.MustNewBeanDefinition(
		reflect.TypeOf((*testApplicationListener)(nil)),
	))
	object, err := app.GetBean("l")
	g.Expect(err).ToNot(HaveOccurred())
	l := object.(*testApplicationListener)

	app.PublishEvent(context.TODO(), ev1)
	app.PublishEvent(context.TODO(), ev2)
	g.Expect(l.ev1).To(Equal(1))
	g.Expect(l.ev2).To(Equal(2))
	g.Expect(l.ev).To(Equal(2))
}
