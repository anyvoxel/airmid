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

package ioc

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/anvil/pointer"
	"github.com/anyvoxel/airmid/anvil/xerrors"
)

type testBeanOnlyPropertyField struct {
	F0 int   `airmid:"value:${props.f0}"`
	f1 *int  `airmid:"value:${props.f1}"`
	F2 []int `airmid:"value:${props.f2:=}"`
	f3 []int `airmid:"value:${props.f3:=1,3,2}"`
}

// Implement the fmt.Stringer.
func (v *testBeanOnlyPropertyField) String() string {
	return ""
}

type testBeanOnlyBeanField struct {
	F0 *testBeanOnlyPropertyField `airmid:"autowire:?"`     // autowire by type explicit
	f1 []fmt.Stringer             `airmid:"autowire:?"`     // autowire by implement
	F2 *testBeanOnlyPropertyField `airmid:"autowire:bean1"` // autowire by name
}

type testBeanWithOptionalBeanField struct {
	F0 fmt.Stringer `airmid:"autowire:?,optional"`     // autowire by type, and be optional
	F1 fmt.Stringer `airmid:"autowire:bean1,optional"` // autowire by name, and be optional
}

type testBeanMixField struct {
	F0 *testBeanOnlyPropertyField `airmid:"autowire:?"`
	f1 int                        `airmid:"value:${props.f1}"`
}

type testBeanNotImplement struct {
	notImplementInterface testNotImplementInterface `airmid:"autowire:?"` //nolint
}

// Implement the fmt.Stringer.
func (v *testBeanNotImplement) String() string {
	return ""
}

