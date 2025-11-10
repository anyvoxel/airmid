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

	"github.com/anyvoxel/airmid/anvil/xreflect"
)

type testDefaultConstructorBuilderBuild0Struct struct {
}

type testDefaultConstructorBuilderBuild1Struct struct {
}

func (t *testDefaultConstructorBuilderBuild1Struct) NewtestDefaultConstructorBuilderBuild1Struct() *testDefaultConstructorBuilderBuild1Struct {
	return nil
}

type testDefaultConstructorBuilderBuild2Struct struct {
}

func (t *testDefaultConstructorBuilderBuild2Struct) NewtestDefaultConstructorBuilderBuild2Struct() *testDefaultConstructorBuilderBuild2Struct {
	return nil
}

func (t *testDefaultConstructorBuilderBuild2Struct) NewTestDefaultConstructorBuilderBuild2Struct() *testDefaultConstructorBuilderBuild2Struct {
	return nil
}

func TestDefaultConstructorBuilderBuild(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		args   []ConstructorArgument
		err    string
		expect func(g *WithT, c Constructor)
	}
	testCases := []testCase{
		{
			desp: "non method",
			typ:  reflect.TypeOf((*testDefaultConstructorBuilderBuild0Struct)(nil)),
			err:  "",
			expect: func(g *WithT, c Constructor) {
				g.Expect(c).To(Equal(
					&ReflectConstructor{
						typ: reflect.TypeOf((*testDefaultConstructorBuilderBuild0Struct)(nil)),
					},
				))
			},
		},
		{
			desp: "one method",
			typ:  reflect.TypeOf((*testDefaultConstructorBuilderBuild1Struct)(nil)),
			err:  "",
			expect: func(g *WithT, c Constructor) {
				mc := c.(*MethodConstructor)
				g.Expect(mc.typ).To(Equal(reflect.TypeOf((*testDefaultConstructorBuilderBuild1Struct)(nil))))
				g.Expect(mc.method.Name).To(Equal("NewtestDefaultConstructorBuilderBuild1Struct"))
			},
		},
		{
			desp:   "many method",
			typ:    reflect.TypeOf((*testDefaultConstructorBuilderBuild2Struct)(nil)),
			err:    "Too many function '2'",
			expect: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			b := NewDefaultConstructorBuilder()

			actual, err := b.Build(tc.typ, tc.args)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			tc.expect(g, actual)
		})
	}
}

