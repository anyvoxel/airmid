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

func TestValidate(t *testing.T) {
	type testCase struct {
		desp string
		o    *beanDefinitionOption
		err  string
	}
	testCases := []testCase{
		{
			desp: "default option test",
			o:    defaultBeanDefinitionOption(),
			err:  "",
		},
		{
			desp: "singleton scope",
			o: &beanDefinitionOption{
				scope: ScopeSingleton,
			},
			err: "",
		},
		{
			desp: "prototype scope",
			o: &beanDefinitionOption{
				scope: ScopePrototype,
			},
			err: "",
		},
		{
			desp: "invalid scope",
			o: &beanDefinitionOption{
				scope: "1",
			},
			err: "Unsupport scope '1'",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.o.Validate()
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestWithBeanName(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt.name).To(Equal(""))

	o := WithBeanName("1")
	o.Apply(opt)

	g.Expect(opt.name).To(Equal("1"))
}

func TestWithBeanScope(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt.scope).To(Equal(ScopeSingleton))

	o := WithBeanScope(ScopePrototype)
	o.Apply(opt)

	g.Expect(opt.scope).To(Equal(ScopePrototype))
}

func TestWithLazyMode(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt.lazy).To(BeFalse())

	o := WithLazyMode()
	o.Apply(opt)

	g.Expect(opt.lazy).To(BeTrue())
}

func TestWithPrimary(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt.primary).To(BeFalse())

	o := WithPrimary()
	o.Apply(opt)

	g.Expect(opt.primary).To(BeTrue())
}