func TestBeanFactoryGetBean(t *testing.T) {
	type testCase struct {
		desp     string
		initFunc func(g *WithT, tc *testCase, bf BeanFactory)
		beanName string
		err      string
		expect   any
	}
	testCases := []testCase{
		{
			desp: "normal get bean only property field",
			initFunc: func(g *WithT, tc *testCase, bf BeanFactory) {
				err := bf.RegisterBeanDefinition("bean2", MustNewBeanDefinition(reflect.TypeOf((*testBeanOnlyPropertyField)(nil)), WithBeanName("bean2")))
				g.Expect(err).ToNot(HaveOccurred())

				err = bf.Set(context.Background(), "props.f0", "100")
				g.Expect(err).ToNot(HaveOccurred())
				err = bf.Set(context.Background(), "props.f1", "200")
				g.Expect(err).ToNot(HaveOccurred())
				err = bf.Set(context.Background(), "props.f2", []string{"1", "2"})
				g.Expect(err).ToNot(HaveOccurred())
			},
			beanName: "bean2",
			err:      "",
			expect: &testBeanOnlyPropertyField{
				F0: 100,
				f1: pointer.IntPtr(int(200)),
				F2: []int{1, 2},
				f3: []int{1, 3, 2},
			},
		},
		{
			desp: "normal get bean only bean field",
			initFunc: func(g *WithT, tc *testCase, bf BeanFactory) {
				err := bf.RegisterBeanDefinition("bean1", MustNewBeanDefinition(reflect.TypeOf((*testBeanOnlyPropertyField)(nil)), WithBeanName("bean1")))
				g.Expect(err).ToNot(HaveOccurred())
				err = bf.RegisterBeanDefinition("bean2", MustNewBeanDefinition(reflect.TypeOf((*testBeanOnlyBeanField)(nil))))
				g.Expect(err).ToNot(HaveOccurred())

				err = bf.Set(context.Background(), "props.f0", "100")
				g.Expect(err).ToNot(HaveOccurred())
				err = bf.Set(context.Background(), "props.f1", "200")
				g.Expect(err).ToNot(HaveOccurred())
				err = bf.Set(context.Background(), "props.f2", []string{"1", "2"})
				g.Expect(err).ToNot(HaveOccurred())
			},
			beanName: "bean2",
			err:      "",
			expect: &testBeanOnlyBeanField{
				F0: &testBeanOnlyPropertyField{
					F0: 100,
					f1: pointer.IntPtr(int(200)),
					F2: []int{1, 2},
					f3: []int{1, 3, 2},
				},
				f1: []fmt.Stringer{
					&testBeanOnlyPropertyField{
						F0: 100,
						f1: pointer.IntPtr(int(200)),
						F2: []int{1, 2},
						f3: []int{1, 3, 2},
					},
				},
				F2: &testBeanOnlyPropertyField{
					F0: 100,
					f1: pointer.IntPtr(int(200)),
					F2: []int{1, 2},
					f3: []int{1, 3, 2},
				},
			},
		},
		{
			desp: "normal get bean only bean field with optional not found",
			initFunc: func(g *WithT, tc *testCase, bf BeanFactory) {
				var err error //nolint
				err = bf.RegisterBeanDefinition("bean2", MustNewBeanDefinition(reflect.TypeOf((*testBeanWithOptionalBeanField)(nil))))
				g.Expect(err).ToNot(HaveOccurred())
			},
			beanName: "bean2",
			err:      "",
			expect: &testBeanWithOptionalBeanField{
				F0: nil,
			},
		},
		{
			desp: "normal get bean with mix field",
			initFunc: func(g *WithT, tc *testCase, bf BeanFactory) {
				err := bf.RegisterBeanDefinition("bean1", MustNewBeanDefinition(reflect.TypeOf((*testBeanOnlyPropertyField)(nil)), WithBeanName("bean1")))
				g.Expect(err).ToNot(HaveOccurred())
				err = bf.RegisterBeanDefinition("bean2", MustNewBeanDefinition(reflect.TypeOf((*testBeanMixField)(nil))))
				g.Expect(err).ToNot(HaveOccurred())

				err = bf.Set(context.Background(), "props.f0", "100")
				g.Expect(err).ToNot(HaveOccurred())
				err = bf.Set(context.Background(), "props.f1", "200")
				g.Expect(err).ToNot(HaveOccurred())
				err = bf.Set(context.Background(), "props.f2", []string{"1", "2"})
				g.Expect(err).ToNot(HaveOccurred())
			},
			beanName: "bean2",
			err:      "",
			expect: &testBeanMixField{
				F0: &testBeanOnlyPropertyField{
					F0: 100,
					f1: pointer.IntPtr(int(200)),
					F2: []int{1, 2},
					f3: []int{1, 3, 2},
				},
				f1: 200,
			},
		},
		{
			desp: "lazymode bean with error",
			initFunc: func(g *WithT, tc *testCase, bf BeanFactory) {
				err := bf.RegisterBeanDefinition("beanNotImplement", MustNewBeanDefinition(reflect.TypeOf((*testBeanNotImplement)(nil)), WithLazyMode()))
				g.Expect(err).ToNot(HaveOccurred())
			},
			beanName: "beanNotImplement",
			err:      "No candidate found for field 'notImplementInterface' with type ioc.testNotImplementInterface",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			f := NewBeanFactory()
			tc.initFunc(g, &tc, f)

			actual, err := f.GetBean(context.Background(), tc.beanName)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

type testSingletonBean struct {
	v     int
	count int
}

func (b *testSingletonBean) AfterPropertiesSet() error {
	b.count++
	if b.count > 1 {
		return xerrors.Errorf("Too many instance with singleton scope object")
	}
	return nil
}

func TestGetBeanSingletonScope(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory()
	bf.RegisterBeanDefinition("bean1", MustNewBeanDefinition(
		reflect.TypeOf((*testSingletonBean)(nil)),
	))

	obj1, err := bf.GetBean(context.Background(), "bean1")
	g.Expect(err).ToNot(HaveOccurred())
	bean1 := obj1.(*testSingletonBean)
	g.Expect(bean1.v).To(Equal(0))
	bean1.v++
	g.Expect(bean1.v).To(Equal(1))

	obj2, err := bf.GetBean(context.Background(), "bean1")
	g.Expect(err).ToNot(HaveOccurred())
	bean2 := obj2.(*testSingletonBean)
	g.Expect(bean2.v).To(Equal(1))

	g.Expect(reflect.ValueOf(bean1).Pointer()).To(Equal(reflect.ValueOf(bean2).Pointer()))
}

var (
	prototypeBeanCount = 0
)

type testPrototypeBean struct {
	v int
}

func (b *testPrototypeBean) AfterPropertiesSet(context.Context) error {
	prototypeBeanCount++
	return nil
}

func TestGetBeanPrototypeScope(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory()
	bf.RegisterBeanDefinition("bean1", MustNewBeanDefinition(
		reflect.TypeOf((*testPrototypeBean)(nil)),
		WithBeanScope(ScopePrototype),
	))

	obj1, err := bf.GetBean(context.Background(), "bean1")
	g.Expect(err).ToNot(HaveOccurred())
	bean1 := obj1.(*testPrototypeBean)
	g.Expect(bean1.v).To(Equal(0))
	g.Expect(prototypeBeanCount).To(Equal(1))
	bean1.v++
	g.Expect(bean1.v).To(Equal(1))

	obj2, err := bf.GetBean(context.Background(), "bean1")
	g.Expect(err).ToNot(HaveOccurred())
	bean2 := obj2.(*testPrototypeBean)
	g.Expect(bean2.v).To(Equal(0))
	g.Expect(prototypeBeanCount).To(Equal(2))

	g.Expect(reflect.ValueOf(bean1).Pointer()).ToNot(Equal(reflect.ValueOf(bean2).Pointer()))
}

func TestResolveBeans(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory()
	bf.RegisterBeanDefinition("bean3", MustNewBeanDefinition(
		reflect.TypeOf((*testBeanOnlyPropertyField)(nil)),
	))
	bf.RegisterBeanDefinition("bean1", MustNewBeanDefinition(
		reflect.TypeOf((*testBeanOnlyPropertyField)(nil)),
	))
	bf.RegisterBeanDefinition("bean2", MustNewBeanDefinition(
		reflect.TypeOf((*testBeanMixField)(nil)),
	))

	actual, err := bf.ResolveBeanNames(context.Background(), reflect.TypeOf((*fmt.Stringer)(nil)).Elem())
	g.Expect(err).ToNot(HaveOccurred())

	sort.Strings(actual)
	g.Expect(actual).To(Equal([]string{"bean1", "bean3"}))
}

func TestResolveBeansWithLazymode(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory()
	bf.RegisterBeanDefinition("testBeanNotImplement", MustNewBeanDefinition(
		reflect.TypeOf((*testBeanNotImplement)(nil)), WithLazyMode(),
	))
	bf.RegisterBeanDefinition("bean1", MustNewBeanDefinition(
		reflect.TypeOf((*testBeanOnlyPropertyField)(nil)),
	))

	actual, err := bf.ResolveBeanNames(context.Background(), reflect.TypeOf((*fmt.Stringer)(nil)).Elem())
	g.Expect(err).ToNot(HaveOccurred())

	sort.Strings(actual)
	g.Expect(actual).To(Equal([]string{"bean1", "testBeanNotImplement"}))
}

type testScope struct {
	name string
}

func (s *testScope) Get(name string, objectFactory ObjectFactory) (any, error) {
	return nil, nil
}

func (s *testScope) Remove(name string) error {
	return nil
}

func TestRegisterScope(t *testing.T) {
	type testCase struct {
		desp  string
		bf    *beanFactoryImpl
		init  func(tc *testCase)
		name  string
		scope Scope
		err   string

		scopes map[string]Scope
	}

	testCases := []testCase{
		{
			desp: "normal register",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			init: func(tc *testCase) {
				tc.bf.scopes = map[string]Scope{
					"s1": &testScope{
						name: "s1",
					},
				}
			},
			name: "s2",
			scope: &testScope{
				name: "s2",
			},
			err: "",
			scopes: map[string]Scope{
				"s1": &testScope{
					name: "s1",
				},
				"s2": &testScope{
					name: "s2",
				},
			},
		},
		{
			desp: "register with replace",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			init: func(tc *testCase) {
				tc.bf.scopes = map[string]Scope{
					"s1": &testScope{
						name: "s1",
					},
				}
			},
			name: "s1",
			scope: &testScope{
				name: "s2",
			},
			err: "",
			scopes: map[string]Scope{
				"s1": &testScope{
					name: "s2",
				},
			},
		},
		{
			desp: "singleton error",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			init: func(tc *testCase) {
				tc.bf.scopes = map[string]Scope{
					"s1": &testScope{
						name: "s1",
					},
				}
			},
			name: "singleton",
			scope: &testScope{
				name: "s2",
			},
			err: "Invalid scope name 'singleton'",
		},
		{
			desp: "prototype error",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			init: func(tc *testCase) {
				tc.bf.scopes = map[string]Scope{
					"s1": &testScope{
						name: "s1",
					},
				}
			},
			name: "prototype",
			scope: &testScope{
				name: "s2",
			},
			err: "Invalid scope name 'prototype'",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			tc.init(&tc)

			err := tc.bf.RegisterScope(tc.name, tc.scope)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.bf.scopes).To(Equal(tc.scopes))
		})
	}
}

func TestRegisterSingleton(t *testing.T) {
	type testCase struct {
		desp string
		bf   *beanFactoryImpl
		init func(tc *testCase)
		name string
		bean any
		err  string
	}
	testCases := []testCase{
		{
			desp: "normal register",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			init: func(tc *testCase) {
				tc.bf.RegisterSingleton("s2", "bean2")
			},
			name: "s1",
			bean: "bean1",
			err:  "",
		},
		{
			desp: "register duplicate",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			init: func(tc *testCase) {
				tc.bf.RegisterSingleton("s1", "bean1")
			},
			name: "s1",
			bean: "bean1",
			err:  "Invalid object 'bean1' under bean name 's1'",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			tc.init(&tc)

			err := tc.bf.RegisterSingleton(tc.name, tc.bean)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.bf.singletonObjects[tc.name].Interface()).To(Equal(tc.bean))
		})
	}
}

