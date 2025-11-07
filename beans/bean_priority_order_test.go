package beans

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
