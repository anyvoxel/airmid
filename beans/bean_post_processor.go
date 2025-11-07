package beans

// BeanPostProcessor hook that allows for custom modification of new bean instances.
// For example, checking for marker interfaces or wrapping beans with proxies.
type BeanPostProcessor interface {
	// PostProcessBeforeInitialization will been apply to the given new bean instance
	// before any bean initialization callbacks.
	// It will called before AfterPropertiesSet & custom init-method,
	// and the bean will already be populated with property values.
	PostProcessBeforeInitialization(obj any, beanName string) (v any, err error)

	// PostProcessAfterInitialization will been apply to the given new bean instance
	// after any bean initialization callbacks.
	PostProcessAfterInitialization(obj any, beanName string) (v any, err error)
}

// FuncBeanPostProcessor is the processor for compose function.
// NOTE: this is used for unit test, be careful to use it.
type FuncBeanPostProcessor struct {
	BeforeInitializationFunc func(obj any, beanName string) (v any, err error)
	AfterInitializationFunc  func(obj any, beanName string) (v any, err error)
}

var _ BeanPostProcessor = &FuncBeanPostProcessor{}

// PostProcessBeforeInitialization implement BeanPostProcessor.PostProcessBeforeInitialization.
func (p *FuncBeanPostProcessor) PostProcessBeforeInitialization(obj any, beanName string) (v any, err error) {
	if p.BeforeInitializationFunc != nil {
		return p.BeforeInitializationFunc(obj, beanName)
	}

	return obj, nil
}

// PostProcessAfterInitialization implement BeanPostProcessor.PostProcessAfterInitialization.
func (p *FuncBeanPostProcessor) PostProcessAfterInitialization(obj any, beanName string) (v any, err error) {
	if p.AfterInitializationFunc != nil {
		return p.AfterInitializationFunc(obj, beanName)
	}

	return obj, nil
}

// DestructionAwareBeanPostProcessor is the processor which will be called before bean destruction.
type DestructionAwareBeanPostProcessor interface {
	// PostProcessBeforeDestruction will be called before bean destruction
	PostProcessBeforeDestruction(beanName string, bean any)
}