func TestGetBeanFromRegisterSingleton(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory()
	bean1 := pointer.StringPtr("1")
	bf.RegisterSingleton("bean1", bean1)

	obj1, err := bf.GetBean(context.Background(), "bean1")
	g.Expect(err).ToNot(HaveOccurred())
	beanGet1 := obj1.(*string)
	g.Expect(beanGet1).To(Equal(bean1))
	g.Expect(reflect.ValueOf(beanGet1).Pointer()).To(Equal(reflect.ValueOf(bean1).Pointer()))

	*bean1 = "2"
	g.Expect(beanGet1).To(Equal(pointer.StringPtr("2")))
}

type testAware struct {
	name        string
	beanFactory BeanFactory
}

func (a *testAware) SetBeanName(beanName string) {
	a.name = beanName
}

func (a *testAware) SetBeanFactory(beanFactory BeanFactory) {
	a.beanFactory = beanFactory
}

func TestWithAwareBeanPostProcessor(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory()

	bf.RegisterBeanDefinition("bean1", MustNewBeanDefinition(
		reflect.TypeOf((*testAware)(nil)),
	))
	obj, err := bf.GetBean(context.Background(), "bean1")
	g.Expect(err).ToNot(HaveOccurred())
	a := obj.(*testAware)
	g.Expect(a).To(Equal(&testAware{
		name:        "bean1",
		beanFactory: bf,
	}))
}

