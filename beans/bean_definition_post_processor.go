package beans

// BeanDefinitionPostProcessor is the callback interface for bean definitions at runtime.
type BeanDefinitionPostProcessor interface {
	// The BeanDefinitionPostProcessor must implement the BeanPostProcessor too.
	BeanPostProcessor

	// PostProcessBeanDefinition will been apply to the given bean and beanDefinition.
	PostProcessBeanDefinition(beanName string, beanDefinition BeanDefinition)
}

// FuncBeanDefinitionPostProcessor is the processor for compose function.
// NOTE: this is used for unit test, be careful to use it.
type FuncBeanDefinitionPostProcessor struct {
	BeanPostProcessor

	PostBeanDefinitionFunc func(string, BeanDefinition)
}

var _ BeanDefinitionPostProcessor = &FuncBeanDefinitionPostProcessor{}

// PostProcessBeanDefinition implement BeanDefinitionPostProcessor.PostProcessBeanDefinition.
func (p *FuncBeanDefinitionPostProcessor) PostProcessBeanDefinition(beanName string, beanDefinition BeanDefinition) {
	if p.PostBeanDefinitionFunc != nil {
		p.PostBeanDefinitionFunc(beanName, beanDefinition)
	}
}
