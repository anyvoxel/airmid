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