func TestMethodNameFilter(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		method reflect.Method
		args   []ConstructorArgument
		expect bool
	}
	testCases := []testCase{
		{
			desp: "normal test",
			typ:  reflect.TypeOf(t),
			method: reflect.Method{
				Name: "NewT",
			},
			expect: true,
		},
		{
			desp: "in-case sensitive match",
			typ:  reflect.TypeOf(t),
			method: reflect.Method{
				Name: "Newt",
			},
			expect: true,
		},
		{
			desp: "prefix mismatch",
			typ:  reflect.TypeOf(t),
			method: reflect.Method{
				Name: "newT",
			},
			expect: false,
		},
		{
			desp: "method name mismatch",
			typ:  reflect.TypeOf(t),
			method: reflect.Method{
				Name: "Newg",
			},
			expect: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := MethodNameFilter(tc.typ, tc.method, tc.args)
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestMethodInputNumFilter(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		method reflect.Method
		args   []ConstructorArgument
		expect bool
	}
	testCases := []testCase{
		{
			desp: "normal test",
			method: reflect.Method{
				Func: reflect.ValueOf(TestMethodInputNumFilter),
			},
			expect: true,
		},
		{
			desp: "with args",
			method: reflect.Method{
				Func: reflect.ValueOf(MethodInputNumFilter),
			},
			args:   []ConstructorArgument{{}, {}},
			expect: true,
		},
		{
			desp: "num mismatch",
			method: reflect.Method{
				Func: reflect.ValueOf(MethodInputNumFilter),
			},
			expect: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := MethodInputNumFilter(tc.typ, tc.method, tc.args)
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestMethodOutputNumFilter(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		method reflect.Method
		args   []ConstructorArgument
		expect bool
	}
	testCases := []testCase{
		{
			desp: "normal 1 test",
			method: reflect.Method{
				Func: reflect.ValueOf(MethodOutputNumFilter),
			},
			expect: true,
		},
		{
			desp: "normal 2 test",
			method: reflect.Method{
				Func: reflect.ValueOf(xreflect.NewValue),
			},
			expect: true,
		},
		{
			desp: "num 0 mismatch",
			method: reflect.Method{
				Func: reflect.ValueOf(TestMethodOutputNumFilter),
			},
			expect: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := MethodOutputNumFilter(tc.typ, tc.method, tc.args)
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

type testMethodOutput0FilterStruct struct {
}

func (*testMethodOutput0FilterStruct) Match() *testMethodOutput0FilterStruct {
	return nil
}

func (*testMethodOutput0FilterStruct) MisMatch() testMethodOutput0FilterStruct {
	return testMethodOutput0FilterStruct{}
}

func TestMethodOutput0Filter(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		method reflect.Method
		args   []ConstructorArgument
		expect bool
	}
	testCases := []testCase{
		{
			desp:   "output match test",
			typ:    reflect.TypeOf((*testMethodOutput0FilterStruct)(nil)),
			method: reflect.TypeOf((*testMethodOutput0FilterStruct)(nil)).Method(0),
			expect: true,
		},
		{
			desp:   "output mismatch test",
			typ:    reflect.TypeOf((*testMethodOutput0FilterStruct)(nil)),
			method: reflect.TypeOf((*testMethodOutput0FilterStruct)(nil)).Method(1),
			expect: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := MethodOutput0Filter(tc.typ, tc.method, tc.args)
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

type testMethodOutput1FilterStruct struct {
}

func (*testMethodOutput1FilterStruct) Non2() *testMethodOutput0FilterStruct {
	return nil
}

func (*testMethodOutput1FilterStruct) Output1Match() (*testMethodOutput0FilterStruct, error) {
	return nil, nil
}

func (*testMethodOutput1FilterStruct) Output1MisMatch() (*testMethodOutput0FilterStruct, *error) {
	return nil, nil
}

func TestMethodOutput1Filter(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		method reflect.Method
		args   []ConstructorArgument
		expect bool
	}
	testCases := []testCase{
		{
			desp:   "output non-2 test",
			typ:    reflect.TypeOf((*testMethodOutput1FilterStruct)(nil)),
			method: reflect.TypeOf((*testMethodOutput1FilterStruct)(nil)).Method(0),
			expect: true,
		},
		{
			desp:   "output match test",
			typ:    reflect.TypeOf((*testMethodOutput1FilterStruct)(nil)),
			method: reflect.TypeOf((*testMethodOutput1FilterStruct)(nil)).Method(1),
			expect: true,
		},
		{
			desp:   "output mismatch test",
			typ:    reflect.TypeOf((*testMethodOutput1FilterStruct)(nil)),
			method: reflect.TypeOf((*testMethodOutput1FilterStruct)(nil)).Method(2),
			expect: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := MethodOutput1Filter(tc.typ, tc.method, tc.args)
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestMethodFiltersFilter(t *testing.T) {
	type testCase struct {
		desp   string
		F      MethodFilters
		expect bool
	}
	testCases := []testCase{
		{
			desp: "all-1 pass",
			F: []MethodFilter{
				&FnMethodFilter{
					FnFilter: func(t reflect.Type, m reflect.Method, args []ConstructorArgument) bool {
						return true
					},
				},
			},
			expect: true,
		},
		{
			desp: "all-2 pass",
			F: []MethodFilter{
				&FnMethodFilter{
					FnFilter: func(t reflect.Type, m reflect.Method, args []ConstructorArgument) bool {
						return true
					},
				},
				&FnMethodFilter{
					FnFilter: func(t reflect.Type, m reflect.Method, args []ConstructorArgument) bool {
						return true
					},
				},
			},
			expect: true,
		},
		{
			desp: "one failed",
			F: []MethodFilter{
				&FnMethodFilter{
					FnFilter: func(t reflect.Type, m reflect.Method, args []ConstructorArgument) bool {
						return true
					},
				},
				&FnMethodFilter{
					FnFilter: func(t reflect.Type, m reflect.Method, args []ConstructorArgument) bool {
						return false
					},
				},
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := tc.F.Filter(nil, reflect.Method{}, nil)
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}
