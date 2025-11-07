package beans

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
