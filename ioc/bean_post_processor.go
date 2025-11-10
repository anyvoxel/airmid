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