func TestPreInstantiateSingletons(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory().(*beanFactoryImpl)

	type testPrototype1 struct {
		val int `airmid:"value:${val:=1}"`
	}
	type testPrototype2 struct {
		notImplementInterface testNotImplementInterface `airmid:"autowire:?"` //nolint
	}
	type testSingleton1 struct {
		p *testPrototype1 `airmid:"autowire:?"`
	}
	type testSingleton2 struct {
		notImplementInterface testNotImplementInterface `airmid:"autowire:?"` //nolint
	}

	err := bf.RegisterBeanDefinition("s1", MustNewBeanDefinition(
		reflect.TypeOf((*testSingleton1)(nil)),
	))
	g.Expect(err).ToNot(HaveOccurred())
	err = bf.RegisterBeanDefinition("s2", MustNewBeanDefinition(
		reflect.TypeOf((*testSingleton2)(nil)),
		WithLazyMode(),
	))
	g.Expect(err).ToNot(HaveOccurred())
	err = bf.RegisterBeanDefinition("p1", MustNewBeanDefinition(
		reflect.TypeOf((*testPrototype1)(nil)),
		WithBeanScope(ScopePrototype),
	))
	g.Expect(err).ToNot(HaveOccurred())
	err = bf.RegisterBeanDefinition("p2", MustNewBeanDefinition(
		reflect.TypeOf((*testPrototype2)(nil)),
		WithBeanScope(ScopePrototype),
		WithLazyMode(),
	))
	g.Expect(err).ToNot(HaveOccurred())

	err = bf.PreInstantiateSingletons(context.Background())
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(bf.singletonObjects).To(HaveLen(1))
	g.Expect(bf.singletonObjects["s1"].Interface()).To(Equal(&testSingleton1{
		p: &testPrototype1{
			val: 1,
		},
	}))
}

