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

func TestFuncPostProcessBeforeInitialization(t *testing.T) {
	g := NewWithT(t)
	p := &FuncBeanPostProcessor{}
	g.Expect(p.BeforeInitializationFunc).To(BeNil())

	obj := int(1)
	actual, err := p.PostProcessBeforeInitialization(&obj, "b1")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(actual).To(Equal(&obj))

	p.BeforeInitializationFunc = func(obj any, beanName string) (v any, err error) {
		vv := obj.(*int)
		*vv++
		return vv, nil
	}
	actual, err = p.PostProcessBeforeInitialization(&obj, "b1")
	g.Expect(err).ToNot(HaveOccurred())
	obj = int(2)
	g.Expect(actual).To(Equal(&obj))
}

func TestFuncPostProcessAfterInitialization(t *testing.T) {
	g := NewWithT(t)
	p := &FuncBeanPostProcessor{}
	g.Expect(p.AfterInitializationFunc).To(BeNil())

	obj := int(1)
	actual, err := p.PostProcessAfterInitialization(&obj, "b1")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(actual).To(Equal(&obj))

	p.AfterInitializationFunc = func(obj any, beanName string) (v any, err error) {
		vv := obj.(*int)
		*vv++
		return vv, nil
	}
	actual, err = p.PostProcessAfterInitialization(&obj, "b1")
	g.Expect(err).ToNot(HaveOccurred())
	obj = int(2)
	g.Expect(actual).To(Equal(&obj))
}
