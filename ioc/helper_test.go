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
)

func TestIsTypeMatched(t *testing.T) {
	type testCase struct {
		desp      string
		targetTyp reflect.Type
		beanTyp   reflect.Type
		expect    bool
	}
	testCases := []testCase{
		{
			desp:      "type assignable",
			targetTyp: reflect.TypeOf(testCase{}),
			beanTyp:   reflect.TypeOf(testCase{}),
			expect:    true,
		},
		{
			desp:      "interface assignable",
			targetTyp: reflect.TypeOf(any(testCase{})),
			beanTyp:   reflect.TypeOf(testCase{}),
			expect:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			g.Expect(IsTypeMatched(tc.targetTyp, tc.beanTyp)).To(Equal(tc.expect))
		})
	}
}

type testNamerI interface {
	Name() string
}

type testNamerI2 interface {
	Name2() string
}

type testNamer1 struct {
	n string
}

func (t *testNamer1) Name() string {
	return t.n
}

func (t *testNamer1) Name2() string {
	return t.n + "/2"
}

type proxyNamer struct {
	testNamerI
}

func (p *proxyNamer) OriginalObject() any {
	return p.testNamerI
}

type nonproxyNamer struct {
	testNamerI
}

func TestIndirectTo(t *testing.T) {
	t.Run("directly implement", func(t *testing.T) {
		g := NewWithT(t)

		v := &testNamer1{
			n: "1",
		}
		actual := IndirectTo[testNamerI](v)
		g.Expect(actual).To(Equal(v))
	})

	t.Run("proxy directly implement", func(t *testing.T) {
		g := NewWithT(t)

		v := &proxyNamer{
			&testNamer1{
				n: "1",
			},
		}
		actual := IndirectTo[testNamerI](v)
		g.Expect(actual).To(Equal(v))
	})

	t.Run("proxy non implement", func(t *testing.T) {
		g := NewWithT(t)

		v := &proxyNamer{
			&testNamer1{
				n: "1",
			},
		}
		actual := IndirectTo[testNamerI2](v)
		g.Expect(actual).To(Equal(v.testNamerI))
	})

	t.Run("non-proxy wrapper", func(t *testing.T) {
		g := NewWithT(t)

		v := &nonproxyNamer{
			&testNamer1{
				n: "1",
			},
		}
		actual := IndirectTo[testNamerI2](v)
		g.Expect(actual).To(BeNil())
	})
}

func TestGetBean(t *testing.T) {
	t.Run("didn't exists", func(t *testing.T) {
		g := NewWithT(t)
		bf := NewBeanFactory()

		o, err := GetBean[*nonproxyNamer](bf, "1")
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`No bean '1' registered: ObjectNotFound`))
		g.Expect(o).To(BeNil())
	})

	t.Run("type mismatch", func(t *testing.T) {
		g := NewWithT(t)
		bf := NewBeanFactory()
		bf.RegisterBeanDefinition("1", MustNewBeanDefinition(reflect.TypeOf((*proxyNamer)(nil))))

		o, err := GetBean[*nonproxyNamer](bf, "1")
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`cannot convert \*ioc\.proxyNamer to \*ioc\.nonproxyNamer`))
		g.Expect(o).To(BeNil())
	})

	t.Run("normal", func(t *testing.T) {
		g := NewWithT(t)
		bf := NewBeanFactory()
		bf.RegisterBeanDefinition("1", MustNewBeanDefinition(reflect.TypeOf((*nonproxyNamer)(nil))))

		o, err := GetBean[*nonproxyNamer](bf, "1")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(o).ToNot(BeNil())
	})
}
