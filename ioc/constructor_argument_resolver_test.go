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
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/anvil/pointer"
)

func TestFactoryResolverResolve(t *testing.T) {
	type testCase struct {
		desp     string
		bf       *beanFactoryImpl
		args     []ConstructorArgument
		initFunc func(g *WithT, tc *testCase, bf BeanFactory)

		expect []any
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal test",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			args: []ConstructorArgument{
				{
					Property: &PropertyFieldDescriptor{
						Name: "p1",
					},
					Type: reflect.TypeOf(""),
				},
				{
					Property: &PropertyFieldDescriptor{
						Name:    "p2",
						Default: pointer.StringPtr("2"),
					},
					Type: reflect.TypeOf(""),
				},
				{
					Bean: &BeanFieldDescriptor{
						Name: "b1",
					},
					Type: reflect.TypeOf(""),
				},
				{
					Value: reflect.ValueOf("value"),
					Type:  reflect.TypeOf(""),
				},
			},
			initFunc: func(g *WithT, tc *testCase, bf BeanFactory) {
				bf.Set("p1", "p1")
				bf.RegisterSingleton("b1", "b1")
			},
			expect: []any{
				"p1",
				"2",
				"b1",
				"value",
			},
			err: "",
		},
		{
			desp: "property not found",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			args: []ConstructorArgument{
				{
					Property: &PropertyFieldDescriptor{
						Name: "p1",
					},
					Type: reflect.TypeOf(""),
				},
				{
					Property: &PropertyFieldDescriptor{
						Name:    "p2",
						Default: pointer.StringPtr("2"),
					},
					Type: reflect.TypeOf(""),
				},
				{
					Bean: &BeanFieldDescriptor{
						Name: "b1",
					},
					Type: reflect.TypeOf(""),
				},
				{
					Value: reflect.ValueOf("value"),
					Type:  reflect.TypeOf(""),
				},
			},
			initFunc: func(g *WithT, tc *testCase, bf BeanFactory) {
				bf.Set("p3", "p1")
				bf.RegisterSingleton("b1", "b1")
			},
			expect: []any{},
			err:    `property with key='p1' not found`,
		},
		{
			desp: "bean not found",
			bf:   NewBeanFactory().(*beanFactoryImpl),
			args: []ConstructorArgument{
				{
					Property: &PropertyFieldDescriptor{
						Name: "p1",
					},
					Type: reflect.TypeOf(""),
				},
				{
					Property: &PropertyFieldDescriptor{
						Name:    "p2",
						Default: pointer.StringPtr("2"),
					},
					Type: reflect.TypeOf(""),
				},
				{
					Bean: &BeanFieldDescriptor{
						Name: "b1",
					},
					Type: reflect.TypeOf(""),
				},
				{
					Value: reflect.ValueOf("value"),
					Type:  reflect.TypeOf(""),
				},
			},
			initFunc: func(g *WithT, tc *testCase, bf BeanFactory) {
				bf.Set("p1", "p1")
			},
			expect: []any{},
			err:    "No bean 'b1' registered: ObjectNotFound",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			tc.initFunc(g, &tc, tc.bf)

			r := &factoryConstructorArgumentResolver{
				bf: tc.bf,
			}
			actual, err := r.Resolve(tc.args)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			for i := 0; i < len(actual); i++ {
				g.Expect(actual[i].Interface()).To(Equal(tc.expect[i]))
			}
		})
	}
}
