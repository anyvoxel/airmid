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

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

func TestReflectConstructor(t *testing.T) {
	g := NewWithT(t)
	actual, err := (&ReflectConstructor{typ: reflect.TypeOf(t)}).NewObject(nil)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(actual.Interface()).To(Equal(&testing.T{}))
}

type testMethodConstructorStruct struct{}

func (t *testMethodConstructorStruct) F1() int {
	return 1
}

func (t *testMethodConstructorStruct) F2WithError() (int, error) {
	return 2, xerrors.ErrContinue
}

func (t *testMethodConstructorStruct) F2WithNilError() (int, error) {
	return 3, nil
}

func (t *testMethodConstructorStruct) F3WithArgs(x int, y int) int {
	return x + y
}

type testFnConstructorArgumentResolver struct {
	fn func(args []ConstructorArgument) ([]reflect.Value, error)
}

func (t *testFnConstructorArgumentResolver) Resolve(args []ConstructorArgument) ([]reflect.Value, error) {
	return t.fn(args)
}

func TestMethodConstructor(t *testing.T) {
	type testCase struct {
		desp     string
		typ      reflect.Type
		method   reflect.Method
		resolver ConstructorArgumentResolver

		err    string
		expect any
	}
	testCases := []testCase{
		{
			desp:   "output 1",
			typ:    reflect.TypeOf((*testMethodConstructorStruct)(nil)),
			method: reflect.TypeOf((*testMethodConstructorStruct)(nil)).Method(0),
			resolver: &testFnConstructorArgumentResolver{
				fn: func(args []ConstructorArgument) ([]reflect.Value, error) {
					return nil, nil
				},
			},
			err:    "",
			expect: 1,
		},
		{
			desp:   "output 2 with error",
			typ:    reflect.TypeOf((*testMethodConstructorStruct)(nil)),
			method: reflect.TypeOf((*testMethodConstructorStruct)(nil)).Method(1),
			resolver: &testFnConstructorArgumentResolver{
				fn: func(args []ConstructorArgument) ([]reflect.Value, error) {
					return nil, nil
				},
			},
			err:    "Continue",
			expect: 2,
		},
		{
			desp:   "output 2 with nil error",
			typ:    reflect.TypeOf((*testMethodConstructorStruct)(nil)),
			method: reflect.TypeOf((*testMethodConstructorStruct)(nil)).Method(2),
			resolver: &testFnConstructorArgumentResolver{
				fn: func(args []ConstructorArgument) ([]reflect.Value, error) {
					return nil, nil
				},
			},
			err:    "",
			expect: 3,
		},
		{
			desp:   "output 5 with args",
			typ:    reflect.TypeOf((*testMethodConstructorStruct)(nil)),
			method: reflect.TypeOf((*testMethodConstructorStruct)(nil)).Method(3),
			resolver: &testFnConstructorArgumentResolver{
				fn: func(args []ConstructorArgument) ([]reflect.Value, error) {
					return []reflect.Value{reflect.ValueOf(2), reflect.ValueOf(3)}, nil
				},
			},
			err:    "",
			expect: 5,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			m := &MethodConstructor{
				typ:    tc.typ,
				method: tc.method,
			}
			actual, err := m.NewObject(tc.resolver)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
			} else {
				g.Expect(err).ToNot(HaveOccurred())
			}

			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}
