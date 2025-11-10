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
