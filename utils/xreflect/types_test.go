package xreflect

import (
	"context"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestContextType(t *testing.T) {
	g := NewWithT(t)

	fn1 := func(ctx context.Context) {}
	g.Expect(reflect.TypeOf(fn1).In(0)).To(Equal(ContextType))

	fn2 := func(*testing.T) {}
	g.Expect(reflect.TypeOf(fn2).In(0)).ToNot(Equal(ContextType))
}

func TestErrorType(t *testing.T) {
	g := NewWithT(t)

	fn1 := func(error) {}
	g.Expect(reflect.TypeOf(fn1).In(0)).To(Equal(ErrorType))

	fn2 := func(*testing.T) {}
	g.Expect(reflect.TypeOf(fn2).In(0)).ToNot(Equal(ErrorType))
}
