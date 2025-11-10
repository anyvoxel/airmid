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
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestRegisterBeanDefinition(t *testing.T) {
	g := NewWithT(t)
	r := NewBeanDefinitionRegistry()

	err := r.RegisterBeanDefinition("bean1", nil)
	g.Expect(err).ToNot(HaveOccurred())

	err = r.RegisterBeanDefinition("bean1", nil)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp("Cannot register bean 'bean1'"))
}

func TestRemoveBeanDefinition(t *testing.T) {
	g := NewWithT(t)
	r := NewBeanDefinitionRegistry()

	err := r.RegisterBeanDefinition("bean1", nil)
	g.Expect(err).ToNot(HaveOccurred())

	err = r.RemoveBeanDefinition("bean1")
	g.Expect(err).ToNot(HaveOccurred())

	err = r.RemoveBeanDefinition("bean1")
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp("No bean 'bean1' registered"))
}

func TestGetBeanDefinition(t *testing.T) {
	g := NewWithT(t)
	r := NewBeanDefinitionRegistry()
	bean := &beanDefinitionHolder{
		name: "bean1",
	}

	err := r.RegisterBeanDefinition("bean1", bean)
	g.Expect(err).ToNot(HaveOccurred())

	actual, err := r.GetBeanDefinition("bean1")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(actual).To(Equal(bean))

	actual, err = r.GetBeanDefinition("bean2")
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(MatchRegexp("No bean 'bean2' registered"))
	g.Expect(actual).To(BeNil())
}

func TestVisitBeanDefinition(t *testing.T) {
	g := NewWithT(t)
	r := NewBeanDefinitionRegistry()
	beans := []BeanDefinition{
		&beanDefinitionHolder{
			name: "bean1",
			Typ:  reflect.TypeOf(int(0)),
		},
		&beanDefinitionHolder{
			name: "bean2",
			Typ:  reflect.TypeOf(""),
		},
	}

	expect := map[string]reflect.Type{
		"bean1": reflect.TypeOf(int(0)),
		"bean2": reflect.TypeOf(""),
	}

	for _, bean := range beans {
		err := r.RegisterBeanDefinition(bean.Name(), bean)
		g.Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("register bean '%v'", bean.Name()))
	}

	actual := map[string]reflect.Type{}
	r.VisitBeanDefinition(FuncVisitor{
		VisitFunc: func(s string, bd BeanDefinition) {
			actual[s] = bd.Type()
		},
	})
	g.Expect(actual).To(Equal(expect))
}
