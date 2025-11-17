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
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAddValue(t *testing.T) {
	g := NewWithT(t)
	p := NewPropertyValues().(*propertyValuesImpl)

	g.Expect(p.values).To(Equal(make(map[int]reflect.Value)))
	p.AddValue(1, reflect.ValueOf(g))
	g.Expect(p.values).To(Equal(map[int]reflect.Value{
		1: reflect.ValueOf(g),
	}))
}

func TestSetProperty(t *testing.T) {
	g := NewWithT(t)
	p := NewPropertyValues()

	type testBean struct {
		v1 int //nolint
		v2 string
		V3 []int
	}
	p.AddValue(1, reflect.ValueOf("v2"))
	p.AddValue(2, reflect.ValueOf([]int{0, 1}))
	b := &testBean{}
	p.SetProperty(context.Background(), reflect.ValueOf(b).Elem(), []FieldDescriptor{
		{
			FieldIndex: 0,
			Name:       "v1",
			Typ:        reflect.TypeOf(int(0)),
			Unexported: true,
		},
		{
			FieldIndex: 1,
			Name:       "v2",
			Typ:        reflect.TypeOf(""),
			Unexported: true,
		},
		{
			FieldIndex: 2,
			Name:       "V3",
			Typ:        reflect.TypeOf([]int{}),
			Unexported: false,
		},
	})
	g.Expect(b).To(Equal(&testBean{
		v2: "v2",
		V3: []int{0, 1},
	}))
}
