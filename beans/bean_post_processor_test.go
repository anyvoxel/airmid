package beans

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
