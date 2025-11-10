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
	"testing"

	. "github.com/onsi/gomega"
)

type testBeanNameAware struct {
	name string
}

func (a *testBeanNameAware) SetBeanName(name string) {
	a.name = name
}

type testBeanFactoryAware struct {
	beanFactory BeanFactory
}

func (a *testBeanFactoryAware) SetBeanFactory(beanFactory BeanFactory) {
	a.beanFactory = beanFactory
}

func TestBeanAwarePostProcessBeforeInitialization(t *testing.T) {
	type testCase struct {
		desp   string
		obj    any
		name   string
		expect any
	}
	testCases := []testCase{
		{
			desp:   "not aware",
			obj:    "1",
			name:   "n1",
			expect: "1",
		},
		{
			desp:   "bean name aware",
			obj:    &testBeanNameAware{name: "1"},
			name:   "n1",
			expect: &testBeanNameAware{name: "n1"},
		},
		{
			desp: "bean factory aware",
			obj:  &testBeanFactoryAware{},
			name: "n1",
			expect: &testBeanFactoryAware{
				beanFactory: NewBeanFactory(),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			p := &BeanAwareProcessor{
				beanFactory: NewBeanFactory(),
			}
			_, err := p.PostProcessBeforeInitialization(tc.obj, tc.name)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.obj).To(Equal(tc.expect))
		})
	}
}

func TestBeanAwarePostProcessAfterInitialization(t *testing.T) {
	type testCase struct {
		desp   string
		obj    any
		name   string
		expect any
	}
	testCases := []testCase{
		{
			desp:   "not aware",
			obj:    "1",
			name:   "n1",
			expect: "1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			p := &BeanAwareProcessor{}
			_, err := p.PostProcessBeforeInitialization(tc.obj, tc.name)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.obj).To(Equal(tc.expect))
		})
	}
}
