package beans

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestFuncPostProcessBeanDefinition(t *testing.T) {
	g := NewWithT(t)
	p := &FuncBeanDefinitionPostProcessor{}

	g.Expect(p.PostBeanDefinitionFunc).To(BeNil())
	beanName := "1"
	p.PostProcessBeanDefinition(beanName, nil)
	g.Expect(beanName).To(Equal("1"))

	p.PostBeanDefinitionFunc = func(s string, bd BeanDefinition) {
		beanName = s + "1"
	}
	p.PostProcessBeanDefinition(beanName, nil)
	g.Expect(beanName).To(Equal("11"))
}
