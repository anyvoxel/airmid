package beans

// BeanAwareProcessor is the processor to invoke aware method, it always execute before any custom BeanPostProcessor.
type BeanAwareProcessor struct {
	beanFactory BeanFactory
}

// PostProcessBeforeInitialization implement the BeanPostProcessor.PostProcessBeforeInitialization.
func (p *BeanAwareProcessor) PostProcessBeforeInitialization(obj any, beanName string) (v any, err error) {
	if vobj := IndirectTo[BeanNameAware](obj); vobj != nil {
		vobj.SetBeanName(beanName)
	}

	if vobj := IndirectTo[BeanFactoryAware](obj); vobj != nil {
		vobj.SetBeanFactory(p.beanFactory)
	}

	return obj, nil
}

// PostProcessAfterInitialization implement the BeanPostProcessor.PostProcessAfterInitialization.
func (*BeanAwareProcessor) PostProcessAfterInitialization(obj any, _ string) (v any, err error) {
	return obj, nil
}
