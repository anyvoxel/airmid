// Package xapp provide the implement to Instantiate an application
package xapp

import (
	"context"

	_ "go.uber.org/automaxprocs" //nolint

	"github.com/anyvoxel/airmid/beans"
	"github.com/anyvoxel/airmid/props"
)

var (
	app = NewApplication()
)

// RegisterBeanDefinition wraps beans.BeanFactory.RegisterBeanDefinition function.
func RegisterBeanDefinition(beanName string, beanDefinition beans.BeanDefinition) error {
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
