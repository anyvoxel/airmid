package beans

import (
	"github.com/anyvoxel/airmid/xerrors"
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
