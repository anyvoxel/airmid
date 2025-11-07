package xapp

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestApplicationEventType(t *testing.T) {
	g := NewWithT(t)

	fn := func(event ApplicationEvent) {}
	g.Expect(reflect.TypeOf(fn).In(0)).To(Equal(applicationEventType))
}

func TestOnApplicationEventFuncType(t *testing.T) {
	g := NewWithT(t)

	obj := &FuncApplicationListener{}
	m := reflect.ValueOf(obj).MethodByName("OnApplicationEvent").Type()

	g.Expect(m).To(Equal(onApplicationEventFuncType))
}

func TestOnApplicationEventAsyncFuncType(t *testing.T) {
	g := NewWithT(t)

	obj := &FuncApplicationListener{}
	m := reflect.ValueOf(obj).MethodByName("OnApplicationEventAsync").Type()

	g.Expect(m).To(Equal(onApplicationEventAsyncFuncType))
}
