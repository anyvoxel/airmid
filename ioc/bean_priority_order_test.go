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
	"sort"
	"testing"

	"github.com/onsi/gomega"
)

type TestOrder interface {
	GetName() string
}

var globPriority BeanPriority

type OrderBean struct {
	priority BeanPriority
}

func (o *OrderBean) GetPriority() BeanPriority {
	return o.priority
}

func (o *OrderBean) GetName() string {
	return "implement BeanPriorityOrder"
}

func (o *OrderBean) AfterPropertiesSet() error {
	globPriority++
	o.priority = globPriority
	return nil
}

type OrderBean2 struct{}

func (o *OrderBean2) GetName() string {
	return "unimplement BeanPriorityOrder"
}

func TestBeanFactory_OrderBean(t *testing.T) {
	g := gomega.NewWithT(t)
	ob1 := &OrderBean{
		priority: func() BeanPriority {
			globPriority++
			return globPriority
		}(),
	}
	ob2 := &OrderBean{
		priority: func() BeanPriority {
			globPriority++
			return globPriority
		}(),
	}
	ob3 := &OrderBean{
		priority: func() BeanPriority {
			globPriority++
			return globPriority
		}(),
	}
	ob4 := &OrderBean2{}

	obs := make([]reflect.Value, 0, 4)
	// add in a inverted order
	obs = append(obs, reflect.ValueOf(ob4))
	obs = append(obs, reflect.ValueOf(ob3))
	obs = append(obs, reflect.ValueOf(ob2))
	obs = append(obs, reflect.ValueOf(ob1))

	var candidates CandidateBeans = obs
	sort.Sort(candidates)

	for i, c := range candidates {
		obj := c.Interface().(TestOrder)
		if i < 3 {
			g.Expect("implement BeanPriorityOrder").To(gomega.Equal(obj.GetName()))
			actualObj := c.Interface().(*OrderBean)
			g.Expect(actualObj.GetPriority()).To(gomega.Equal(BeanPriority(3 - i)))
		} else {
			g.Expect("unimplement BeanPriorityOrder").To(gomega.Equal(obj.GetName()))
		}
	}
}
