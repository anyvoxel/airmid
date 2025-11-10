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

// Package app provide the implement to Instantiate an application
package app

import (
	"context"

	_ "go.uber.org/automaxprocs" //nolint

	"github.com/anyvoxel/airmid/ioc"
	"github.com/anyvoxel/airmid/ioc/props"
)

var (
	app = NewApplication()
)

// RegisterBeanDefinition wraps beans.BeanFactory.RegisterBeanDefinition function.
func RegisterBeanDefinition(beanName string, beanDefinition ioc.BeanDefinition) error {
	return app.RegisterBeanDefinition(beanName, beanDefinition)
}

// Set wraps beans.BeanFactory.Set function.
func Set(key string, value any) error {
	return app.Set(key, value)
}

// Get wraps beans.BeanFactory.Get function.
func Get(key string, opts ...props.GetOption) (any, error) {
	return app.Get(key, opts...)
}

// Run wraps Application.Run function.
func Run(ctx context.Context, opts ...Option) error {
	return app.Run(ctx, opts...)
}

// Shutdown wraps Application.Shutdown function.
func Shutdown() {
	app.Shutdown()
}

// DefaultApp return the default application implement.
func DefaultApp() Application {
	return app
}
