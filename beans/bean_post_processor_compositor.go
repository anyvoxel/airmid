package beans

// BeanPostProcessorCompositor is the compositor of BeanPostProcessor.
type BeanPostProcessorCompositor interface {
	// The compositor should implement the BeanPostProcessor
	BeanPostProcessor

	// The compositor should implement the BeanDefinitionPostProcessor
	BeanDefinitionPostProcessor

	// AddBeanPostProcessor will add beanPostProcessor to compositor
	AddBeanPostProcessor(BeanPostProcessor)
}

// beanPostProcessorCompositorImpl is the implement of BeanPostProcessorCompositor.
type beanPostProcessorCompositorImpl struct {
	// beanPostProcessors is the array of BeanPostProcessor
	beanPostProcessors []BeanPostProcessor

	// beanDefinitionPostProcessors is the array of BeanDefinitionPostProcessor
	beanDefinitionPostProcessors []BeanDefinitionPostProcessor
}

// NewBeanPostProcessorCompositor return the implement of BeanPostProcessorCompositor.
func NewBeanPostProcessorCompositor(beanPostProcessors ...BeanPostProcessor) BeanPostProcessorCompositor {
	compositor := &beanPostProcessorCompositorImpl{
		beanPostProcessors:           make([]BeanPostProcessor, 0, len(beanPostProcessors)),
		beanDefinitionPostProcessors: make([]BeanDefinitionPostProcessor, 0, len(beanPostProcessors)),
	}
	for _, beanPostProcessor := range beanPostProcessors {
		compositor.AddBeanPostProcessor(beanPostProcessor)
	}

	return compositor
}

func (p *beanPostProcessorCompositorImpl) PostProcessBeforeInitialization(obj any, beanName string) (v any, err error) {
	for _, beanPostProcessor := range p.beanPostProcessors {
		obj, err = beanPostProcessor.PostProcessBeforeInitialization(obj, beanName)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func (p *beanPostProcessorCompositorImpl) PostProcessAfterInitialization(obj any, beanName string) (v any, err error) {
	for _, beanPostProcessor := range p.beanPostProcessors {
		obj, err = beanPostProcessor.PostProcessAfterInitialization(obj, beanName)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func (p *beanPostProcessorCompositorImpl) PostProcessBeanDefinition(beanName string, beanDefinition BeanDefinition) {
	for _, beanDefinitionPostProcessor := range p.beanDefinitionPostProcessors {
		beanDefinitionPostProcessor.PostProcessBeanDefinition(beanName, beanDefinition)
	}
}

func (p *beanPostProcessorCompositorImpl) AddBeanPostProcessor(beanPostProcessor BeanPostProcessor) {
	// TODO: optimize this
	arr := make([]BeanPostProcessor, 0, len(p.beanPostProcessors)+1)
	for _, v := range p.beanPostProcessors {
		if v != beanPostProcessor {
			arr = append(arr, v)
		}
	}
	arr = append(arr, beanPostProcessor)
	p.beanPostProcessors = arr

	p.AddBeanDefinitionPostProcessor(beanPostProcessor)
}

func (p *beanPostProcessorCompositorImpl) AddBeanDefinitionPostProcessor(beanPostProcessor BeanPostProcessor) {
	beanDefinitionPostProcessor := IndirectTo[BeanDefinitionPostProcessor](beanPostProcessor)
	if beanDefinitionPostProcessor == nil {
		return
	}

	// TODO: optimize this
	arr := make([]BeanDefinitionPostProcessor, 0, len(p.beanDefinitionPostProcessors)+1)
	for _, v := range p.beanDefinitionPostProcessors {
		if v != beanDefinitionPostProcessor {
			arr = append(arr, v)
		}
	}
	arr = append(arr, beanDefinitionPostProcessor)
	p.beanDefinitionPostProcessors = arr
}
