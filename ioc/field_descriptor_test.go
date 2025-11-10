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

func TestValueRegexpMatch(t *testing.T) {
	type testCase struct {
		desp   string
		value  string
		expect [][]string
	}
	testCases := []testCase{
		{
			desp:  "normal test",
			value: "${123.123}",
			expect: [][]string{
				{"${123.123}", "123.123"},
			},
		},
		{
			desp:  "normal with default",
			value: "${x.y:=d}",
			expect: [][]string{
				{"${x.y:=d}", "x.y:=d"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			g.Expect(valueRegex.FindAllStringSubmatch(tc.value, -1)).To(Equal(tc.expect))
		})
	}
}

func TestNewFieldDescriptor(t *testing.T) {
	type testCase struct {
		desp   string
		field  reflect.StructField
		idx    int
		expect *FieldDescriptor
		err    string
	}
	testCases := []testCase{
		{
			desp: "with property field",
			field: reflect.StructField{
				Name:    "f1",
				PkgPath: "pkg",
				Tag:     reflect.StructTag(`airmid:"value:${v1:=d1}"`),
			},
			idx: 0,
			expect: &FieldDescriptor{
				FieldIndex: 0,
				Name:       "f1",
				Unexported: true,
				Property: &PropertyFieldDescriptor{
					Name:    "v1",
					Default: pointer.StringPtr("d1"),
				},
			},
		},
		{
			desp: "with bean field",
			field: reflect.StructField{
				Name:    "f1",
				PkgPath: "",
				Tag:     reflect.StructTag(`airmid:"autowire:b1,optional"`),
			},
			idx: 1,
			expect: &FieldDescriptor{
				FieldIndex: 1,
				Name:       "f1",
				Unexported: false,
				Bean: &BeanFieldDescriptor{
					Name:     "b1",
					Optional: true,
				},
			},
		},
		{
			desp: "non field",
			field: reflect.StructField{
				Name:    "f1",
				PkgPath: "",
				Tag:     reflect.StructTag(``),
			},
			idx:    0,
			expect: nil,
		},
		{
			desp: "wrong prefix",
			field: reflect.StructField{
				Name:    "f1",
				PkgPath: "",
				Tag:     reflect.StructTag(`airmid:"v"`),
			},
			idx:    0,
			expect: nil,
			err:    "Invalid tag 'v'",
		},
		{
			desp: "wrong value field",
			field: reflect.StructField{
				Name:    "f1",
				PkgPath: "",
				Tag:     reflect.StructTag(`airmid:"value:1"`),
			},
			idx:    0,
			expect: nil,
			err:    "Invalid value '1'",
		},
		{
			desp: "wrong bean field",
			field: reflect.StructField{
				Name:    "f1",
				PkgPath: "",
				Tag:     reflect.StructTag(`airmid:"autowire:1,vv"`),
			},
			idx:    0,
			expect: nil,
			err:    "Invalid autowire 'vv'",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			fd, err := NewFieldDescriptor(tc.field, tc.idx)

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(fd).To(Equal(tc.expect))
		})
	}
}

func TestNewBeanFieldDescriptor(t *testing.T) {
	type testCase struct {
		desp   string
		value  string
		expect *BeanFieldDescriptor
		err    string
	}
	testCases := []testCase{
		{
			desp:  "only bean name",
			value: "b1",
			expect: &BeanFieldDescriptor{
				Name: "b1",
			},
			err: "",
		},
		{
			desp:  "bean type matcher",
			value: "?",
			expect: &BeanFieldDescriptor{
				Name: "?",
			},
			err: "",
		},
		{
			desp:  "bean name with optional",
			value: "b1,optional",
			expect: &BeanFieldDescriptor{
				Name:     "b1",
				Optional: true,
			},
			err: "",
		},
		{
			desp:   "empty content",
			value:  "",
			expect: nil,
			err:    "Required autowire content",
		},
		{
			desp:   "non optional",
			value:  "b1,required",
			expect: nil,
			err:    "Invalid autowire 'required'",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			fd, err := NewBeanFieldDescriptor(tc.value)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(fd).To(Equal(tc.expect))
		})
	}
}

func TestNewPropertyFieldDescriptor(t *testing.T) {
	type testCase struct {
		desp   string
		value  string
		expect *PropertyFieldDescriptor
		err    string
	}
	testCases := []testCase{
		{
			desp:  "only property name",
			value: "${b1}",
			expect: &PropertyFieldDescriptor{
				Name: "b1",
			},
			err: "",
		},
		{
			desp:  "property with default",
			value: "${b1:=1}",
			expect: &PropertyFieldDescriptor{
				Name:    "b1",
				Default: pointer.StringPtr("1"),
			},
			err: "",
		},
		{
			desp:   "format error",
			value:  "b1",
			expect: nil,
			err:    "Invalid value 'b1'",
		},
		{
			desp:   "empty value",
			value:  "${}",
			expect: nil,
			err:    "Required value content",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			fd, err := NewPropertyFieldDescriptor(tc.value)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(fd).To(Equal(tc.expect))
		})
	}
}