type singleton struct{}

func TestPreInstantiateSingletonsWithMoreBeans(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)
		br := NewBeanFactory().(*beanFactoryImpl)

		// We make this size bigger than runtime.NumCPU
		beansCount := 200
		for i := 0; i < beansCount; i++ {
			err := br.RegisterBeanDefinition(fmt.Sprintf("test%v", i), MustNewBeanDefinition(
				reflect.TypeOf((*singleton)(nil)),
			))
			g.Expect(err).ShouldNot(HaveOccurred())
		}
		err := br.PreInstantiateSingletons(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred())
		g.Expect(br.singletonObjects).To(HaveLen(beansCount))
		for i := 0; i < beansCount; i++ {
			name := fmt.Sprintf("test%v", i)
			obj := br.singletonObjects[name].Interface()
			g.Expect(obj).To(Equal(&singleton{}))
		}
	})

	t.Run("empty singleton test", func(t *testing.T) {
		g := NewWithT(t)
		br := NewBeanFactory().(*beanFactoryImpl)

		// We make this size bigger than runtime.NumCPU
		err := br.PreInstantiateSingletons(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred())
		g.Expect(br.singletonObjects).To(BeEmpty())
	})

	t.Run("some failed", func(t *testing.T) {
		g := NewWithT(t)
		br := NewBeanFactory().(*beanFactoryImpl)

		type failedSingleton struct {
			F1 int `airmid:"value:${props.1}"`
		}

		// We make this size bigger than runtime.NumCPU
		beansCount := 200
		for i := 0; i < beansCount; i++ {
			var err error
			if i != 30 {
				err = br.RegisterBeanDefinition(fmt.Sprintf("test%v", i), MustNewBeanDefinition(
					reflect.TypeOf((*singleton)(nil)),
				))
			} else {
				err = br.RegisterBeanDefinition(fmt.Sprintf("test%v", i), MustNewBeanDefinition(
					reflect.TypeOf((*failedSingleton)(nil)),
				))
			}

			g.Expect(err).ShouldNot(HaveOccurred())
		}
		err := br.PreInstantiateSingletons(context.Background())
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`property with key='props.1' not found`))
	})
}

type smartSingleton struct {
	a string
}

func (s *smartSingleton) AfterSingletonsInstantiated(_ context.Context) error {
	s.a = "AfterSingletonsInstantiated"
	return nil
}

func TestSmartInstantiateSingleton(t *testing.T) {
	g := NewWithT(t)
	br := NewBeanFactory().(*beanFactoryImpl)
	err := br.RegisterBeanDefinition("smart_singleton", MustNewBeanDefinition(
		reflect.TypeOf((*smartSingleton)(nil)),
	))
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.PreInstantiateSingletons(context.Background())
	g.Expect(err).ShouldNot(HaveOccurred())
	smartObj := br.singletonObjects["smart_singleton"]
	g.Expect(smartObj.Interface().(*smartSingleton).a).To(Equal("AfterSingletonsInstantiated"))
}

