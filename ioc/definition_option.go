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

import (
	"github.com/anyvoxel/airmid/anvil/xerrors"
)

// BeanDefinitionOption is the configuration helper for build bean definition.
type BeanDefinitionOption interface {
	Apply(*beanDefinitionOption)
}

type fnBeanDefinitionOption struct {
	fn func(*beanDefinitionOption)
}

func (f *fnBeanDefinitionOption) Apply(opt *beanDefinitionOption) {
	f.fn(opt)
}

type beanDefinitionOption struct {
	name    string
	scope   string
	lazy    bool
	primary bool

	construtorArguments []ConstructorArgument
}

// WithBeanName will set the bean definition name.
type WithBeanName string

// Apply will apply the name to bean definition.
func (w WithBeanName) Apply(opt *beanDefinitionOption) {
	opt.name = string(w)
}

// WithBeanScope will set the bean scope.
type WithBeanScope string

// Apply will apply the scope to bean definition.
func (w WithBeanScope) Apply(opt *beanDefinitionOption) {
	opt.scope = string(w)
}

// WithLazyMode will set the bean lazy mode.
func WithLazyMode() BeanDefinitionOption {
	return &fnBeanDefinitionOption{
		fn: func(opt *beanDefinitionOption) {
			opt.lazy = true
		},
	}
}

// WithPrimary will set the bean primary true.
func WithPrimary() BeanDefinitionOption {
	return &fnBeanDefinitionOption{
		fn: func(opt *beanDefinitionOption) {
			opt.primary = true
		},
	}
}

// WithConstructorArguments will set the constructor arguments.
func WithConstructorArguments(args []ConstructorArgument) BeanDefinitionOption {
	return &fnBeanDefinitionOption{
		fn: func(opt *beanDefinitionOption) {
			opt.construtorArguments = args
		},
	}
}

func (o *beanDefinitionOption) Validate() error {
	if o.scope != ScopePrototype && o.scope != ScopeSingleton {
		return xerrors.Errorf("Unsupport scope '%v'", o.scope)
	}

	return nil
}

func defaultBeanDefinitionOption() *beanDefinitionOption {
	return &beanDefinitionOption{
		scope: ScopeSingleton,
		lazy:  false,
	}
}