type BeanA struct {
	BeanB *BeanB `airmid:"autowire:?"`
	a     string `airmid:"value:${a:=}"`
}

type BeanB struct {
	BeanA []*BeanA `airmid:"autowire:?"`
	b     string   `airmid:"value:${b:=}"`
}

func TestAutowireCircularly_singleton(t *testing.T) {
	g := NewWithT(t)
	br := NewBeanFactory()

	err := br.Set(context.Background(), "a", "aaa")
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.Set(context.Background(), "b", "bbb")
	g.Expect(err).ShouldNot(HaveOccurred())

	err = br.RegisterBeanDefinition("BeanA", MustNewBeanDefinition(reflect.TypeOf((*BeanA)(nil))))
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.RegisterBeanDefinition("BeanB", MustNewBeanDefinition(reflect.TypeOf((*BeanB)(nil))))
	g.Expect(err).ShouldNot(HaveOccurred())

	a, err := br.GetBean(context.Background(), "BeanA")
	g.Expect(err).ShouldNot(HaveOccurred())

	objA := a.(*BeanA)
	objB := objA.BeanB
	g.Expect(objA.a).To(Equal("aaa"))
	g.Expect(objB.BeanA[0].a).To(Equal("aaa"))
	g.Expect(objB.b).To(Equal("bbb"))
}

func TestAutowireCircularly_prototype(t *testing.T) {
	g := NewWithT(t)
	br := NewBeanFactory()

	err := br.Set(context.Background(), "a", "aaa")
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.Set(context.Background(), "b", "bbb")
	g.Expect(err).ShouldNot(HaveOccurred())

	err = br.RegisterBeanDefinition("BeanA", MustNewBeanDefinition(reflect.TypeOf((*BeanA)(nil)), WithBeanScope(ScopePrototype)))
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.RegisterBeanDefinition("BeanB", MustNewBeanDefinition(reflect.TypeOf((*BeanB)(nil)), WithBeanScope(ScopePrototype)))
	g.Expect(err).ShouldNot(HaveOccurred())

	_, err = br.GetBean(context.Background(), "BeanA")
	g.Expect(err).Should(HaveOccurred())
	g.Expect(err.Error()).Should(Equal("cannot get bean 'BeanA' circularly"))
}

func TestGetBeanConcurrently_singleton(t *testing.T) {
	g := NewWithT(t)
	br := NewBeanFactory()

	err := br.Set(context.Background(), "a", "aaa")
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.Set(context.Background(), "b", "bbb")
	g.Expect(err).ShouldNot(HaveOccurred())

	err = br.RegisterBeanDefinition("BeanA", MustNewBeanDefinition(reflect.TypeOf((*BeanA)(nil))))
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.RegisterBeanDefinition("BeanB", MustNewBeanDefinition(reflect.TypeOf((*BeanB)(nil))))
	g.Expect(err).ShouldNot(HaveOccurred())

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			objA, err := br.GetBean(context.Background(), "BeanA")
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objA.(*BeanA).a).To(Equal("aaa"))
			g.Expect(objA.(*BeanA).BeanB.BeanA[0].a).To(Equal("aaa"))
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGetBeanConcurrently_prototype(t *testing.T) {
	g := NewWithT(t)
	br := NewBeanFactory()

	err := br.Set(context.Background(), "a", "aaa")
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.Set(context.Background(), "b", "bbb")
	g.Expect(err).ShouldNot(HaveOccurred())

	err = br.RegisterBeanDefinition("BeanA", MustNewBeanDefinition(reflect.TypeOf((*BeanA)(nil)), WithBeanScope(ScopePrototype)))
	g.Expect(err).ShouldNot(HaveOccurred())
	err = br.RegisterBeanDefinition("BeanB", MustNewBeanDefinition(reflect.TypeOf((*BeanB)(nil))))
	g.Expect(err).ShouldNot(HaveOccurred())

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			res, err := br.GetBean(context.Background(), "BeanA")
			g.Expect(res).To(BeNil())
			g.Expect(err).Should(HaveOccurred())
			g.Expect(err.Error()).Should(Equal("cannot get bean 'BeanA' circularly"))
			wg.Done()
		}()
	}
	wg.Wait()
}

type testDestruction struct {
	str string `airmid:"value:${test.destruction:=test}"`
}

func (t *testDestruction) PostProcessBeforeDestruction(beanName string, bean any) {
	t.str = ""
}

func TestBeanFactory_Destroy(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory().(*beanFactoryImpl)
	bf.RegisterBeanDefinition("test-destruction", MustNewBeanDefinition(reflect.TypeOf((*testDestruction)(nil))))
	obj, err := bf.GetBean(context.Background(), "test-destruction")
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(obj.(*testDestruction).str).To(Equal("test"))
	bf.Destroy()
	obj = bf.allBeans["test-destruction"]
	g.Expect(obj.(*testDestruction).str).To(Equal(""))
}

type testInterface interface {
	Test() string
}

type impl1 struct {
	str string `airmid:"value:${test.primary:=primary}"`
}

func (i1 *impl1) Test() string {
	return i1.str
}

type impl2 struct {
	str string `airmid:"value:${test.primary:=normal}"`
}

func (i2 *impl2) Test() string {
	return i2.str
}

type testNotImplementInterface interface { //nolint
	NotImplement() string
}

type impl3 struct {
	str                   string                    `airmid:"value:${test.primary:=NotImplement}"`
	notImplementInterface testNotImplementInterface `airmid:"autowire:?"` //nolint
}

func (i3 *impl3) Test() string {
	return i3.str
}

type testBean struct {
	test testInterface `airmid:"autowire:?"`
}

func TestAutowireMultipleCandidatesWithOnePrimary(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory().(*beanFactoryImpl)
	bf.RegisterBeanDefinition("test-bean", MustNewBeanDefinition(reflect.TypeOf((*testBean)(nil))))
	bf.RegisterBeanDefinition("test-impl1", MustNewBeanDefinition(reflect.TypeOf((*impl1)(nil)), WithPrimary()))
	bf.RegisterBeanDefinition("test-impl2", MustNewBeanDefinition(reflect.TypeOf((*impl2)(nil))))
	obj, err := bf.GetBean(context.Background(), "test-bean")
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(obj.(*testBean).test.Test()).To(Equal("primary"))
}

func TestAutowireMultipleCandidatesWithMultiplePrimary(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory().(*beanFactoryImpl)
	bf.RegisterBeanDefinition("test-bean", MustNewBeanDefinition(reflect.TypeOf((*testBean)(nil))))
	bf.RegisterBeanDefinition("test-impl1", MustNewBeanDefinition(reflect.TypeOf((*impl1)(nil)), WithPrimary()))
	bf.RegisterBeanDefinition("test-impl2", MustNewBeanDefinition(reflect.TypeOf((*impl2)(nil)), WithPrimary()))
	_, err := bf.GetBean(context.Background(), "test-bean")
	g.Expect(err).Should(HaveOccurred())
	g.Expect(err.Error()).To(Equal("'2' primary candidates found for field 'test' with type ioc.testInterface"))
}

func TestAutowireMultipleCandidatesWithNotImplementInterfaceLazyLoadBean(t *testing.T) {
	g := NewWithT(t)
	bf := NewBeanFactory().(*beanFactoryImpl)
	bf.RegisterBeanDefinition("test-bean", MustNewBeanDefinition(reflect.TypeOf((*testBean)(nil))))
	bf.RegisterBeanDefinition("test-impl1", MustNewBeanDefinition(reflect.TypeOf((*impl1)(nil)), WithLazyMode()))
	bf.RegisterBeanDefinition("test-impl3", MustNewBeanDefinition(reflect.TypeOf((*impl3)(nil)), WithLazyMode()))
	_, err := bf.GetBean(context.Background(), "test-bean")
	g.Expect(err).Should(HaveOccurred())
	g.Expect(err.Error()).To(Equal("'2' candidates found for field 'test' with type ioc.testInterface"))
}
